package sites

import (
	"testing"
	"os"
	"io"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/pkg/flags"
)

func initUpdateTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./update_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupUpdateTest() {
	os.Remove("./update_test.db")
}

var updateConfig = configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml").SiteConfigs["test"]


func Test_Sites_Site_Update(t *testing.T) {
	updateConfig.DatabaseLocation = "./update_test.db"
	updateConfig.BookMeta.TitleRegex = "(title-.*?) "
	updateConfig.BookMeta.WriterRegex = "(writer-.*?) "
	updateConfig.BookMeta.TypeRegex = "(type-.*?) "
	updateConfig.BookMeta.LastUpdateRegex = " (last-update-.*?) "
	updateConfig.BookMeta.LastChapterRegex = "(last-chapter-.*?)$"
	site := NewSite("test", updateConfig)
	site.OpenDatabase()
	defer site.database.Close()

	server := mock.UpdateServer()
	defer server.Close()

	t.Run("func updateBook", func(t *testing.T) {
		site.config.BookMeta.CONST_SLEEP = 0

		book := books.LoadBook(site.database, "test", 1, 100, site.config.BookMeta)
		book.SetTitle("title-regex")
		book.SetWriter("writer-regex")
		book.SetType("type-regex")
		book.SetUpdateDate("last-update-regex")
		book.SetUpdateChapter("last-chapter-regex")
		book.Save(site.database)
		site.database.Commit()

		t.Run("do nothing if books does not get updated", func(t *testing.T) {
			site.config.BookMeta.BaseUrl = server.URL + "/success/%v"

			rows := site.database.QueryBookBySiteIdHash("test", 1, 100)
			record, err := rows.Scan()
			if err != nil {
				t.Fatalf("%v-%v-%v not exist in database", "test", 1, 100)
			}
			rows.Close()
			err = site.updateBook(record.(*database.BookRecord))
			site.database.Commit()
			if err != nil {
				t.Fatalf("site updateBook return error: %v", err)
			}

			rows = site.database.QueryBookBySiteIdHash("test", 1, 100)
			record, err = rows.Scan()
			if err != nil || rows.Next() {
				t.Fatalf("%v-%v-%v not exist in database", "test", 1, 100)
			}
			rows.Close()
			bookRecord := record.(*database.BookRecord)
			if bookRecord.Title != "title-regex" || bookRecord.WriterId != 4 ||
				bookRecord.Type != "type-regex" || bookRecord.UpdateDate != "last-update-regex" ||
				bookRecord.UpdateChapter != "last-chapter-regex" {
					t.Fatalf("bookRecord had been modified: %v", bookRecord)
			}
		})

		t.Run("success save the updated book to database", func(t *testing.T) {
			site.config.BookMeta.BaseUrl = server.URL + "/success/%v"
			book.SetUpdateChapter("hello")
			book.Save(site.database)
			site.database.Commit()
			rows := site.database.QueryBookBySiteIdHash("test", 1, -1)
			record, err := rows.Scan()
			if err != nil {
				t.Fatalf("%v-%v-%v not exist in database", "test", 1, 100)
			}
			rows.Close()
			err = site.updateBook(record.(*database.BookRecord))
			site.database.Commit()
			if err != nil {
				t.Fatalf("site updateBook return error: %v", err)
			}

			rows = site.database.QueryBookBySiteIdHash("test", 1, 100)
			record, err = rows.Scan()
			if err != nil || rows.Next() {
				t.Fatalf("%v-%v-%v not exist in database", "test", 1, 100)
			}
			rows.Close()
			bookRecord := record.(*database.BookRecord)
			if bookRecord.Title != "title-regex" || bookRecord.WriterId != 4 ||
				bookRecord.Type != "type-regex" || bookRecord.UpdateDate != "last-update-regex" ||
				bookRecord.UpdateChapter != "last-chapter-regex" {
					t.Fatalf("bookRecord had not been modified: %v", bookRecord)
			}
		})

		t.Run("return error if failed to load the books writer", func(t *testing.T) {
			site.config.BookMeta.BaseUrl = server.URL + "/partial_fail"
			record := &database.BookRecord{
				Site: "not-exist",
				Id: 100,
				HashCode: 100,
				WriterId: 100,
			}
			err := site.updateBook(record)
			site.database.Commit()
			if err == nil {
				t.Fatalf("site updateBook not return error when record not exist: %v", err)
			}
		})
	})

	t.Run("func update", func(t *testing.T) {
		site.config.BookMeta.BaseUrl = server.URL + "/success/%v"

		t.Run("loop add the books in database (ignore error)", func(t *testing.T) {
			summary := site.database.Summary(site.Name)
			if summary.BookCount != 6 || summary.ErrorCount != 3 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 5 ||
				summary.MaxBookId != 5 || summary.LatestSuccessId != 3 ||
				summary.StatusCount[database.Error] != 3 ||
				summary.StatusCount[database.InProgress] != 1 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}

			err := site.update()
			if err != nil {
				t.Fatalf("site update() return error: %v", err)
			}
			// downloaded book will get a new record created
			book := books.LoadBook(site.database, "test", 3, 200, site.config.BookMeta)
			if book == nil {
				t.Fatalf("cannot query %v-%v-%v, it was removed", "test", 3, 200)
			}
			book = books.LoadBook(site.database, "test", 3, -1, site.config.BookMeta)
			if book == nil {
				t.Fatalf("cannot query %v-%v", "test", 3)
			}
			_, _, hashCode := book.GetInfo()
			if hashCode == 200 || book.GetTitle() != "title-regex" {
				t.Fatalf("site update does not create new book for already download content: %v", book.GetTitle())
			}

			// error book will not be looped
			book = books.LoadBook(site.database, "test", 2, -1, site.config.BookMeta)
			if book == nil || book.GetStatus() != database.Error {
				t.Fatalf("cannot query %v-%v", "test", 3)
			}

			summary = site.database.Summary(site.Name)
			if summary.BookCount != 7 || summary.ErrorCount != 3 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 5 ||
				summary.MaxBookId != 5 || summary.LatestSuccessId != 3 ||
				summary.StatusCount[database.Error] != 3 ||
				summary.StatusCount[database.InProgress] != 2 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("book update generate wrong summary: %v", summary)
				}
		})
	})
	
	t.Run("func Update", func(t *testing.T) {
		site.OpenDatabase()
		defer site.database.Close()

		t.Run("success for full site update", func(t *testing.T) {
			site.config.BookMeta.BaseUrl = server.URL + "/success/%v"

			f := &flags.Flags{}
			err := site.Update(f)
			if err != nil {
				t.Fatalf("site Update return error for specific book - error: %v", err)
			}

			summary := site.database.Summary(site.Name)
			if summary.BookCount != 7 || summary.ErrorCount != 3 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 5 ||
				summary.MaxBookId != 5 || summary.LatestSuccessId != 3 ||
				summary.StatusCount[database.Error] != 3 ||
				summary.StatusCount[database.InProgress] != 2 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("book update generate wrong summary: %v", summary)
				}
		})

		t.Run("success for specific book", func(t *testing.T) {
			site.config.BookMeta.LastUpdateRegex = " (\\d+) "
			site.config.BookMeta.LastChapterRegex = "(chapter-.*?)$"
			site.config.BookMeta.BaseUrl = server.URL + "/specific_success/%v"
			flagSite, flagId, flagHash := "test", 1, "-1"
			f := &flags.Flags{
				Site: &flagSite,
				Id: &flagId,
				HashCode: &flagHash,
			}

			err := site.Update(f)
			if err != nil {
				t.Fatalf("site Update return error for specific book - error: %v", err)
			}
			book := books.LoadBook(site.database, "test", 1, -1, site.config.BookMeta)
			summary := site.database.Summary(site.Name)
			if book == nil || book.GetUpdateDate() != "104" ||
				book.GetUpdateChapter() != "chapter-1" || summary.BookCount != 7 {
					t.Fatalf("wrong result: book count: %v, book: %v", summary.BookCount, book.GetUpdateChapter())
			}
		})

		t.Run("fail for invalid arguments", func(t *testing.T) {
			flagId := 123
			f := &flags.Flags{
				Id: &flagId,
			}

			err := site.Update(f)
			if err == nil {
				t.Fatalf("site Update not return error for invalid arguments")
			}
		})

		t.Run("skip if arguments provide site name but not matched", func(t *testing.T) {
			flagSite := "others"
			f := &flags.Flags{
				Site: &flagSite,
			}

			err := site.Update(f)
			if err != nil {
				t.Fatalf("site Update return error for not matching site name- error: %v", err)
			}
		})
	})
}