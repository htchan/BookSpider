package model

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/text/encoding"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
	"strconv"
	"strings"
	"../helper"
	"os"
	"time"
	//"runtime"
)

const SITE_MAX_THREAD = 1000;

type Site struct {
	SiteName string
	database *sql.DB
	MetaBaseUrl, metaDownloadUrl, metaChapterUrl, chapterPattern string
	decoder *encoding.Decoder
	titleRegex, writerRegex, typeRegex, lastUpdateRegex, lastChapterRegex string
	chapterUrlRegex, chapterTitleRegex string
	chapterContentRegex string
	databaseLocation, downloadLocation string
	bookTx *sql.Tx
}

func NewSite(siteName string, decoder *encoding.Decoder, configFileLocation string, databaseLocation string, downloadLocation string) (Site) {
	database, err := sql.Open("sqlite3", databaseLocation)
	helper.CheckError(err);
	database.SetMaxIdleConns(10);
	database.SetMaxOpenConns(99999);
	data, err := ioutil.ReadFile(configFileLocation)
	helper.CheckError(err);
	var info map[string]interface{};
	if err = json.Unmarshal(data, &info); err != nil {
        panic(err);
	}
	site := Site{
		SiteName: siteName,
		database: database,
		MetaBaseUrl: info["metaBaseUrl"].(string),
		metaDownloadUrl: info["metaDownloadUrl"].(string),
		metaChapterUrl: info["metaChapterUrl"].(string),
		chapterPattern: info["chapterPattern"].(string),
		decoder: decoder,
		titleRegex: info["titleRegex"].(string),
		writerRegex: info["writerRegex"].(string),
		typeRegex: info["typeRegex"].(string),
		lastUpdateRegex: info["lastUpdateRegex"].(string),
		lastChapterRegex: info["lastChapterRegex"].(string),
		chapterUrlRegex: info["chapterUrlRegex"].(string),
		chapterTitleRegex: info["chapterTitleRegex"].(string),
		chapterContentRegex: info["chapterContentRegex"].(string),
		databaseLocation: databaseLocation,
		downloadLocation: downloadLocation};
	return site;
}

func (site *Site) Book(id int) (Book) {
	baseUrl := fmt.Sprintf(site.MetaBaseUrl, id);
	downloadUrl := fmt.Sprintf(site.metaDownloadUrl, id);
	var siteName string;
	var temp int;
	version := -1;
	title := "";
	writer := "";
	typeName := "";
	lastUpdate := "";
	lastChapter := "";
	end := false;
	download := false;
	read := false;
	for i := 0; i < 1; i++ {
		rows, err := site.bookTx.Query("select site, num, version, name, writer, "+
						"type, date, chapter, end, download, read from books where "+
						"num="+strconv.Itoa(id) +
						" order by version desc");
		if (err != nil) {
			fmt.Println(id)
			fmt.Println(err)
			//time.Sleep(1000)
			continue
		}
		if (rows.Next()) {
			rows.Scan(&siteName, &temp, &version, &title, &writer, &typeName,
						&lastUpdate, &lastChapter, &end, &download, &read);
		} else {
			fmt.Println("retry (" + strconv.Itoa(i) + ") Cannot load " + strconv.Itoa(id) + " from database")
			//time.Sleep(1000)
			continue
		}
		rows.Close()
		break
		//panic(err)
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

func (site *Site) Update() () {
	// init concurrent variable
	ctx := context.Background()
	site.bookTx, _ = site.database.Begin()
	var s = semaphore.NewWeighted(int64(SITE_MAX_THREAD))
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
	update, err := site.database.Prepare("update books set version=?, name=?, writer=?, type=?,"+
					"date=?, chapter=?, end=?, download=?, read=? where site=? and num=?");
	helper.CheckError(err);
	// update all normal books
	rows, err := site.database.Query("SELECT site, num FROM books order by date desc");
	helper.CheckError(err)
	for rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1);
		rows.Scan(&siteName, &id);
		book := site.Book(id)
		if (book.Version == -1) {
			panic(strconv.Itoa(book.Version) + " cannot cache from database")
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
}
func (site *Site) updateThread(book Book, s *semaphore.Weighted, wg *sync.WaitGroup, tx *sql.Tx, save *sql.Stmt, update *sql.Stmt) () {
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
			fmt.Println("new version update " + strconv.Itoa(checkVersion) + " -> " + strconv.Itoa(book.Version) + " - - - - - - - -\n"+book.String());
			fmt.Println();
		} else { // update old record
			tx.Stmt(update).Exec(book.Version, book.Title, book.Writer, book.Type,
						book.LastUpdate, book.LastChapter,
						book.EndFlag, book.DownloadFlag, book.ReadFlag,
						book.SiteName, book.Id);
			fmt.Println("regular version update - - - - - - - - - -\n"+book.String());
			fmt.Println();
		}
	} else {
		// tell others nothing updated
		fmt.Println("Not updated - - - - - - - - - - - - - - -\n" + book.String())
		fmt.Println()
	}
}

