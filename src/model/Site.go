package model

import (
	"fmt"
	"strconv"
	"strings"
	"os"
	"time"
	"runtime"
	"math/rand"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"encoding/json"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"context"
	"sync"
	"golang.org/x/sync/semaphore"

	"../helper"
)

const MAX_THREAD_COUNT = 1000;

type Site struct {
	SiteName string
	database *sql.DB
	metaBaseUrl, metaDownloadUrl, metaChapterUrl, chapterPattern string
	decoder *encoding.Decoder
	titleRegex, writerRegex, typeRegex, lastUpdateRegex, lastChapterRegex string
	chapterUrlRegex, chapterTitleRegex string
	chapterContentRegex string
	databaseLocation, downloadLocation string
	bookTx *sql.Tx
	MAX_THREAD_COUNT int
}

func (site *Site) Book(id, version int) Book {
	baseUrl := fmt.Sprintf(site.metaBaseUrl, id);
	downloadUrl := fmt.Sprintf(site.metaDownloadUrl, id);
	var siteName, title, writer, typeName, lastUpdate, lastChapter string;
	var temp int;
	var end, download, read bool
	if site.bookTx == nil {
		var err error
		site.OpenDatabase()
		defer site.CloseDatabase()
		site.bookTx, err = site.database.Begin()
		helper.CheckError(err)
		defer site.bookTx.Commit()
	}
	var i int;
	for i = 0; i < 2; i++ {
		var rows *sql.Rows
		var err error
		if version < 0 {
			rows, err = site.bookTx.Query("select site, num, version, name, writer, "+
							"type, date, chapter, end, download, read from books where "+
							"num=? order by version desc", id);
		} else {
			rows, err = site.bookTx.Query("select site, num, version, name, writer, "+
							"type, date, chapter, end, download, read from books where "+
							"num=? and version=?", id, version)
		}
		if (err != nil) {
			fmt.Print(id, "---", err)
			//time.Sleep(1000)
			continue
		}
		if (rows.Next()) {
			rows.Scan(&siteName, &temp, &version, &title, &writer, &typeName,
						&lastUpdate, &lastChapter, &end, &download, &read);
		} else {
			//time.Sleep(1000)
			continue
		}
		rows.Close()
		//panic(err)
	}
	if (siteName == "") {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": site.SiteName,
			"id": strconv.Itoa(id),
			"retry": strconv.Itoa(i),
			"message": "cannot load from database",
		})

		helper.CheckError(err)
		fmt.Println(string(strByte))
	}
	book := Book{
		SiteName: site.SiteName,
		Id: id,
		Version: version,
		Title: title,
		Writer: writer,
		Type : typeName,
		LastUpdate: lastUpdate,
		LastChapter: lastChapter,
		EndFlag: end,
		DownloadFlag: download,
		ReadFlag: read,
		decoder: site.decoder,
		baseUrl: baseUrl,
		downloadUrl: downloadUrl,
		chapterUrl: site.metaChapterUrl,
		chapterPattern: site.chapterPattern,
		titleRegex: site.titleRegex,
		writerRegex: site.writerRegex,
		typeRegex: site.typeRegex,
		lastUpdateRegex: site.lastUpdateRegex,
		lastChapterRegex: site.lastChapterRegex,
		chapterUrlRegex: site.chapterUrlRegex,
		chapterTitleRegex: site.chapterTitleRegex,
		chapterContentRegex: site.chapterContentRegex};
	return book;
}

func (site *Site) BookContent(book Book) string {
	if book.Title == "" || !book.DownloadFlag {
		return ""
	}
	path := site.downloadLocation + "/" + strconv.Itoa(book.Id)
	if book.Version > 0 {
		path += "-v" + strconv.Itoa(book.Version)
	}
	path += ".txt"
	content, err := ioutil.ReadFile(path)
	helper.CheckError(err)
	return string(content)
}

func (site *Site) OpenDatabase() {
	var err error
	site.database, err = sql.Open("sqlite3", site.databaseLocation)
	helper.CheckError(err)
	site.database.SetMaxIdleConns(10);
	site.database.SetMaxOpenConns(99999);
}
func (site *Site) CloseDatabase() {
	err := site.database.Close()
	helper.CheckError(err)
	site.database = nil
}

