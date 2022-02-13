package books

import (
	"os"
	"io"
	"testing"
	"errors"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/database/sqlite"
)

func initBookTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./book_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupBookTest() {
	os.Remove("./book_test.db")
}

var config = configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml").SiteConfigs["test"].BookMeta

func bookEqual(
	book Book, site string, id, hash int, title string, writerId int,
	typeString, updateDate, updateChapter string, status database.StatusCode) bool {
	return book.bookRecord.Equal(database.BookRecord{
		Site: site, Id: id, HashCode: hash, Title: title, WriterId: writerId, Type: typeString,
		UpdateDate: updateDate, UpdateChapter: updateChapter, Status: status,
	})
}
func writerEqual(book Book, id int, name string) bool {
	return book.writerRecord.Equal(database.WriterRecord { Id: id, Name: name })
}

func Test_Books_Book_Constructor(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	t.Run("func NewBook", func(t *testing.T) {
		t.Parallel()
		t.Run("success", func(t *testing.T) {
			book := NewBook("test", 1, 500, config)
			if book.bookRecord == nil ||
				!bookEqual(*book,
					"test", 1, 500, "", 0, "", "", "", database.Error) {
				t.Fatalf("book.bookRecord init wrong record: %v", book.bookRecord)
			}
			if book.writerRecord == nil || !writerEqual(*book, 0, "") {
				t.Fatalf("book.writerRecord init wrong record: %v", book.writerRecord)
			}
			if book.errorRecord != nil {
				t.Fatalf("book.errorRecord was initialized error: %v", book.errorRecord)
			}

			if book.config.BaseUrl != "https://base-url/1" ||
				book.config.DownloadUrl != "https://download-url/1" ||
				book.config.ChapterUrl != "https://chapter-url/%v" ||
				book.config.ChapterUrlPattern != "chapter-url-pattern" ||
				book.config.TitleRegex != "(title-regex)" ||
				book.config.WriterRegex != "(writer-regex)" ||
				book.config.TypeRegex != "(type-regex)" ||
				book.config.LastUpdateRegex != "(last-update-regex)" ||
				book.config.LastChapterRegex != "(last-chapter-regex)" ||
				book.config.ChapterUrlRegex != "chapter-url-regex-(\\d)" ||
				book.config.ChapterTitleRegex != "chapter-title-regex-(\\d)" ||
				book.config.ChapterContentRegex != "chapter-content-(.*)-content-regex" ||
				book.config.Decoder == nil ||
				book.config.CONST_SLEEP != 1000 ||
				book.config.StorageDirectory != "/test-data/storage/" {
				t.Fatalf("book.config init wrongly config: %v", book.config)
			}
		})
	})

	t.Run("func LoadBook", func(t *testing.T) {
		t.Run("success with hash code > 1", func(t *testing.T) {
			t.Parallel()
			book := LoadBook(db, "test", 3, 102, config)
			
			if book.bookRecord == nil ||
				!bookEqual(*book,
					"test", 3, 102, "title-3", 2, "type-3",
					"102", "chapter-3", database.Download) {
				t.Fatalf("LoadBook load the wrong book record: %v", book.bookRecord)
			}

			if book.writerRecord == nil || !writerEqual(*book, 2, "writer-2") {
				t.Fatalf("LoadBook load the wrong writer record: %v", book.writerRecord)
			}

			if book.errorRecord != nil {
				t.Fatalf("LoadBook load the error record: %v", book.errorRecord)
			}

			if book.config.BaseUrl != "https://base-url/3" ||
				book.config.DownloadUrl != "https://download-url/3" {
				t.Fatalf("book.config init wrongly config: %v", book.config)
			}
		})

		t.Run("success for getting latest book", func(t *testing.T) {
			t.Parallel()
			book := LoadBook(db, "test", 3, -1, config)
			
			if book.bookRecord == nil ||
				!bookEqual(*book,
					"test", 3, 200, "title-3-new", 3, "type-3-new",
					"100", "chapter-3-new", database.End) {
				t.Fatalf("LoadBook load the wrong book record: %v", book.bookRecord)
			}

			if book.writerRecord == nil || !writerEqual(*book, 3, "writer-3") {
				t.Fatalf("LoadBook load the wrong writer record: %v", book.writerRecord)
			}

			if book.errorRecord != nil {
				t.Fatalf("LoadBook load the error record: %v", book.errorRecord)
			}

			if book.config.BaseUrl != "https://base-url/3" ||
				book.config.DownloadUrl != "https://download-url/3" {
				t.Fatalf("book.config init wrongly config: %v", book.config)
			}
		})

		t.Run("success for getting book with error", func(t *testing.T) {
			t.Parallel()
			book := LoadBook(db, "test", 2, 101, config)
			
			if book.bookRecord == nil ||
				!bookEqual(*book,
					"test", 2, 101, "", 0, "", "", "", database.Error) {
				t.Fatalf("LoadBook load the wrong book record: %v", book.bookRecord)
			}

			if book.writerRecord == nil || !writerEqual(*book, 0, "")  {
				t.Fatalf("LoadBook load the wrong writer record: %v", book.writerRecord)
			}

			if book.errorRecord == nil || book.errorRecord.Site != "test" ||
				book.errorRecord.Id != 2 || book.errorRecord.Error.Error() != "error-2" {
				t.Fatalf("LoadBook load the error record: %v", book.errorRecord)
			}

			if book.config.BaseUrl != "https://base-url/2" ||
				book.config.DownloadUrl != "https://download-url/2" {
				t.Fatalf("book.config init wrongly config: %v", book.config)
			}
		})

		t.Run("Fail if book not exist", func(t *testing.T) {
			t.Parallel()
			book := LoadBook(db, "not-exist", 1, -1, config)

			if book != nil {
				t.Fatalf(
					"LoadBook load sth when book not exist\nbook: %v\nwriter: %v\nerror: %v",
					book.bookRecord, book.writerRecord, book.errorRecord)
			}
		})
	})

	t.Run("func LoadBookByRecord", func(t *testing.T) {
		t.Run("success with some record", func(t *testing.T) {
			t.Parallel()
			record := &database.BookRecord{ Site: "test", Id: 3, HashCode: 102, WriterId: 2 }
			book := LoadBookByRecord(db, record, config)
			
			if book.bookRecord == nil || 
				!bookEqual(*book, 
					"test", 3, 102, "", 2, "", "", "", database.Error) {
				t.Fatalf("LoadBook load the wrong book record: %v", book.bookRecord)
			}

			if book.writerRecord == nil || !writerEqual(*book, 2, "writer-2") {
				t.Fatalf("LoadBook load the wrong writer record: %v", book.writerRecord)
			}

			if book.errorRecord != nil {
				t.Fatalf("LoadBook load the error record: %v", book.errorRecord)
			}

			if book.config.BaseUrl != "https://base-url/3" ||
				book.config.DownloadUrl != "https://download-url/3" {
				t.Fatalf("book.config init wrongly config: %v", book.config)
			}
		})

		t.Run("success for getting book with error", func(t *testing.T) {
			t.Parallel()
			record := &database.BookRecord{ Site: "test", Id: 2, HashCode: 101, Title: "hello" }
			book := LoadBookByRecord(db, record, config)
			
			if book.bookRecord == nil || 
				!bookEqual(*book, 
					"test", 2, 101, "hello", 0, "", "", "", database.Error) {
				t.Fatalf("LoadBook load the wrong book record: %v", book.bookRecord)
			}

			if book.writerRecord == nil || !writerEqual(*book, 0, "") {
				t.Fatalf("LoadBook load the wrong writer record: %v", book.writerRecord)
			}

			if book.errorRecord == nil || book.errorRecord.Site != "test" ||
				book.errorRecord.Id != 2 || book.errorRecord.Error.Error() != "error-2" {
				t.Fatalf("LoadBook load the error record: %v", book.errorRecord)
			}

			if book.config.BaseUrl != "https://base-url/2" ||
				book.config.DownloadUrl != "https://download-url/2" {
				t.Fatalf("book.config init wrongly config: %v", book.config)
			}
		})

		t.Run("success even if book not exist", func(t *testing.T) {
			t.Parallel()
			record := &database.BookRecord{ Site: "not-exist", Id: 3, HashCode: 102 }
			book := LoadBookByRecord(db, record, config)
			
			if book.bookRecord == nil || 
				!bookEqual(*book, 
					"not-exist", 3, 102, "", 0, "", "", "", database.Error) {
				t.Fatalf("LoadBook load the wrong book record: %v", book.bookRecord)
			}

			if book.writerRecord == nil || !writerEqual(*book, 0, "") {
				t.Fatalf("LoadBook load the wrong writer record: %v", book.writerRecord)
			}

			if book.errorRecord != nil {
				t.Fatalf("LoadBook load the error record: %v", book.errorRecord)
			}

			if book.config.BaseUrl != "https://base-url/3" ||
				book.config.DownloadUrl != "https://download-url/3" {
				t.Fatalf("book.config init wrongly config: %v", book.config)
			}
		})
	})
}