func (site *Site) Explore(maxError int) () {
	// init concurrent variable
	ctx := context.Background()
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	var s = semaphore.NewWeighted(int64(SITE_MAX_THREAD))
	var wg sync.WaitGroup
	// prepare transaction and statement
	tx, err := site.database.Begin();
	helper.CheckError(err)
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
	deleteError, err := site.database.Prepare("delete from books "+
					"where site=? and num=?");
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
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
}
func (site *Site) exploreThread(book Book, errorCount *int, s *semaphore.Weighted, wg *sync.WaitGroup, tx *sql.Tx, save, saveError, deleteError *sql.Stmt) () {
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
		fmt.Println("Explore - - - - - - - - - - - -\n" + book.String())
		fmt.Println();
		*errorCount = 0;
	} else { // increase error Count
		tx.Stmt(saveError).Exec(book.SiteName, book.Id)
		_, err := fmt.Println("Unreachable - - - - - - - - - - -\n" + book.String());
		helper.CheckError(err)
		fmt.Println();
		*errorCount++;
	}
}

func (site *Site) Download() () {
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	rows, err := site.database.Query("select num from books where end=true and download=false")
	helper.CheckError(err)
	update, err := site.database.Prepare("update books set download=true where num=?")
	helper.CheckError(err)
	tx, err := site.database.Begin()
	helper.CheckError(err)
	var id int;
	if (rows.Next()) {
		rows.Scan(&id);
		book := site.Book(id);
		check := book.Download(site.downloadLocation)
		if (! check) {
			fmt.Println("download failure\t" + strconv.Itoa(book.Id) + "\t" + book.Title)
		} else {
			tx.Stmt(update).Exec(id)
		}
	}
	err = rows.Close()
	helper.CheckError(err)
	err = site.bookTx.Commit()
	helper.CheckError(err)
	err = tx.Commit()
	helper.CheckError(err)
}

func (site *Site) UpdateError() () {
	// init concurrent variable
	var err error
	ctx := context.Background()
	site.bookTx, err = site.database.Begin()
	var s = semaphore.NewWeighted(int64(SITE_MAX_THREAD))
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
		book := site.Book(id)
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
}
func (site *Site) updateErrorThread(book Book, s *semaphore.Weighted, wg *sync.WaitGroup, tx *sql.Tx, delete, save *sql.Stmt) () {
	defer wg.Done()
	defer s.Release(1)
	// try to update book
	//book := site.Book(id)
	updated := book.Update();
	if (updated) {
		// if update successfully
		tx.Stmt(delete).Exec(site.SiteName, book.Id);
		tx.Stmt(save).Exec(site.SiteName, book.Id, book.Version,
					book.Title, book.Writer, book.Type,
					book.LastUpdate, book.LastChapter,
					book.EndFlag, book.DownloadFlag, book.ReadFlag);
		fmt.Println("Error update - - - - - - - - - -\n"+book.String());
		fmt.Println();
	} else {
		// tell others nothing updated
		fmt.Println("Not updated - - - - - - - - - - -\n" + book.String())
		fmt.Println()
	}
}

func (site Site) Info() () {
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
	rows, _ = site.database.Query("select num from books order by num desc limit 1");
	for rows.Next() {
		rows.Scan(&normalCount);
	}
	rows.Close()
	rows, _ = site.database.Query("select num from error order by num desc limit 1");
	for rows.Next() {
		rows.Scan(&errorCount);
	}
	rows.Close()
	var max int;
	if (normalCount > errorCount) {
		max = normalCount;
	} else {
		max = errorCount;
	}
	fmt.Println("Max Book Number :\t" + strconv.Itoa(max));
}

func (site *Site) fixStroageError() () {

}

func (site *Site) fixDatabaseDuplicateError() () {
	// init variable
	tx, err := site.database.Begin()
	// check any duplicate record in books table and show them
	rows, err := tx.Query("select num, version, count(*) as c from books group by num, version order by c desc")
	helper.CheckError(err)
	booksDuplicate := make([]Book, 0)
	for rows.Next() {
		var book Book
		var count int
		rows.Scan(&book.Id, &book.Version, &count)
		if (count > 1) {
			booksDuplicate = append(booksDuplicate, book)
		}  
	}
	err = rows.Close()
	helper.CheckError(err)
	// TODO delete duplicate record
	stmt, err := tx.Prepare("delete from books where num=? and version=?")
	helper.CheckError(err)
	fmt.Println("duplicate book count : " + strconv.Itoa(len(booksDuplicate)))
	for _, book := range booksDuplicate {
		fmt.Println("duplicate book - - - - - - - - - -\n"+book.String())
		_, err := tx.Stmt(stmt).Exec(book.Id, book.Version)
		helper.CheckError(err)
	}
	err = stmt.Close()
	helper.CheckError(err)
	// check any duplicate record in books table and show them
	rows, err = tx.Query("select num, count(*) as c from error group by num order by c desc")
	helper.CheckError(err)
	errorDuplicate := make([]Book, 0)
	for rows.Next() {
		var book Book
		var count int
		rows.Scan(&book.Id, &count)
		if (count > 1) {
			errorDuplicate = append(errorDuplicate, book)
		}
	}
	err = rows.Close()
	helper.CheckError(err)
	// TODO delete duplicate record
	stmt, err = tx.Prepare("delete from error where num=?")
	helper.CheckError(err)
	fmt.Println("duplicate error count : " + strconv.Itoa(len(errorDuplicate)))
	for _, book := range errorDuplicate {
		fmt.Println("duplicate error - - - - - - - - - -\n"+book.String())
		tx.Stmt(stmt).Exec(book.Id)
	}
	err = stmt.Close()
	helper.CheckError(err)
	// TODO check if any record in book table duplicate in error table
	rows, err = tx.Query("select books.num from books, error where books.num=error.num")
	helper.CheckError(err)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		fmt.Println(id)
	}
	err = tx.Commit()
	helper.CheckError(err)
}

