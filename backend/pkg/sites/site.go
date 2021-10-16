package sites

import (
	"fmt"
	"golang.org/x/text/encoding"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"errors"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/sync/semaphore"
	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

const MAX_THREAD_COUNT = 1000

type Site struct {
	SiteName                                                         string
	database                                                         *sql.DB
	decoder                                                          *encoding.Decoder
	meta                                                             books.MetaInfo
	databaseLocation, DownloadLocation                               string
	bookLoadTx, bookOperateTx                                        *sql.Tx
	insertBookStmt, updateBookStmt, insertErrorStmt, deleteErrorStmt *sql.Stmt
	MAX_THREAD_COUNT                                                 int
	CONST_SLEEP                                                      int
	semaphore                                                        *semaphore.Weighted
}

func LoadSite(siteName string, source map[string]string, metaMap map[string]string) (*Site, error) {
	expectKey := [4]string{ "databaseLocation", "downloadLocation", "threadsCount", "decode" }
	for _, key := range expectKey {
		if _, ok := source[key]; !ok {
			return nil, errors.New("missing key " + key)
		}
	}
	site := new(Site)
	site.SiteName = siteName
	site.databaseLocation = source["databaseLocation"]
	site.DownloadLocation = source["downloadLocation"]
	site.MAX_THREAD_COUNT, _ = strconv.Atoi(source["threadsCount"])
	if (source["decode"] == "big5") {
		site.decoder = traditionalchinese.Big5.NewDecoder()
	}
	meta, err := books.NewMetaInfo(metaMap)
	if err != nil {
		return site, err
	}
	site.meta = *meta
	return site, nil
}

func (site Site) Book(id, version int) (*books.Book, error) {
	site.OpenDatabase()
	var rows *sql.Rows
	var book *books.Book
	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)
	if version >= 0 {
		rows, err = site.bookQuery(" where site=? and num=? and version=?", site.SiteName, id, version)
	} else {
		rows, err = site.bookQuery(" where site=? and num=? group by site, num", site.SiteName, id)
	}
	utils.CheckError(err)
	if rows.Next() {
		book, err = books.LoadBook(rows, site.meta, site.decoder, site.CONST_SLEEP)
	} else {
		err = errors.New("book not found")
	}
	utils.CheckError(rows.Close())
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseDatabase()
	return book, err
}

func (site *Site) Info() {
	site.OpenDatabase()
	log.Println("Site :\t" + site.SiteName)
	distinctBookcount, _ := site.bookCount()
	distinctErrorCount, _ := site.errorCount()
	log.Println(site.SiteName, "Normal Book Count :\t", strconv.Itoa(distinctBookcount))
	log.Println(site.SiteName, "Error Book Count :\t", strconv.Itoa(distinctErrorCount))
	log.Println(site.SiteName, "Total Book Count :\t",
		strconv.Itoa(distinctBookcount+distinctErrorCount))

	log.Println("Max Book id :\t" + strconv.Itoa(site.maxBookId()))
	site.CloseDatabase()
}

func (site *Site) RandomSuggestBook(size int) []*books.Book {
	// find number of book to random
	site.OpenDatabase()
	_, downloadCount := site.downloadCount()
	if downloadCount < size {
		size = downloadCount
	}
	// random the book and put it to array
	var result = make([]*books.Book, size)
	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)
	rows, err := site.bookQuery(" where download=? order by random() limit ?", true, size)
	for i := 0; rows.Next() && i < size; i++ {
		result[i], err = books.LoadBook(rows, site.meta, site.decoder, site.CONST_SLEEP)
		if err != nil {
			result[i].Log(map[string]interface{}{
				"error": "cannot load book from database", "stage": "random",
			})
		}
	}
	utils.CheckError(rows.Close())
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseDatabase()
	return result
}

