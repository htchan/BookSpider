package model

import (
	"fmt"
	"log"
	"strconv"
	"errors"
	"os"
	"time"
	"runtime"
	"golang.org/x/text/encoding"
	"encoding/json"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"context"
	"sync"
	"golang.org/x/sync/semaphore"

	"github.com/htchan/BookSpider/helper"
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
	databaseLocation, DownloadLocation string
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
		if err != nil || !rows.Next() {
			if err == nil { err = errors.New("no record found") }
			log.Print(id, "---", err)
			continue //time.Sleep(1000)
		}
		rows.Scan(&siteName, &temp, &version, &title, &writer, &typeName,
					&lastUpdate, &lastChapter, &end, &download, &read);
		rows.Close()
	}
	if (siteName == "") {
		outputByte, err := json.Marshal(map[string]interface{} {
			"site": site.SiteName, "id": id, "retry": i,
			"message": "cannot load from database",
		})
		helper.CheckError(err)
		log.Println(string(outputByte))
	}
	book := Book{
		SiteName: site.SiteName,
		Id: id,						Version: version,
		Title: title,				Writer: writer,					Type : typeName,
		LastUpdate: lastUpdate,		LastChapter: lastChapter,
		EndFlag: end,				DownloadFlag: download,			ReadFlag: read,
		decoder: site.decoder,
		baseUrl: baseUrl,			downloadUrl: downloadUrl,		chapterUrl: site.metaChapterUrl,
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
	if s == nil { s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT)) }
	var wg sync.WaitGroup
	var siteName string
	var bookId int
	// prepare transaction and statements
	tx, err := site.database.Begin()
	helper.CheckError(err)
	save, err := site.database.Prepare("insert into books "+
					"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
					" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	update, err := site.database.Prepare("update books set name=?, writer=?, type=?,"+
					"date=?, chapter=?, end=?, download=?, read=? where site=? and num=? and version=?");
	helper.CheckError(err);
	// update all normal books
	rows, err := site.database.Query("SELECT site, num FROM books group by num order by date desc");
	helper.CheckError(err)
	for rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1)
		rows.Scan(&siteName, &bookId)
		book := Book{ Version: -1 }
		for i := 0; i < 10 && book.Version < 0; i += 1 { book = site.Book(bookId, -1) }
		if (book.Version == -1) {
			outputByte, err := json.Marshal(map[string]interface{} {
				"site": site.SiteName, "id": bookId, "version": book.Version,
				"message": "cannot load from database",
			})
			helper.CheckError(err)
			log.Println(string(outputByte))
		}
		go site.updateThread(book, s, &wg, tx, save, update);
	}
	rows.Close()
	wg.Wait()
	helper.CheckError(save.Close())
	helper.CheckError(update.Close())
	helper.CheckError(site.bookTx.Commit())
	helper.CheckError(tx.Commit())
	site.CloseDatabase()
}
func (site *Site) updateThread(book Book, s *semaphore.Weighted, wg *sync.WaitGroup, 
	tx *sql.Tx, save *sql.Stmt, update *sql.Stmt) {
	defer wg.Done()
	defer s.Release(1)
	checkVersion := book.Version;
	updated := book.Update();
	if (updated) {
		if (book.Version != checkVersion) {
			tx.Stmt(save).Exec(site.SiteName, book.Id, book.Version,
						book.Title, book.Writer, book.Type,
						book.LastUpdate, book.LastChapter,
						book.EndFlag, book.DownloadFlag, book.ReadFlag);
			book.Log(map[string]interface{} {
				"title": book.Title, "message": "new version updated", "stage": "update",
			})
		} else { // update old record
			tx.Stmt(update).Exec(book.Title, book.Writer, book.Type,
						book.LastUpdate, book.LastChapter,
						book.EndFlag, book.DownloadFlag, book.ReadFlag,
						book.SiteName, book.Id, book.Version);
			book.Log(map[string]interface{} {
				"message": "regular update", "stage": "update",
			})
			log.Println();
		}
	} else {
		// tell others nothing updated
		book.Log(map[string]interface{} {
			"message": "not updated", "stage": "update",
		})
		log.Println()
	}
}