func (site *Site) fixDatabaseMissingError() () {
	// init variable
	tx, err := site.database.Begin()
	// check any id missing in the database
	rows, err := tx.Query("select num from books group by num")
	helper.CheckError(err)
	nums := make([]int, 0)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		nums = append(nums, id)
	}
	rows.Close()
	rows, err = tx.Query("select num from error group by num")
	helper.CheckError(err)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		if !(helper.Contains(nums, id)) {
			nums = append(nums, id)
		}
	}
	rows.Close()
	max := -1
	rows, err = tx.Query("select num from books order by num desc")
	helper.CheckError(err)
	if rows.Next() {
		var num int
		rows.Scan(&num)
		if (num > max) {
			max = num
		}
	}
	err = rows.Close()
	helper.CheckError(err)
	rows, err = tx.Query("select num from error order by num desc")
	helper.CheckError(err)
	if rows.Next() {
		var num int
		rows.Scan(&num)
		if (num > max) {
			max = num
		}
	}
	err = rows.Close()
	helper.CheckError(err)
	missingNum := make([]int, 0)
	for i := 1; i < max; i += 1 {
		if !(helper.Contains(nums, i)) {
			missingNum = append(missingNum, i)
		}
	}
	// insert missing record
	// init concurrent variable
	ctx := context.Background()
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	var s = semaphore.NewWeighted(int64(SITE_MAX_THREAD))
	var wg sync.WaitGroup
	var errorCount int
	save, err := site.database.Prepare("insert into books "+
		"(site, num, version, name, writer, type, date, chapter, end, download, read)"+
		" values "+
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)");
	helper.CheckError(err);
	saveError, err := site.database.Prepare("insert into error (site, num) values (?,?)");
	helper.CheckError(err);
	deleteError, err := site.database.Prepare("delete from error "+
		"where site=? and num=?");
	helper.CheckError(err);
	fmt.Println("start add missing count " + strconv.Itoa(len(missingNum)))
	for _, num := range missingNum {
		fmt.Println(num)
		wg.Add(1)
		s.Acquire(ctx, 1);
		book := site.Book(num)
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
	fmt.Println("finish add missing count "+strconv.Itoa(len(missingNum)))

}

func (site *Site) Fix() () {
	fmt.Println("Add Missing Record")
	site.fixDatabaseMissingError()
	fmt.Println("fix duplicate record")
	site.fixDatabaseDuplicateError()
	fmt.Println()
}

func (site *Site) Check() () {
	// init variable
	tx, err := site.database.Begin()
	// check duplicate record
	// check missing record
	rows, err := tx.Query("select num from books group by num")
	helper.CheckError(err)
	nums := make([]int, 0)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		nums = append(nums, id)
	}
	rows.Close()
	rows, err = tx.Query("select num from error group by num")
	helper.CheckError(err)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		if !(helper.Contains(nums, id)) {
			nums = append(nums, id)
		}
	}
	rows.Close()
	// get max num from database
	max := -1
	rows, err = tx.Query("select num from books order by num desc")
	helper.CheckError(err)
	if rows.Next() {
		var num int
		rows.Scan(&num)
		if (num > max) {
			max = num
		}
	}
	err = rows.Close()
	helper.CheckError(err)
	rows, err = tx.Query("select num from error order by num desc")
	helper.CheckError(err)
	if rows.Next() {
		var num int
		rows.Scan(&num)
		if (num > max) {
			max = num
		}
	}
	err = rows.Close()
	helper.CheckError(err)
	missingNum := make([]int, 0)
	for i := 1; i < max; i += 1 {
		if !(helper.Contains(nums, i)) {
			missingNum = append(missingNum, i)
		}
	}
	err = tx.Commit()
	helper.CheckError(err)
	fmt.Println("missing count :\t" + strconv.Itoa(len(missingNum)))
	fmt.Print("missing num :\t")
	fmt.Println(missingNum)
}

func (site *Site) Backup() (bool) {
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

func (site Site) Test() () {
	maxError := 5
	ctx := context.Background()
	var err error
	site.bookTx, err = site.database.Begin()
	helper.CheckError(err)
	var s = semaphore.NewWeighted(int64(SITE_MAX_THREAD))
	var wg sync.WaitGroup
	// prepare transaction and statement
	tx, err := site.database.Begin();
	helper.CheckError(err)
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
	deleteError, err := site.database.Prepare("delete from books "+
					"where site=? and num=?");
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
