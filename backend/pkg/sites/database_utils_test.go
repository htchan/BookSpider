package sites

import (
	"testing"
	"os"
	"io"

	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

var testSite_database_utils = Site{
	SiteName: "ck101",
	database: nil,
	meta: *meta,
	decoder: traditionalchinese.Big5.NewDecoder(),
	databaseLocation: "../../test/site-test-data/ck101_database_utils.db",
	DownloadLocation: "./test_res/site-test-data/",
	MAX_THREAD_COUNT: 100,
}

func init() {
	source, err := os.Open("../../test/site-test-data/ck101_template.db")
	utils.CheckError(err)
	destination, err := os.Create("../../test/site-test-data/ck101_database_utils.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func TestDatabase(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	if testSite_database_utils.database == nil {
		t.Fatalf("Site.OpenDatabase() failed")
	}

	testSite_database_utils.CloseDatabase()
	if testSite_database_utils.database != nil {
		t.Fatalf("Site.CloseDatabase() failed")
	}
}

func Test_query(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	rows, tx := testSite_database_utils.query("select * from books")
	if rows == nil || tx == nil {
		t.Fatalf("Site.query(\"\") return nil value")
	}
	if rows.Err() != nil {
		t.Fatalf("Site.query(\"\") contains error %v", rows.Err())
	}
	closeQuery(rows, tx)
	if rows.Err() != nil {
		t.Fatalf("Site.closeQuery(\"\") contains error %v", rows.Err())
	}
}

func TestbookQuery(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	var err error
	testSite_database_utils.bookLoadTx, err = testSite_database_utils.database.Begin()
	rows, err := testSite_database_utils.bookQuery(" where num=?", 5)
	if err != nil {
		t.Fatalf("Site.bookQuery(\" where num=?\", 5) returns %v", err)
	}
	_, err = books.LoadBook(rows, testSite_database_utils.meta,
		testSite_database_utils.decoder, testSite_database_utils.CONST_SLEEP)
	if err != nil {
		t.Fatalf("books.LoadBook cannot load rows %v", err)
	}
	err = rows.Close()
	if err != nil {
		t.Fatalf("Site.bookQuery(\" where num=?\", 5) returns %v", err)
	}
	err = testSite_database_utils.bookLoadTx.Rollback()
	if err != nil {
		t.Fatalf("Site.bookQuery(\" where num=?\", 5) returns %v", err)
	}
}

func TestbookQueryGroup(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	var err error
	testSite_database_utils.bookLoadTx, err = testSite_database_utils.database.Begin()
	rows, err := testSite_database_utils.bookQuery(" group by testSite_database_utils, num")
	if err != nil {
		t.Fatalf("Site.bookQuery(\" where id=?\", 5) returns %v", err)
	}
	book, err := books.LoadBook(rows, testSite_database_utils.meta,
		testSite_database_utils.decoder, testSite_database_utils.CONST_SLEEP)
	if err != nil {
		t.Fatalf("books.LoadBook cannot load rows %v", err)
	}
	if book.Version != 1 {
		t.Fatalf("result book's version is %v, but not 1", book.Version)
	}
	err = rows.Close()
	if err != nil {
		t.Fatalf("Site.bookQuery(\" where id=?\", 5) returns %v", err)
	}
	err = testSite_database_utils.bookLoadTx.Rollback()
	if err != nil {
		t.Fatalf("Site.bookQuery(\" where id=?\", 5) returns %v", err)
	}
}

func Test_bookCount(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	distinctBookCount, bookCount := testSite_database_utils.bookCount()
	if distinctBookCount != 1 || bookCount != 2 {
		t.Fatalf("Site.bookCount() returns\n(%v, %v)\nbut not\n(1, 2)", distinctBookCount, bookCount)
	}
}

func Test_errorCount(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	distinctErrorCount, errorCount := testSite_database_utils.errorCount()
	if distinctErrorCount != 1 || errorCount != 1 {
		t.Fatalf("Site.errorCount() returns\n(%v, %v)\nbut not\n(1, 1)", distinctErrorCount, errorCount)
	}
}

func Test_downloadCount(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	distinctDownloadCount, downloadCount := testSite_database_utils.downloadCount()
	if distinctDownloadCount != 1 || downloadCount != 1 {
		t.Fatalf("Site.downloadCount() returns\n(%v, %v)\nbut not\n(1, 1)", distinctDownloadCount, downloadCount)
	}
}

func Test_endCount(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	distinctEndCount, endCount := testSite_database_utils.endCount()
	if distinctEndCount != 1 || endCount != 1 {
		t.Fatalf("Site.endCount() returns\n(%v, %v)\nbut not\n(1, 1)", distinctEndCount, endCount)
	}
}

func Test_readCount(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	distinctReadCount, readCount := testSite_database_utils.readCount()
	if distinctReadCount != 1 || readCount != 1 {
		t.Fatalf("Site.readCount() returns\n(%v, %v)\nbut not\n(1, 1)", distinctReadCount, readCount)
	}
}

func TestmaxBookId(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	maxBookId := testSite_database_utils.maxBookId()
	if maxBookId != 5 {
		t.Fatalf("Site.maxBookId() returns %v but not 5", maxBookId)
	}
}

func TestmaxErrorId(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	maxErrorId := testSite_database_utils.maxErrorId()
	if maxErrorId != 1 {
		t.Fatalf("Site.maxErrorId() returns %v but not 1", maxErrorId)
	}
}

func TestmaxId(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	maxId := testSite_database_utils.maxId()
	if maxId != 5 {
		t.Fatalf("Site.maxErrorId() returns %v but not 5", maxId)
	}
}

func TestStmt(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	testSite_database_utils.PrepareStmt()
	if testSite_database_utils.bookOperateTx == nil ||
		testSite_database_utils.insertBookStmt == nil ||
		testSite_database_utils.updateBookStmt == nil ||
		testSite_database_utils.insertErrorStmt == nil ||
		testSite_database_utils.deleteErrorStmt == nil {
		t.Fatalf("Site.PrepareStmt return some stmt as nil")
	}
	testSite_database_utils.CloseStmt()
	if testSite_database_utils.bookOperateTx != nil ||
		testSite_database_utils.insertBookStmt != nil ||
		testSite_database_utils.updateBookStmt != nil ||
		testSite_database_utils.insertErrorStmt != nil ||
		testSite_database_utils.deleteErrorStmt != nil {
		t.Fatalf("Site.CloseStmt return some stmt is not nil")
	}
}

func TestInsertBook(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	testSite_database_utils.PrepareStmt()
	var err error
	testSite_database_utils.bookLoadTx, err = testSite_database_utils.database.Begin()
	book := books.NewBook(testSite_database_utils.SiteName, 999,
		testSite_database_utils.meta, testSite_database_utils.decoder,
		testSite_database_utils.bookLoadTx)
	book.Title = "999_title"
	book.Writer = "999_writer"
	book.Type = "999_type"
	book.LastChapter = "999_chapter"
	book.LastUpdate = "2010-10-10"
	book.Version = 0
	testSite_database_utils.InsertBook(*book)
	utils.CheckError(err)
	rows, err := testSite_database_utils.bookQuery(" where num=?", 999)
	if err != nil {
		t.Fatalf("Site.bookQuery(\"where num = ?\", 999) fail %v", err)
	}
	if rows.Next() {
		book, err = books.LoadBook(rows, testSite_database_utils.meta,
			testSite_database_utils.decoder, testSite_database_utils.CONST_SLEEP)
	}
	if err != nil{
		t.Fatalf("fail %v", err)
	}
	if book.Title != "999_title" || book.Writer != "999_writer" || book.Type != "999_type" ||
		book.LastChapter != "999_chapter" || book.LastUpdate != "2010-10-10" {
		t.Fatalf("Site.InsertBook fail the book with id 999 has value %v", book.Map())
	}
	err = testSite_database_utils.bookLoadTx.Rollback()
	testSite_database_utils.CloseStmt()
}

func TestUpdateBook(t *testing.T) {
	testSite_database_utils.OpenDatabase()
	defer testSite_database_utils.CloseDatabase()
	testSite_database_utils.PrepareStmt()
	testSite_database_utils.bookLoadTx, _ = testSite_database_utils.database.Begin()
	rows, _ := testSite_database_utils.bookQuery(" where num=?", 999)
	var book *books.Book
	if rows.Next() {
		book, _ = books.LoadBook(rows, testSite_database_utils.meta,
			testSite_database_utils.decoder, testSite_database_utils.CONST_SLEEP)
	}
	book.Title = "999_new_title"
	testSite_database_utils.UpdateBook(*book)
	testSite_database_utils.bookLoadTx.Rollback()
	testSite_database_utils.bookOperateTx.Commit()
	testSite_database_utils.bookLoadTx, _ = testSite_database_utils.database.Begin()
	rows, _ = testSite_database_utils.bookQuery(" where num=?", 999)
	if rows.Next() {
		book, _ = books.LoadBook(rows, testSite_database_utils.meta,
			testSite_database_utils.decoder, testSite_database_utils.CONST_SLEEP)
	}
	if book.Title != "999_new_title" {
		t.Fatalf("Site.Update fail result %v", book.Map())
	}
	testSite_database_utils.bookLoadTx.Rollback()
}