func (site *Site) Explore(maxError int, s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	ctx := context.Background()
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	if s == nil { s = semaphore.NewWeighted(int64(maxError)) }
	var wg sync.WaitGroup
	// prepare transaction and statement
	tx, err := site.database.Begin();
	helper.CheckError(err)
	save, err := site.database.Prepare("insert into books "+
					"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
					" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	saveError, err := site.database.Prepare("insert into error(site, num) " +
					"select ?, ? from error where not exists(select 1 from error where num=?) limit 1");
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
	log.Println(maxId);
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
	helper.CheckError(deleteError.Close())
	helper.CheckError(saveError.Close())
	helper.CheckError(save.Close())
	helper.CheckError(site.bookTx.Commit())
	helper.CheckError(tx.Commit())
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
		book.Log(map[string]interface{} {
			"title": book.Title, "writer": book.Writer, "type": book.Type,
			"lastUpdate": book.LastUpdate, "lastChapter": book.LastChapter,
			"message": "explored", "stage": "explore",
		})
		*errorCount = 0;
	} else { // increase error Count
		_, err := tx.Stmt(saveError).Exec(book.SiteName, book.Id, book.Id)
		helper.CheckError(err)
		book.Log(map[string]interface{} {
			"message": "no such book", "stage": "explore",
		})
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
		if book.DownloadFlag { continue }
		book.Log(map[string]interface{} {
			"title": book.Title, "message": "start download", "stage": "download",
		})

		check := book.Download(site.DownloadLocation, site.MAX_THREAD_COUNT)
		if (! check) {
			book.Log(map[string]interface{} {
				"title": book.Title, "message": "download failure", "stage": "download",
			})
		} else {
			tx.Stmt(update).Exec(true, book.Id, book.Version)
		}
		runtime.GC()
	}
	helper.CheckError(rows.Close())
	helper.CheckError(site.bookTx.Commit())
	helper.CheckError(tx.Commit())
	site.CloseDatabase()
}

func (site *Site) UpdateError(s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	var err error
	ctx := context.Background()
	site.bookTx, err = site.database.Begin()
	if s == nil { s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT)) }
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
					" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
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
	helper.CheckError(save.Close())
	helper.CheckError(delete.Close())
	helper.CheckError(site.bookTx.Commit())
	helper.CheckError(tx.Commit())
	site.CloseDatabase()
}
func (site *Site) updateErrorThread(book Book, s *semaphore.Weighted, 
	wg *sync.WaitGroup, tx *sql.Tx, delete, save *sql.Stmt) {
	defer wg.Done()
	defer s.Release(1)
	updated := book.Update();
	if (updated) {
		// if update successfully
		tx.Stmt(save).Exec(site.SiteName, book.Id, book.Version,
					book.Title, book.Writer, book.Type,
					book.LastUpdate, book.LastChapter,
					book.EndFlag, book.DownloadFlag, book.ReadFlag);
		tx.Stmt(delete).Exec(site.SiteName, book.Id);
		book.Log(map[string]interface{} {
			"title": book.Title, "message": "error updated", "stage": "update",
		})
	} else {
		// tell others nothing updated
		book.Log(map[string]interface{} {
			"message": "error not updated", "stage": "update",
		})
	}
}

func (site *Site) Info() {
	site.OpenDatabase()
	log.Println("Site :\t" + site.SiteName);
	var normalCount, errorCount int;
	var rows *sql.Rows;
	rows, _ = site.database.Query("select count(DISTINCT num) as c from books");
	for rows.Next() { rows.Scan(&normalCount); }
	helper.CheckError(rows.Close())
	log.Println(site.SiteName, "Normal Book Count :\t" + strconv.Itoa(normalCount));
	rows, _ = site.database.Query("select count(num) as c from error");
	for rows.Next() { rows.Scan(&errorCount); }
	rows.Close()
	log.Println(site.SiteName, "Error Book Count :\t" + strconv.Itoa(errorCount));
	log.Println(site.SiteName, "Total Book Count :\t" + strconv.Itoa(normalCount + errorCount));
	
	maxId := site.getMaxBookId()
	log.Println("Max Book id :\t" + strconv.Itoa(maxId));
	site.CloseDatabase()
}

func (site *Site) CheckEnd() {
	site.OpenDatabase()
	tx, err :=site.database.Begin()
	helper.CheckError(err);
	matchingChapterCriteria := []string{"后记", "後記", "新书", "新書", "结局", "結局", "感言", 
                "尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本", "结束", "結束", "完結", 
				"完结", "终结", "終結", "番外", "结尾", "結尾", "全书完", "全書完", "全本完"}
	sqlStmt := "update books set end=true, download=false where ("
	for _, criteria := range matchingChapterCriteria { sqlStmt += "chapter like '%" + criteria + "%' or " }
	sqlStmt += "date < '" + strconv.Itoa(time.Now().Year() - 1) + "') and (end <> true or end is null)"
	result, err := tx.Exec(sqlStmt)
	helper.CheckError(err)
	rowAffect, err := result.RowsAffected()
	helper.CheckError(err)
	helper.CheckError(tx.Commit())
	log.Println(site.SiteName, "Row affected: ", rowAffect)
	site.CloseDatabase()
}