func (site *Site) Update(s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	ctx := context.Background()
	site.bookTx, _ = site.database.Begin()
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	var wg sync.WaitGroup
	var siteName string;
	var id int;
	// prepare transaction and statements
	tx, err := site.database.Begin()
	helper.CheckError(err)
	save, err := site.database.Prepare("insert into books "+
					"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
					" values "+
					"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	update, err := site.database.Prepare("update books set name=?, writer=?, type=?,"+
					"date=?, chapter=?, end=?, download=?, read=? where site=? and num=? and version=?");
	helper.CheckError(err);
	// update all normal books
	rows, err := site.database.Query("SELECT site, num FROM books group by num order by date desc");
	helper.CheckError(err)
	for rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1);
		rows.Scan(&siteName, &id);
		var book Book
		for i := 0; i < 10; i += 1 {
			book = site.Book(id, -1)
			if book.Version >= 0 {
				break
			}
		}
		if (book.Version == -1) {
			strByte, err := json.Marshal(map[string]interface{} {
				"site": site.SiteName,
				"id": strconv.Itoa(id),
				"version": strconv.Itoa(book.Version),
				"message": "cannot load from database",
			})

			helper.CheckError(err)
			fmt.Println(string(strByte))
		}
		go site.updateThread(book, s, &wg, tx, save, update);
	}
	rows.Close()
	wg.Wait()
	err = save.Close()
	helper.CheckError(err)
	err = update.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
	site.CloseDatabase()
}
func (site *Site) updateThread(book Book, s *semaphore.Weighted, wg *sync.WaitGroup, 
	tx *sql.Tx, save *sql.Stmt, update *sql.Stmt) {
	defer wg.Done()
	defer s.Release(1)
	//book := site.Book(id)
	checkVersion := book.Version;
	// try to update book
	updated := book.Update();
	if (updated) {
		// if version different, save a new record
		if (book.Version != checkVersion) {
			tx.Stmt(save).Exec(site.SiteName, book.Id, book.Version,
						book.Title, book.Writer, book.Type,
						book.LastUpdate, book.LastChapter,
						book.EndFlag, book.DownloadFlag, book.ReadFlag);
			
			strByte, err := json.Marshal(map[string]interface{} {
				"site": book.SiteName,
				"id": strconv.Itoa(book.Id),
				"version": strconv.Itoa(book.Version),
				"title": book.Title,
				"message": "new version updated",
			})

			helper.CheckError(err)
			fmt.Println(string(strByte))
		} else { // update old record
			tx.Stmt(update).Exec(book.Title, book.Writer, book.Type,
						book.LastUpdate, book.LastChapter,
						book.EndFlag, book.DownloadFlag, book.ReadFlag,
						book.SiteName, book.Id, book.Version);
			strByte, err := json.Marshal(map[string]interface{} {
				"site": book.SiteName,
				"id": strconv.Itoa(book.Id),
				"version": strconv.Itoa(book.Version),
				"message": "regular update",
			})

			helper.CheckError(err)
			fmt.Println(string(strByte))
			fmt.Println();
		}
	} else {
		// tell others nothing updated
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": strconv.Itoa(book.Id),
			"version": strconv.Itoa(book.Version),
			"message": "not updated",
		})
		
		helper.CheckError(err)
		fmt.Println(string(strByte))
		fmt.Println()
	}
}

func (site *Site) Explore(maxError int, s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	ctx := context.Background()
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	if s == nil {
		s = semaphore.NewWeighted(int64(maxError))
	}
	var wg sync.WaitGroup
	// prepare transaction and statement
	tx, err := site.database.Begin();
	helper.CheckError(err)
	save, err := site.database.Prepare("insert into books "+
					"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
					" values "+
					"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	saveError, err := site.database.Prepare("insert into error(site, num) select ?, ? from error where not exists(select 1 from error where num=?) limit 1");
	helper.CheckError(err);
	deleteError, err := site.database.Prepare("delete from error where site=? and num=?");
	helper.CheckError(err);
	// find max id
	rows, err := site.database.Query("select site, num from books order by num desc");
	helper.CheckError(err)
	var siteName string
	var maxId int;
	if (rows.Next()) {
		rows.Scan(&siteName, &maxId);
		maxId++;
	} else {
		maxId = 1;
	}
	rows.Close();
	fmt.Println(maxId);
	// keep explore until reach max error count
	errorCount := 0
	for (errorCount < maxError) {
		wg.Add(1)
		s.Acquire(ctx, 1);
		book := site.Book(maxId, -1);
		go site.exploreThread(book, &errorCount, s, &wg, tx, save, saveError, deleteError);
		maxId += 1;
	}
	wg.Wait()
	err = deleteError.Close()
	helper.CheckError(err)
	err = saveError.Close()
	helper.CheckError(err)
	err = save.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
	site.CloseDatabase()
}
func (site *Site) exploreThread(book Book, errorCount *int, s *semaphore.Weighted, 
	wg *sync.WaitGroup, tx *sql.Tx, save, saveError, deleteError *sql.Stmt) {
	defer wg.Done()
	defer s.Release(1)
	//book := site.Book(id)
	updated := book.Update();
	// if updated, save in books table, else, save in error table and **reset error count**
	if (updated) {
		_, err := tx.Stmt(save).Exec(book.SiteName, book.Id, book.Version,
					book.Title, book.Writer, book.Type,
					book.LastUpdate, book.LastChapter,
					book.EndFlag, book.DownloadFlag, book.ReadFlag);
		helper.CheckError(err)
		_, err = tx.Stmt(deleteError).Exec(book.SiteName, book.Id)
		helper.CheckError(err)
		strByte, err := json.Marshal(map[string]interface{} {
			"book": book.Map(),
			"message": "explored",
		})

		helper.CheckError(err)
		fmt.Println(string(strByte))
		*errorCount = 0;
	} else { // increase error Count
		_, err := tx.Stmt(saveError).Exec(book.SiteName, book.Id, book.Id)
		helper.CheckError(err)
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": strconv.Itoa(book.Id),
			"version": strconv.Itoa(book.Version),
			"message": "no such book",
		})

		helper.CheckError(err)
		fmt.Println(string(strByte))
		*errorCount++;
	}
}