func Test_Books_Book_Save(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()

	t.Run("success to create a completely new book with writer", func(t *testing.T) {
		book := NewBook("test", 1000, -1, config)
		book.SetWriter("new-writer")
		book.SetStatus(database.InProgress)
		result := book.Save(db)
		db.Commit()
		if !result {
			t.Fatalf("Save fail")
		}

		query := db.QueryWriterByName("new-writer")
		defer query.Close()
		record, err := query.Scan()
		if err != nil || record.(*database.WriterRecord).Id != 4 ||
			record.(*database.WriterRecord).Name != book.GetWriter() {
				t.Fatalf(
					"Save does not create not exist writer: %v, err: %v",
					record, err)
		}

		query = db.QueryBookBySiteIdHash("test", 1000, -1)
		defer query.Close()
		record, err = query.Scan()
		if err != nil {
			t.Fatalf("cannot query test-1000: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		book.bookRecord.WriterId = 4
		if !actualRecord.Equal(*book.bookRecord) {
			t.Fatalf("Save does not create not exist book: %v, err: %v", record, err)
		}

		query = db.QueryErrorBySiteId("test", 1000)
		defer query.Close()
		if query.Next() {
			t.Fatalf("Save create error: %v, err: %v", record, err)
		}
	})

	t.Run("success to create a completely new book with existing writer", func(t *testing.T) {
		book := NewBook("test", 1001, -1, config)
		book.SetWriter("writer-1")
		book.SetStatus(database.InProgress)
		result := book.Save(db)
		db.Commit()
		if !result {
			t.Fatalf("Save fail")
		}

		query := db.QueryWriterByName("writer-1")
		defer query.Close()
		record, _ := query.Scan()
		_, err := query.Scan()
		if err == nil || *(record.(*database.WriterRecord)) != *book.writerRecord {
			t.Fatalf(
				"Save create new writer in book: %v, writer in db: %v, err: %v",
				book.writerRecord, record, err)
		}

		query = db.QueryBookBySiteIdHash("test", 1001, -1)
		defer query.Close()
		record, err = query.Scan()
		if err != nil {
			t.Fatalf(
				"cannot query test-1001: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		
		if !book.bookRecord.Equal(*actualRecord) {
			t.Fatalf("Save does not create not exist book: %v, err: %v", actualRecord, err)
		}

		query = db.QueryErrorBySiteId("test", 100)
		defer query.Close()
		if query.Next() {
			t.Fatalf("Save create error: %v, err: %v", record, err)
		}
	})

	t.Run("success to create a new book with error", func(t *testing.T) {
		book := NewBook("test", 1002, -1, config)
		book.SetError(errors.New("test error"))
		result := book.Save(db)
		db.Commit()
		if !result {
			t.Fatalf("Save fail")
		}

		query := db.QueryBookBySiteIdHash("test", 1002, -1)
		defer query.Close()
		record, err := query.Scan()
		if err != nil {
			t.Fatalf("cannot query test-1002: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		book.bookRecord.WriterId = 0
		if !book.bookRecord.Equal(*actualRecord) {
			t.Fatalf("Save does not create not exist book: %v, err: %v", record, err)
		}

		query = db.QueryErrorBySiteId("test", 1002)
		defer query.Close()
		record, err = query.Scan()
		if err != nil || record.(*database.ErrorRecord).Site != book.errorRecord.Site ||
			record.(*database.ErrorRecord).Id != book.errorRecord.Id ||
			record.(*database.ErrorRecord).Error.Error() != book.errorRecord.Error.Error() {
				t.Fatalf(
					"Save does not create error in book: %v, error in db: %v, err: %v",
					book.errorRecord, record, err)
		}
	})

	t.Run("success to update an existing book", func(t *testing.T) {
		book := LoadBook(db, "test", 1000, -1, config)
		book.SetTitle("title-1-new")
		result := book.Save(db)
		db.Commit()
		if !result {
			t.Fatalf("Save fail")
		}

		query := db.QueryBookBySiteIdHash("test", 1000, -1)
		defer query.Close()
		record, err := query.Scan()
		if err != nil {
			t.Fatalf("cannot query test-1000: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		book.bookRecord.WriterId = 4
		if !book.bookRecord.Equal(*actualRecord) {
			t.Fatalf("Save does not create not exist book: %v, err: %v", record, err)
		}

		query = db.QueryErrorBySiteId("test", 1000)
		defer query.Close()
		if query.Next() {
			t.Fatalf( "Save create error: %v, err: %v", record, err)
		}
	})

	t.Run("success to update an existing book from error to in progress", func(t *testing.T) {
		book := LoadBook(db, "test", 1002, -1, config)
		book.SetStatus(database.InProgress)
		result := book.Save(db)
		db.Commit()
		if !result {
			t.Fatalf("Save fail")
		}

		query := db.QueryBookBySiteIdHash("test", 1002, -1)
		record, err := query.Scan()
		if err != nil {
			t.Fatalf("cannot query test-1000: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		defer query.Close()
		book.bookRecord.WriterId = 0
		if !book.bookRecord.Equal(*actualRecord) {
			t.Fatalf("Save does not create not exist book: %v, err: %v", record, err)
		}

		query = db.QueryErrorBySiteId("test", 1002)
		defer query.Close()
		if query.Next() {
			t.Fatalf("Save does not remove error: %v, err: %v", record, err)
		}
	})

	t.Run("success to update an existing book with new writer", func(t *testing.T) {
		book := LoadBook(db, "test", 1000, -1, config)
		book.SetWriter("new-writer-2")
		result := book.Save(db)
		db.Commit()
		if !result {
			t.Fatalf("Save fail")
		}

		query := db.QueryWriterByName("new-writer-2")
		defer query.Close()
		record, err := query.Scan()
		if err != nil {
			t.Fatalf("cannot query new-writer-2: %v", err)
		}
		writerRecord := record.(*database.WriterRecord)
		book.writerRecord.Id = 5
		if err != nil || !writerRecord.Equal(*book.writerRecord) {
			t.Fatalf(
				"Save does not create not exist writer: %v, err: %v",
				record, err)
		}

		query = db.QueryBookBySiteIdHash("test", 1000, -1)
		defer query.Close()
		record, err = query.Scan()
		if err != nil {
			t.Fatalf("cannot query test-1000: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		book.bookRecord.WriterId = 5
		if !book.bookRecord.Equal(*actualRecord) {
				t.Fatalf(
					"Save does not create not exist book: %v, err: %v",
					record, err)
		}

		query = db.QueryErrorBySiteId("test", 1002)
		defer query.Close()
		if query.Next() {
			t.Fatalf("Save create error: %v, err: %v", record, err)
		}
	})

	t.Run("success to create new book with existing same site and id, but different hash", func(t *testing.T) {
		book := NewBook("test", 1000, -1, config)
		book.SetWriter("new-writer-2")
		book.SetStatus(database.InProgress)
		result := book.Save(db)
		db.Commit()
		if !result {
			t.Fatalf("Save fail")
		}

		query := db.QueryBookBySiteIdHash("test", 1000, -1)
		defer query.Close()
		record, _ := query.Scan()
		_, err := query.Scan()
		if err != nil {
			t.Fatalf("cannot query test-1000: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		if !book.bookRecord.Equal(*actualRecord) {
				t.Fatalf(
					"Save does not create not exist book: %v, err: %v",
					record, book.bookRecord)
		}

		query = db.QueryErrorBySiteId("test", 1002)
		defer query.Close()
		if query.Next() {
			t.Fatalf(
				"Save create error: %v, err: %v",
				record, err)
		}
	})
}

func Test_Books_Book_validHTML(t *testing.T) {
	book := NewBook("test", 1, 500, config)

	t.Run("success", func(t *testing.T) {
		err := book.validHTML("hello")
		if err != nil {
			t.Fatalf("validate normal string as html return err: %v", err)
		}
	})

	t.Run("fail if input html is empty", func(t *testing.T) {
		err := book.validHTML("")
		if err == nil {
			t.Fatalf("valid empty string as html not return error")
		}
	})

	t.Run("fail if input html is number", func(t *testing.T) {
		err := book.validHTML("200")
		if err == nil {
			t.Fatalf("valid number string as html not return error")
		}
	})
}

func TestMain(m *testing.M) {
	initConcurrentTest()
	initBookTest()
	
	code := m.Run()

	cleanupBookTest()
	cleanupConcurrentTest()
	os.Exit(code)
}