func (site *Site) RandomSuggestBook(size int) []Book {
	var downloadCount int
	site.OpenDatabase()
	tx, err := site.database.Begin()
	helper.CheckError(err)
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	rows, err := tx.Query("select count(*) from books where download=?", true)
	helper.CheckError(err)
	if err == nil && rows.Next() { rows.Scan(&downloadCount) }
	helper.CheckError(rows.Close())
	if (downloadCount < size) { size = downloadCount }
	var result = make([]Book, size)
	rows, err = tx.Query("select num, version from books where download=? order by random() limit ?", 
							true, size)
	helper.CheckError(err)
	var tempBookId, tempBookVersion int
	for i := 0; rows.Next() && i < size; i++ {
		rows.Scan(&tempBookId, &tempBookVersion)
		result[i] = site.Book(tempBookId, tempBookVersion)
	}
	helper.CheckError(rows.Close())
	helper.CheckError(site.bookTx.Commit())
	helper.CheckError(tx.Commit())
	site.CloseDatabase()
	return result;
}

func (site *Site) fixStroageError(s *semaphore.Weighted) {
	// init var for concurrency
	ctx := context.Background()
	if s == nil { s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT)) }
	var wg sync.WaitGroup
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	tx, err := site.database.Begin()
	helper.CheckError(err)
	// get book from database
	saveDownload, err := tx.Prepare("update books set end=?, download=? where num=? and version=?")
	helper.CheckError(err)
	saveNotDownload, err := tx.Prepare("update books set download=? where num=? and version=?")
	helper.CheckError(err)
	rows, err := tx.Query("select num, version from books")
	helper.CheckError(err)
	// loop all book
	var bookId, bookVersion int
	for rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1);
		rows.Scan(&bookId, &bookVersion)
		go site.CheckDownloadExistThread(site.Book(bookId, bookVersion), s, &wg, 
											tx, saveDownload, saveNotDownload)
	}
	wg.Wait()
	// commit changes to database
	helper.CheckError(rows.Close())
	helper.CheckError(saveDownload.Close())
	helper.CheckError(saveNotDownload.Close())
	helper.CheckError(tx.Commit())
	helper.CheckError(site.bookTx.Commit())
	site.bookTx = nil
}
func (site *Site)CheckDownloadExistThread(book Book, s *semaphore.Weighted, 
	wg *sync.WaitGroup, tx *sql.Tx, markDownload, markNotDownload *sql.Stmt) {
	defer wg.Done()
	defer s.Release(1)
	bookLocation := book.storageLocation(site.DownloadLocation)
	// check book file exist
	exist := helper.Exists(bookLocation)
	if exist && !book.DownloadFlag {
		// if book mark as not download, but it exist, mark as download
		tx.Stmt(markDownload).Exec(true, true, book.Id, book.Version)
		log.Println(site.SiteName + "\t" + strconv.Itoa(book.Id) + "\t" + 
					strconv.Itoa(book.Version) + "\t" + "mark to download")
	} else if !exist && book.DownloadFlag {
		// if book mark as download, but not exist, mark as not download
		tx.Stmt(markNotDownload).Exec(false, book.Id, book.Version)
		log.Println(site.SiteName + "\t" + strconv.Itoa(book.Id) + "\t" + 
					strconv.Itoa(book.Version) + "\t" + "mark to not download")
	}
}

func (site *Site) fixDatabaseDuplicateError() {
	// init variable
	tx, err := site.database.Begin()
	var bookId, bookVersion, bookCount, errorCount int
	// check any duplicate record in books table and show them
	rows, err := tx.Query("select num, version from books group by num, version having count(*) > 1")
	helper.CheckError(err)
	for rows.Next() {
		rows.Scan(&bookId, &bookVersion)
		log.Println(site.SiteName, "(" + strconv.Itoa(bookId) + ", " + strconv.Itoa(bookVersion) + ")")
		bookCount += 1
	}
	log.Println(site.SiteName, "duplicate book count : " + strconv.Itoa(bookCount))
	helper.CheckError(rows.Close())
	// delete duplicate record in book
	deleteStmt, err := tx.Prepare("delete from books where rowid not in " +
					"(select min(rowid) from books group by num, version)")
	helper.CheckError(err)
	_, err = tx.Stmt(deleteStmt).Exec()
	helper.CheckError(err)
	helper.CheckError(deleteStmt.Close())
	// check any duplicate record in error table and show them
	rows, err = tx.Query("select num from error group by num having count(*) > 1")
	helper.CheckError(err)
	for rows.Next() {
		rows.Scan(&bookId)
		log.Println(site.SiteName, "(" + strconv.Itoa(bookId) + ")")
	}
	log.Println(site.SiteName, "duplicate error count : " + strconv.Itoa(errorCount))
	helper.CheckError(rows.Close())
	// delete duplicate record
	deleteStmt, err = tx.Prepare("delete from error where rowid not in " +
					"(select min(rowid) from books group by site, num)")
	helper.CheckError(err)
	_, err = tx.Stmt(deleteStmt).Exec()
	helper.CheckError(err)
	helper.CheckError(deleteStmt.Close())
	// check if any record in book table duplicate in error table
	log.Println(site.SiteName, "duplicate cross - - - - - - - - - -")
	deleteStmt, err = tx.Prepare("delete from error where num in (select distinct num from books)")
	helper.CheckError(err)
	tx.Stmt(deleteStmt).Exec()
	helper.CheckError(deleteStmt.Close())
	helper.CheckError(tx.Commit())
}