func (site *Site) Download() {
	site.OpenDatabase()
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	rows, err := site.database.Query("select num, version from books where end=? and download=?", true, false)
	helper.CheckError(err)
	update, err := site.database.Prepare("update books set download=? where num=? and version=?")
	helper.CheckError(err)
	tx, err := site.database.Begin()
	helper.CheckError(err)
	var id, version int;
	for (rows.Next()) {
		rows.Scan(&id, &version);
		book := site.Book(id, version);
		if book.DownloadFlag {
			continue
		}
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": strconv.Itoa(book.Id),
			"version": strconv.Itoa(book.Version),
			"title": book.Title,
			"message": "start download",
		})

		helper.CheckError(err)
		fmt.Println(string(strByte))

		check := book.Download(site.downloadLocation, site.MAX_THREAD_COUNT)
		if (! check) {
			strByte, err := json.Marshal(map[string]interface{} {
				"site": book.SiteName,
				"id": strconv.Itoa(book.Id),
				"version": strconv.Itoa(book.Version),
				"title": book.Title,
				"message": "download failure",
			})

			helper.CheckError(err)
			fmt.Println(string(strByte))
		} else {
			tx.Stmt(update).Exec(true, book.Id, book.Version)
		}
		runtime.GC()
	}
	err = rows.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
	site.CloseDatabase()
}

func (site *Site) UpdateError(s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	var err error
	ctx := context.Background()
	site.bookTx, err = site.database.Begin()
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	var wg sync.WaitGroup
	var siteName string;
	var id int;
	// prepare transaction and statements
	tx, err := site.database.Begin()
	helper.CheckError(err)
	delete, err := site.database.Prepare("delete from error where site=? and num=?");
	helper.CheckError(err);
	save, err := site.database.Prepare("insert into books "+
					"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
					" values "+
					"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	// try update all error books
	rows, err := site.database.Query("SELECT site, num FROM error order by num desc")
	helper.CheckError(err)
	for rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1);
		rows.Scan(&siteName, &id);
		book := site.Book(id, -1)
		go site.updateErrorThread(book, s, &wg, tx, delete, save);
	}
	rows.Close()
	wg.Wait()
	err = save.Close()
	helper.CheckError(err)
	err = delete.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
	site.CloseDatabase()
}
func (site *Site) updateErrorThread(book Book, s *semaphore.Weighted, 
	wg *sync.WaitGroup, tx *sql.Tx, delete, save *sql.Stmt) {
	defer wg.Done()
	defer s.Release(1)
	// try to update book
	//book := site.Book(id)
	updated := book.Update();
	if (updated) {
		// if update successfully
		tx.Stmt(save).Exec(site.SiteName, book.Id, book.Version,
					book.Title, book.Writer, book.Type,
					book.LastUpdate, book.LastChapter,
					book.EndFlag, book.DownloadFlag, book.ReadFlag);
		tx.Stmt(delete).Exec(site.SiteName, book.Id);
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": strconv.Itoa(book.Id),
			"version": strconv.Itoa(book.Version),
			"title": book.Title,
			"message": "error updated",
		})
		
		helper.CheckError(err)
		fmt.Println(string(strByte))
	} else {
		// tell others nothing updated
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": strconv.Itoa(book.Id),
			"message": "error not updated",
		})
		
		helper.CheckError(err)
		fmt.Println(string(strByte))
	}
}

func (site *Site) Info() {
	site.OpenDatabase()
	fmt.Println("Site :\t" + site.SiteName);
	var normalCount, errorCount int;
	var rows *sql.Rows;
	rows, _ = site.database.Query("select count(DISTINCT num) as c from books");
	for rows.Next() {
		rows.Scan(&normalCount);
	}
	rows.Close()
	fmt.Println("Normal Book Count :\t" + strconv.Itoa(normalCount));
	rows, _ = site.database.Query("select count(num) as c from error");
	for rows.Next() {
		rows.Scan(&errorCount);
	}
	rows.Close()
	fmt.Println("Error Book Count :\t" + strconv.Itoa(errorCount));
	fmt.Println("Total Book Count :\t" + strconv.Itoa(normalCount + errorCount));
	
	maxId := site.getMaxId()
	fmt.Println("Max Book id :\t" + strconv.Itoa(maxId));
	site.CloseDatabase()
}