func (site *Site) Search(title, writer string, page int) []*books.Book {
	site.OpenDatabase()
	results := make([]*books.Book, 0)
	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)
	// return empty if all search fields are empty
	if title == "" && writer == "" {
		return results
	}
	const n = 20
	rows, err := site.bookQuery(" where name like ? and writer like ? limit ?, ?",
		"%"+title+"%", "%"+writer+"%", page*n, n)
	utils.CheckError(err)
	// load the book match search requirements
	for rows.Next() {
		book, err := books.LoadBook(rows, site.meta, site.decoder, site.CONST_SLEEP)
		if err != nil {
			book.Log(map[string]interface{}{
				"error": err.Error(), "stage": "search",
			})
			continue
		}
		results = append(results, book)
	}
	utils.CheckError(rows.Close())
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseDatabase()
	return results
}

func (site Site) BackupSql(destinationDirectory string) {
	// construct and create the backup path
	destinationFileName := filepath.Join(destinationDirectory, time.Now().Format("2006-01-02"))
	os.MkdirAll(destinationFileName, os.ModePerm)
	destinationFileName = filepath.Join(destinationFileName, site.SiteName+".sql")
	// the create table sql
	createBooksTableSql := "CREATE TABLE books " +
		"( `name` varchar ( 50 ), `writer` varchar ( 30 ), " +
		"`date` varchar ( 30 ), `chapter` varchar ( 50 ), `type` varchar ( 20 ), " +
		"`end` boolean, `download` boolean, `read` boolean, " +
		"`site` varchar ( 15 ), `num` INTEGER , version integer)"
	createErrorTableSql := "CREATE TABLE error " +
		"( `type` varchar ( 10 ), `site` varchar ( 15 ), `num` INTEGER )"
	site.OpenDatabase()
	file, err := os.OpenFile(destinationFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	utils.CheckError(err)
	// write create table sql
	file.WriteString(fmt.Sprintf("%v\n%v\n", createBooksTableSql, createErrorTableSql))
	rows, err := site.database.Query("select name, writer, date, chapter, type, " +
		"end, download, read, site, num, version from books")
	utils.CheckError(err)
	var tempName, tempWriter, tempDate, tempChapter, tempType, tempSite string
	var tempEnd, tempDownload, tempRead bool
	var tempNum, tempVersion int
	// write insert book record sql
	for rows.Next() {
		rows.Scan(&tempName, &tempWriter, &tempDate, &tempChapter, &tempType,
			&tempEnd, &tempDownload, &tempRead, &tempSite, &tempNum, &tempVersion)
		file.WriteString(fmt.Sprintf("insert into books values ("+
			"'%v', '%v', '%v', '%v', '%v', %v, %v, %v, '%v', %v, %v);\n",
			tempName, tempWriter, tempDate, tempChapter, tempType,
			tempEnd, tempDownload, tempRead, tempSite, tempNum, tempVersion))
	}
	rows, err = site.database.Query("select type, site, num from books")
	utils.CheckError(err)
	// write insert error record sql
	for rows.Next() {
		rows.Scan(&tempType, &tempSite, &tempNum)
		file.WriteString(fmt.Sprintf("insert into error values ('%v', '%v', %v);\n",
			tempType, tempSite, tempNum))
	}
	file.Close()
	site.CloseDatabase()
}

func (site Site) Map() map[string]interface{} {
	site.OpenDatabase()
	distinctBookCount, bookCount := site.bookCount()
	distinctErrorCount, errorCount := site.errorCount()
	distinctEndCount, endCount := site.endCount()
	distinctDownloadCount, downloadCount := site.downloadCount()
	maxId := site.maxId()
	site.CloseDatabase()
	return map[string]interface{}{
		"name":                site.SiteName,
		"bookCount":           distinctBookCount,
		"errorCount":          distinctErrorCount,
		"bookRecordCount":     bookCount,
		"errorRecordCount":    errorCount,
		"endCount":            distinctEndCount,
		"endRecordCount":      endCount,
		"downloadCount":       distinctDownloadCount,
		"downloadRecordCount": downloadCount,
		"maxid":               maxId,
		"maxThread":           site.MAX_THREAD_COUNT,
	}
}