func (site *Site) fixDatabaseMissingError(s *semaphore.Weighted) {
	// init variable
	missingBookIds := site.getMissingBookId()
	tx, err := site.database.Begin()
	// insert missing record by thread
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	if s == nil { s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT)) }
	ctx := context.Background()
	var wg sync.WaitGroup
	var errorCount int
	save, err := site.database.Prepare("insert into books "+
		"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
		" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	saveError, err := site.database.Prepare("insert into error (site, num) " +
					"select ?, ? from error where not exists(select 1 from error where num=?) limit 1");
	helper.CheckError(err);
	deleteError, err := site.database.Prepare("delete from error where site=? and num=?");
	helper.CheckError(err);
	log.Println(site.SiteName, "missing count : " + strconv.Itoa(len(missingBookIds)))
	for _, bookId := range missingBookIds {
		log.Println(site.SiteName, bookId)
		wg.Add(1)
		s.Acquire(ctx, 1);
		book := site.Book(bookId, -1)
		go site.exploreThread(book, &errorCount, s, &wg, tx, save, saveError, deleteError);
	}
	wg.Wait()
	helper.CheckError(deleteError.Close())
	helper.CheckError(saveError.Close())
	helper.CheckError(save.Close())
	helper.CheckError(site.bookTx.Commit())
	helper.CheckError(tx.Commit())
	// print missing record count
	log.Println(site.SiteName, "finish add missing count", len(missingBookIds))

}

func (site *Site) Fix(s *semaphore.Weighted) {
	site.OpenDatabase()
	log.Println(site.SiteName, "Add Missing Record")
	site.fixDatabaseMissingError(s)
	log.Println(site.SiteName, "Fix duplicate record")
	site.fixDatabaseDuplicateError()
	log.Println(site.SiteName, "Fix storage error")
	site.fixStroageError(s)
	log.Println()
	site.CloseDatabase()
}