func (site *Site) CheckEnd() {
	site.OpenDatabase()
	tx, err :=site.database.Begin()
	helper.CheckError(err);
	criteria := []string{"后记", "後記", "新书", "新書", "结局", "結局", "感言", 
                "尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本",
                "结束", "結束", "完結", "完结", "终结", "終結", "番外",
				"结尾", "結尾", "全书完", "全書完", "全本完"}
	sql := "update books set end=true, download=false where ("
	for _, str := range criteria {
		sql += "chapter like '%" + str + "%' or "
	}
	sql += "date < '"+strconv.Itoa(time.Now().Year()-1)+"') and (end <> true or end is null)"
	result, err := tx.Exec(sql)
	helper.CheckError(err)
	rowAffect, err := result.RowsAffected()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
	fmt.Println("Row affected: ", rowAffect)
	site.CloseDatabase()
}

func (site *Site) Random(size int) []Book {
	var downloadCount int
	site.OpenDatabase()
	tx, err := site.database.Begin()
	site.bookTx, err = site.database.Begin()
	rows, err := tx.Query("select count(*) from books where download=?", true)
	if err == nil && rows.Next() {
		rows.Scan(&downloadCount)
	}
	rows.Close()
	if (downloadCount < size) { size = downloadCount }
	var result = make([]Book, size)
	var tempId, tempVersion int
	for i := 0; i < size; i++ {
		rows, err := tx.Query("select num, version from books where download=? order by num limit ?, 1", true, rand.Intn(downloadCount))
		if err == nil && rows.Next() {
			rows.Scan(&tempId, &tempVersion)
		}
		rows.Close()
		result[i] = site.Book(tempId, tempVersion)
	}
	site.bookTx.Commit()
	tx.Commit()
	site.CloseDatabase()
	return result;
}

func (site *Site) fixStroageError(s *semaphore.Weighted) {
	// init var for concurrency
	ctx := context.Background()
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	var wg sync.WaitGroup

	tx, err := site.database.Begin()
	helper.CheckError(err)
	// get book from database
	markDownload, err := tx.Prepare("update books set end=?, download=? where num=? and version=?")
	helper.CheckError(err)
	markNotDownload, err := tx.Prepare("update books set download=? where num=? and version=?")
	helper.CheckError(err)
	rows, err := tx.Query("select num, version, download from books")
	helper.CheckError(err)
	// loop all book
	var id, version int
	var recordDownload bool
	for rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1);
		rows.Scan(&id, &version, &recordDownload)
		go site.CheckDownloadExistThread(id, version, recordDownload, s, &wg, tx, markDownload, markNotDownload)
		/*
		path := site.downloadLocation + strconv.Itoa(id)
		if version > 0 {
			path += "-v" + strconv.Itoa(version)
		}
		path += ".txt"
		// check book file exist
		exist := helper.Exists(path)
		if exist && !recordDownload {
			// if book mark as not download, but it exist, mark as download
			tx.Stmt(markDownload).Exec(true, true, id, version)
			fmt.Println(site.SiteName + "\t" + strconv.Itoa(id) + "\t" + strconv.Itoa(version) + "\t" + "mark to download")
		} else if !exist && recordDownload {
			// if book mark as download, but not exist, mark as not download
			tx.Stmt(markNotDownload).Exec(false, id, version)
			fmt.Println(site.SiteName + "\t" + strconv.Itoa(id) + "\t" + strconv.Itoa(version) + "\t" + "mark to not download")
		}
		*/
	}
	wg.Wait()
	// commit changes to database
	err = rows.Close()
	helper.CheckError(err)
	err = markDownload.Close()
	helper.CheckError(err)
	err = markNotDownload.Close()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
}
func (site *Site)CheckDownloadExistThread(id, version int, recordDownload bool, 
	s *semaphore.Weighted, wg *sync.WaitGroup, tx *sql.Tx, markDownload, 
	markNotDownload *sql.Stmt) {
	defer wg.Done()
	defer s.Release(1)
	path := site.downloadLocation + strconv.Itoa(id)
	if version > 0 {
		path += "-v" + strconv.Itoa(version)
	}
	path += ".txt"
	// check book file exist
	exist := helper.Exists(path)
	if exist && !recordDownload {
		// if book mark as not download, but it exist, mark as download
		tx.Stmt(markDownload).Exec(true, true, id, version)
		fmt.Println(site.SiteName + "\t" + strconv.Itoa(id) + "\t" + strconv.Itoa(version) + "\t" + "mark to download")
	} else if !exist && recordDownload {
		// if book mark as download, but not exist, mark as not download
		tx.Stmt(markNotDownload).Exec(false, id, version)
		fmt.Println(site.SiteName + "\t" + strconv.Itoa(id) + "\t" + strconv.Itoa(version) + "\t" + "mark to not download")
	}
}