func (site *Site) checkDuplicateBook() {
	// check duplicate record
	// check any duplicate record in books table and show them
	tx, err := site.database.Begin()
	rows, err := tx.Query("select num, version from books group by num, version having count(*) > 1")
	helper.CheckError(err)
	log.Print(site.SiteName, "duplicate books id : [")
	var id, version, count int
	for rows.Next() {
		if count > 0 { log.Println(", ") }
		rows.Scan(&id, &version)
		log.Print("(" + strconv.Itoa(id) + ", " + strconv.Itoa(version) + ")")
		count++
	}
	log.Println("]")
	log.Println(site.SiteName, "duplicate books count : " + strconv.Itoa(count))
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
	log.Print(site.SiteName, "duplicate error id : [")
	var id, count int
	for rows.Next() {
		if count > 0 { log.Print(", ") }
		rows.Scan(&id)
		log.Print(strconv.Itoa(id))
		count++
	}
	log.Println("]")
	log.Println(site.SiteName, "duplicate error count : " + strconv.Itoa(count))
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
	log.Print(site.SiteName, "duplicate cross id : [")
	var id, count int
	for rows.Next() {
		if count > 0 { log.Print(", ") }
		rows.Scan(&id)
		log.Print(strconv.Itoa(id))
		count++
	}
	log.Println("]")
	log.Println(site.SiteName, "duplicate cross count : " + strconv.Itoa(count))
	err = rows.Close()
	helper.CheckError(err)
	err = tx.Rollback()
	helper.CheckError(err)
}
func (site *Site) getMaxBookId() int {
	// get max id from database
	tx, err := site.database.Begin()
	var maxErrorId, maxBookId int
	rows, err := tx.Query("select num from books order by num desc")
	helper.CheckError(err)
	if rows.Next() { rows.Scan(&maxBookId) }
	helper.CheckError(rows.Close())
	rows, err = tx.Query("select num from error order by num desc")
	helper.CheckError(err)
	if rows.Next() { rows.Scan(&maxErrorId) }
	helper.CheckError(rows.Close())
	helper.CheckError(tx.Rollback())
	if maxBookId > maxErrorId {
		return maxBookId
	} else {
		return maxErrorId
	}
}
func (site *Site) getMissingBookId() []int {
	maxBookId := site.getMaxBookId()
	missingBookIds := make([]int, 0)
	tx, err := site.database.Begin()
	helper.CheckError(err)
	// check missing record
	rows, err := tx.Query("select num from " +
		"(select num from error union select num from books) order by num")
	helper.CheckError(err)
	var currentId, bookId int
	currentId = 1
	for rows.Next() {
		rows.Scan(&bookId)
		for ; currentId < bookId; currentId++ { missingBookIds = append(missingBookIds, currentId) }
		currentId++
	}
	for ; currentId < maxBookId; currentId++ { missingBookIds = append(missingBookIds, currentId) }
	helper.CheckError(rows.Close())
	helper.CheckError(tx.Rollback())
	return missingBookIds
}
func (site *Site) checkMissingId() {
	missingBookIds := site.getMissingBookId()
	jsonByte, err := json.Marshal(missingBookIds)
	helper.CheckError(err)
	log.Println(site.SiteName, "missing id : ", string(jsonByte))
	log.Println(site.SiteName, "missing count : " + strconv.Itoa(len(missingBookIds)))
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

func (site *Site) Search(title, writer string, page int) []Book {
	site.OpenDatabase()
	results := make([]Book, 0)
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	if title == "" && writer == "" { return results }
	const n = 20
	rows, err := site.database.Query("select num, version from books where "+
		"name like ? and writer like ? limit ?, ?", "%"+title+"%", "%"+writer+"%", page*n, n)
	helper.CheckError(err)
	var bookId, bookVersion int
	for rows.Next() {
		rows.Scan(&bookId, &bookVersion)
		results = append(results, site.Book(bookId, bookVersion))
	}
	helper.CheckError(rows.Close())
	helper.CheckError(site.bookTx.Commit())
	site.CloseDatabase()
	return results
}

func (site Site) Validate() float64 {
	os.Mkdir("./validate-download/", 0755)
	site.OpenDatabase()

	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)

	rows, err := site.database.Query("select num, version from books " +
		"where download=? order by random()", true)
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
	if tried / success > 90 { return -1 }
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
	rows, err := site.database.Query("select count(DISTINCT books.num), count(*) from books")
	if err == nil && rows.Next() { rows.Scan(&bookCount, &bookRecordCount) }
	rows.Close()
	rows, err = site.database.Query("select count(DISTINCT error.num), count(*) from error")
	if err == nil && rows.Next() { rows.Scan(&errorCount, &errorRecordCount) }
	rows.Close()
	rows, err = site.database.Query("select count(DISTINCT books.num), count(*) " +
		"from books where end=?", true)
	if err == nil && rows.Next() { rows.Scan(&endCount, &endRecordCount) }
	rows.Close()
	rows, err = site.database.Query("select count(DISTINCT books.num), count(*) " +
		"from books where download=?", true)
	if err == nil && rows.Next() { rows.Scan(&downloadCount, &downloadRecordCount) }
	rows.Close()
	rows, err = site.database.Query("select count(num) from books where read=?", true)
	if err == nil && rows.Next() { rows.Scan(&readCount) }
	rows.Close()
	rows, err = site.database.Query("select num from books order by num desc")
	if err == nil && rows.Next() { rows.Scan(&maxid) }
	rows.Close()
	rows, err = site.database.Query("select num from error order by num desc")
	if err == nil && rows.Next() {
		var temp int
		rows.Scan(&temp)
		if temp > maxid { maxid = temp }
	}
	rows.Close()
	site.CloseDatabase()
	return map[string]interface{} {
		"name": site.SiteName,
		"bookCount": bookCount,
		"errorCount": errorCount,
		"bookRecordCount": bookRecordCount,
		"errorRecordCount": errorRecordCount,
		"endCount": endCount,
		"endRecordCount": endRecordCount,
		"downloadCount": downloadCount,
		"downloadRecordCount": downloadRecordCount,
		"readCount": readCount,
		"maxid": maxid,
		"maxThread": site.MAX_THREAD_COUNT,
	}
}