func (site *Site) fixDatabaseDuplicateError() {
	// init variable
	tx, err := site.database.Begin()
	// check any duplicate record in books table and show them
	rows, err := tx.Query("select num, version from books group by num, version having count(*) > 1")
	helper.CheckError(err)
	booksDuplicate := make([]Book, 0)
	for rows.Next() {
		var book Book
		rows.Scan(&book.Id, &book.Version)
		booksDuplicate = append(booksDuplicate, book)
	}
	err = rows.Close()
	helper.CheckError(err)
	// delete duplicate record in book
	stmt, err := tx.Prepare("delete from books where num=? and version=?")
	helper.CheckError(err)
	add, err := tx.Prepare("insert into books " +
				"(site, num, version, name, writer, type, date, chapter, end, download, read) " +
				"values " +
				"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	helper.CheckError(err)
	fmt.Println("duplicate book count : " + strconv.Itoa(len(booksDuplicate)))
	for _, book := range booksDuplicate {
		// fmt.Println("duplicate book - - - - - - - - - -\n"+book.String())
		fmt.Println("(" + book.SiteName + ", " + strconv.Itoa(book.Id) + ", " + strconv.Itoa(book.Version) + ")")
		_, err := tx.Stmt(stmt).Exec(book.Id, book.Version)
		helper.CheckError(err)
		tx.Stmt(stmt).Exec(book.SiteName, book.Id, book.Version,
					book.Title, book.Writer, book.Type,
					book.LastUpdate, book.LastChapter,
					book.EndFlag, book.DownloadFlag, book.ReadFlag)
		helper.CheckError(err)
	}
	err = stmt.Close()
	helper.CheckError(err)
	// check any duplicate record in error table and show them
	rows, err = tx.Query("select num from error group by num having count(*) > 1")
	helper.CheckError(err)
	errorDuplicate := make([]Book, 0)
	for rows.Next() {
		var book Book
		rows.Scan(&book.Id)
		errorDuplicate = append(errorDuplicate, book)
	}
	err = rows.Close()
	helper.CheckError(err)
	// delete duplicate record
	stmt, err = tx.Prepare("delete from error where num=?")
	helper.CheckError(err)
	add, err = tx.Prepare("insert into error (num) values (?)")
	helper.CheckError(err)
	fmt.Println("duplicate error count : " + strconv.Itoa(len(errorDuplicate)))
	for _, book := range errorDuplicate {
		// fmt.Println("duplicate error - - - - - - - - - -\n"+book.String())
		fmt.Println("(" + book.SiteName + ", " + strconv.Itoa(book.Id) + ")")
		_, err = tx.Stmt(stmt).Exec(book.Id)
		helper.CheckError(err)
		_, err = tx.Stmt(add).Exec(book.Id)
		helper.CheckError(err)
	}
	err = stmt.Close()
	helper.CheckError(err)
	// check if any record in book table duplicate in error table
	fmt.Println("duplicate cross - - - - - - - - - -\n")
	stmt, err = tx.Prepare("delete from error where num in (select distinct num from books)")
	helper.CheckError(err)
	tx.Stmt(stmt).Exec()
	err = stmt.Close()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
}

func (site *Site) fixDatabaseMissingError(s *semaphore.Weighted) {
	// init variable
	missingIds := site.getMissingId()
	tx, err := site.database.Begin()
	// insert missing record by thread
	ctx := context.Background()
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	var wg sync.WaitGroup
	var errorCount int
	save, err := site.database.Prepare("insert into books "+
		"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
		" values "+
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	saveError, err := site.database.Prepare("insert into error "+
					"(site, num)"+
					" values "+
					"(?, ?)");
	helper.CheckError(err);
	deleteError, err := site.database.Prepare("delete from error "+
		"where site=? and num=?");
	helper.CheckError(err);
	fmt.Println("missing count : " + strconv.Itoa(len(missingIds)))
	for _, id := range missingIds {
		fmt.Println(id)
		wg.Add(1)
		s.Acquire(ctx, 1);
		book := site.Book(id, -1)
		go site.exploreThread(book, &errorCount, s, &wg, tx, save, saveError, deleteError);
	}
	wg.Wait()
	err = deleteError.Close()
	helper.CheckError(err)
	err = saveError.Close()
	helper.CheckError(err)
	err = save.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
	// print missing record count
	fmt.Println("finish add missing count "+strconv.Itoa(len(missingIds)))

}

func (site *Site) Fix(s *semaphore.Weighted) {
	site.OpenDatabase()
	fmt.Println("Add Missing Record")
	site.fixDatabaseMissingError(s)
	fmt.Println("Fix duplicate record")
	site.fixDatabaseDuplicateError()
	fmt.Println("Fix storage error")
	site.fixStroageError(s)
	fmt.Println()
	site.CloseDatabase()
}

func (site *Site) checkDuplicateBook() {
	// check duplicate record
	// check any duplicate record in books table and show them
	tx, err := site.database.Begin()
	rows, err := tx.Query("select num, version from books group by num, version having count(*) > 1")
	helper.CheckError(err)
	fmt.Print("duplicate books id : [")
	var id, version, count int
	for rows.Next() {
		if count > 0 {
			fmt.Println(", ")
		}
		rows.Scan(&id, &version)
		fmt.Print("(" + strconv.Itoa(id) + ", " + strconv.Itoa(version) + ")")
		count++
	}
	fmt.Println("]")
	fmt.Println("duplicate books count : " + strconv.Itoa(count))
	err = rows.Close()
	helper.CheckError(err)
	err = tx.Rollback()
	helper.CheckError(err)
}
func (site *Site) checkDuplicateError() {
	// check duplicate record
	// check any duplicate record in books table and show them
	tx, err := site.database.Begin()
	rows, err := tx.Query("select num from error group by num having count(*) > 1")
	helper.CheckError(err)
	fmt.Print("duplicate error id : [")
	var id, count int
	for rows.Next() {
		if count > 0 {
			fmt.Print(", ")
		}
		rows.Scan(&id)
		fmt.Print(strconv.Itoa(id))
		count++
	}
	fmt.Println("]")
	fmt.Println("duplicate error count : " + strconv.Itoa(count))
	err = rows.Close()
	helper.CheckError(err)
	err = tx.Rollback()
	helper.CheckError(err)
}
func (site *Site) checkDuplicateCrossTable() {
	// check duplicate record
	// check if any record in book table duplicate in error table
	tx, err := site.database.Begin()
	rows, err := tx.Query("select distinct num from books where num in (select num from error)")
	helper.CheckError(err)
	fmt.Print("duplicate cross id : [")
	var id, count int
	for rows.Next() {
		if count > 0 {
			fmt.Print(", ")
		}
		rows.Scan(&id)
		fmt.Print(strconv.Itoa(id))
		count++
	}
	fmt.Println("]")
	fmt.Println("duplicate cross count : " + strconv.Itoa(count))
	err = rows.Close()
	helper.CheckError(err)
	err = tx.Rollback()
	helper.CheckError(err)
}
func (site *Site) getMaxId() int {
	// get max id from database
	tx, err := site.database.Begin()
	var maxErrorId, maxBookId int
	rows, err := tx.Query("select num from books order by num desc")
	helper.CheckError(err)
	if rows.Next() {
		rows.Scan(&maxBookId)
	}
	err = rows.Close()
	helper.CheckError(err)
	rows, err = tx.Query("select num from error order by num desc")
	helper.CheckError(err)
	if rows.Next() {
		rows.Scan(&maxErrorId)
	}
	err = rows.Close()
	helper.CheckError(err)
	err = tx.Rollback()
	helper.CheckError(err)
	if maxBookId > maxErrorId {
		return maxBookId
	} else {
		return maxErrorId
	}
}
func (site *Site) getMissingId() []int {
	maxId := site.getMaxId()
	missingIds := make([]int, 0)
	tx, err := site.database.Begin()
	helper.CheckError(err)
	// check missing record
	rows, err := tx.Query("select num from (select num from error union select num from books) order by num")
	helper.CheckError(err)
	var i, id int
	i = 1
	for rows.Next() {
		rows.Scan(&id)
		for ; i < id; i++ {
			missingIds = append(missingIds, i)
		}
		i++
	}
	for ; i < maxId; i++ {
		missingIds = append(missingIds, i)
	}
	err = rows.Close()
	helper.CheckError(err)
	err = tx.Rollback()
	helper.CheckError(err)
	return missingIds
}
func (site *Site) checkMissingId() {
	missingIds := site.getMissingId()
	b, err := json.Marshal(missingIds)
	helper.CheckError(err)
	fmt.Println("missing id : ", string(b))
	fmt.Println("missing count : " + strconv.Itoa(len(missingIds)))
}

func (site *Site) Check() {
	// init variable
	site.OpenDatabase()
	site.checkDuplicateBook()
	site.checkDuplicateError()
	site.checkDuplicateCrossTable()

	// check missing record
	site.checkMissingId()
	site.CloseDatabase()
}

func (site *Site) Backup() bool {
	y, m, d := time.Now().Date()
	folderName := fmt.Sprintf("%d-%d-%d/", y, m, d)
	path := strings.Replace(site.databaseLocation, "database/", "backup/"+folderName, -1)
	// create folder of today if not exist
	if _, err := os.Stat(strings.Replace(path, site.SiteName+".db", "", -1)); os.IsNotExist(err) {
		err := os.Mkdir(strings.Replace(path, site.SiteName+".db", "", -1), os.ModeDir+0755)
		helper.CheckError(err)
		fmt.Println(strings.Replace(path, site.SiteName+".db", "", -1)+" created")
	}
	// save as day-time-site.db for backup
	data, err := ioutil.ReadFile(site.databaseLocation)
	helper.CheckError(err)
	err = ioutil.WriteFile(path, data, 0644)
	helper.CheckError(err)
	return true
}
func (site *Site) table2StringSlice(table string) []string {
	result := make([]string, 0)
	rows, err := site.database.Query("select * from " + table)
	helper.CheckError(err)
	cols, _ := rows.Columns()
	containers := make([]string, len(cols))
	values := make([]interface{}, len(cols))
	for i, _ := range containers {
		values[i] = &containers[i]
	}
	for rows.Next() {
		rows.Scan(values...)
		insertStmt := "insert into " + table + " values ("
		for i, value := range values {
			insertStmt += "\"" + *value.(*string) + "\""
			if i < len(cols) - 1 {
				insertStmt += ", "
			}
		}
		insertStmt += ")"
		result = append(result, insertStmt + ";")
	}
	return result
}
func (site *Site) BackupString() bool {
	y, m, d := time.Now().Date()
	folderName := fmt.Sprintf("%d-%d-%d/", y, m, d)
	path := strings.Replace(
		strings.Replace(site.databaseLocation, "database/", "backup/"+folderName, -1),
		".db", ".sql", -1)
	// create folder of today if not exist
	if _, err := os.Stat(strings.Replace(path, site.SiteName+".sql", "", -1)); os.IsNotExist(err) {
		err := os.Mkdir(strings.Replace(path, site.SiteName+".sql", "", -1), os.ModeDir+0755)
		helper.CheckError(err)
		fmt.Println(strings.Replace(path, site.SiteName+".sql", "", -1)+" created")
	}
	// save as day-time-site.db for backup
	site.OpenDatabase()
	// load table name and sql
	tableNames := make([]string, 0)
	sqlStmts := make([]string, 0)
	rows, err := site.database.Query("SELECT name, sql FROM sqlite_master")
	helper.CheckError(err)
	for rows.Next() {
		var name, sql string
		rows.Scan(&name, &sql)
		tableNames = append(tableNames, name)
		sqlStmts = append(sqlStmts, sql + ";")
	}
	// load each row in table
	for _, tableName := range tableNames {
		sqlStmts = append(sqlStmts, site.table2StringSlice(tableName)...)
	}
	site.CloseDatabase()
	ioutil.WriteFile(path, []byte(strings.Join(sqlStmts, "\n")), 0644)
	return true
}

func (site *Site) Search(title, writer string, page int) []Book {
	site.OpenDatabase()
	results := make([]Book, 0)
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	if title == "" && writer == "" {
		return results
	}
	title = "%" + title + "%"
	writer = "%" + writer + "%"
	const n = 20
	rows, err := site.database.Query("select num, version from books where name like ? and writer like ? limit ?, ?", title, writer, page*n, n)
	helper.CheckError(err)
	var id, version int
	for rows.Next() {
		rows.Scan(&id, &version)
		results = append(results, site.Book(id, version))
	}
	err = rows.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	helper.CheckError(err)
	site.CloseDatabase()
	return results
}
/*
func (site Site) Test() {
	maxError := 5
	ctx := context.Background()
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	var wg sync.WaitGroup
	// prepare transaction and statement
	tx, err := site.database.Begin();
	helper.CheckError(err)
	save, err := site.database.Prepare("insert into books "+
					"(site, id, version, name, writer, type, date, chapter, end, download, read)"+
					" values "+
					"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	saveError, err := site.database.Prepare("insert into error "+
					"(site, id)"+
					" values "+
					"(?, ?)");
	helper.CheckError(err);
	deleteError, err := site.database.Prepare("delete from books "+
					"where site=? and id=?");
	helper.CheckError(err);
	// find max id
	rows, err := site.database.Query("select site, id from books order by id desc");
	helper.CheckError(err)
	var siteName string
	var maxId int;
	if (rows.Next()) {
		rows.Scan(&siteName, &maxId);
		maxId++;
	} else {
		maxId = 1;
	}
	rows.Close();
	fmt.Println(maxId);
	// keep explore until reach max error count
	errorCount := 0
	for (errorCount < maxError) {
		wg.Add(1)
		s.Acquire(ctx, 1);
		book := site.Book(maxId);
		go site.exploreThread(book, &errorCount, s, &wg, tx, save, saveError, deleteError);
		maxId++;
	}
	wg.Wait()
	err = deleteError.Close()
	helper.CheckError(err)
	err = saveError.Close()
	helper.CheckError(err)
	err = save.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	err = tx.Commit()
	helper.CheckError(err)
}
*/

func (site Site) Validate() float64 {
	os.Mkdir("./validate-download/", 0755)
	site.OpenDatabase()

	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)

	rows, err := site.database.Query("select num, version from books where download=? order by random()", true)
	helper.CheckError(err)

	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(20))
	var wg sync.WaitGroup

	success, tried := 0.0, 1.0
	for success < 10 && tried < 1000 && rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1);
		var num, version int
		rows.Scan(&num, &version)
		go site.validateThread(num, version, &success, &tried, s, &wg)
	}
	rows.Close()
	wg.Wait()
	site.CloseDatabase()
	os.RemoveAll("./validate-download/")
	if tried / success > 90 {
		return -1
	}
	return tried / success
}
func (site Site) validateThread(num int, version int, success *float64, 
	tried *float64, s *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	book := site.Book(num, version)
	book.Title = ""
	if *success < 10 && book.Update() && *success < 10 {
		*success ++
	} else {
		s.Release(1)
	}
	if *success < 10 {
		*tried ++
	}
}
func (site Site) ValidateDownload() float64 {
	os.Mkdir("./validate-download/", 0755)
	site.OpenDatabase()

	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)

	rows, err := site.database.Query("select num, version from books where download=? order by random()", true)
	helper.CheckError(err)

	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(4))
	var wg sync.WaitGroup

	success, tried := 0.0, 1.0
	for success < 2 && tried < 100 && rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1);
		var num, version int
		rows.Scan(&num, &version)
		go site.validateDownloadThread(num, version, &success, &tried, s, &wg)
	}
	rows.Close()
	wg.Wait()
	site.CloseDatabase()
	os.RemoveAll("./validate-download/")
	if tried / success > 90 {
		return -1
	}
	return tried / success
}
func (site Site) validateDownloadThread(num int, version int, success *float64, 
	tried *float64, s *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	book := site.Book(num, version)
	// here have two same condition because `book.Download` take a long time
	// the success may change after funush download
	if *success < 2 && book.Download("./validate-download/", site.MAX_THREAD_COUNT) && *success < 2 {
		*success ++
	} else {
		s.Release(1)
	}
	if *success < 2 {
		*tried ++
	}
}
func (site Site)Map() map[string]interface{} {
	site.OpenDatabase()
	var bookCount, bookRecordCount, errorCount, errorRecordCount, endCount, endRecordCount int
	var downloadCount, downloadRecordCount, readCount, maxid int
	rows, err := site.database.Query("select count(DISTINCT books.num) from books")
	if err == nil && rows.Next() {
		rows.Scan(&bookCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(*) from books")
	if err == nil && rows.Next() {
		rows.Scan(&bookRecordCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(DISTINCT error.num) from error")
	if err == nil && rows.Next() {
		rows.Scan(&errorCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(*) from error")
	if err == nil && rows.Next() {
		rows.Scan(&errorRecordCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(DISTINCT books.num) from books where end=?", true)
	if err == nil && rows.Next() {
		rows.Scan(&endCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(*) from books where end=?", true)
	if err == nil && rows.Next() {
		rows.Scan(&endRecordCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(DISTINCT books.num) from books where download=?", true)
	if err == nil && rows.Next() {
		rows.Scan(&downloadCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(*) from books where download=?", true)
	if err == nil && rows.Next() {
		rows.Scan(&downloadRecordCount)
	}
	rows.Close()
	rows, err = site.database.Query("select count(num) from books where read=?", true)
	if err == nil && rows.Next() {
		rows.Scan(&readCount)
	}
	rows.Close()
	rows, err = site.database.Query("select num from books order by num desc")
	if err == nil && rows.Next() {
		rows.Scan(&maxid)
	}
	rows.Close()
	rows, err = site.database.Query("select num from error order by num desc")
	if err == nil && rows.Next() {
		var temp int
		rows.Scan(&temp)
		if temp > maxid {
			maxid = temp
		}
	}
	rows.Close()
	site.CloseDatabase()
	return map[string]interface{} {
		"name": site.SiteName + "\"",
		"bookCount": strconv.Itoa(bookCount),
		"errorCount": strconv.Itoa(errorCount),
		"bookRecordCount": strconv.Itoa(bookRecordCount),
		"errorRecordCount": strconv.Itoa(errorRecordCount),
		"endCount": strconv.Itoa(endCount),
		"endRecordCount": strconv.Itoa(endRecordCount),
		"downloadCount": strconv.Itoa(downloadCount),
		"downloadRecordCount": strconv.Itoa(downloadRecordCount),
		"readCount": strconv.Itoa(readCount),
		"maxid": strconv.Itoa(maxid),
		"maxThread": site.MAX_THREAD_COUNT,
	}
}
