package sites

import (
	"testing"
	"io"
	"os"
	"errors"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/pkg/flags"
)

func initFixTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./fix_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupFixTest() {
	os.Remove("./fix_test.db")
}

var fixConfig = configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml").SiteConfigs["test"]

func Test_Sites_Site_Fix(t *testing.T) {
	fixConfig.DatabaseLocation = "./fix_test.db"
	site := NewSite("test", fixConfig)
	site.OpenDatabase()
	defer site.CloseDatabase()

	var operation SiteOperation
	operation = Fix

	server := mock.UpdateServer()
	defer server.Close()
	t.Run("func addMissingRecords", func(t *testing.T) {
		t.Run("success with adding error record", func(t *testing.T) {
			book := books.NewBook("test", 7, 300, site.config.BookMeta)
			book.SetError(errors.New("test"))
			book.Save(site.database)
			site.CommitDatabase()
			summary := site.database.Summary(site.Name)
			if summary.BookCount != 7 || summary.ErrorCount != 4 ||
				summary.WriterCount != 3 || summary.UniqueBookCount != 6 ||
				summary.MaxBookId != 7 || summary.LatestSuccessId != 3 ||
				summary.StatusCount[database.Error] != 4 ||
				summary.StatusCount[database.InProgress] != 1 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
			site.config.BookMeta.BaseUrl = server.URL + "/partial_fail/%v"
			err := site.addMissingRecords()
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.addMissingRecords return error: %v", err)
			}
			summary = site.database.Summary(site.Name)
			if summary.BookCount != 8 || summary.ErrorCount != 5 ||
				summary.WriterCount != 3 || summary.UniqueBookCount != 7 ||
				summary.MaxBookId != 7 || summary.LatestSuccessId != 3 ||
				summary.StatusCount[database.Error] != 5 ||
				summary.StatusCount[database.InProgress] != 1 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
		})

		t.Run("success with adding in progress record", func(t *testing.T) {
			book := books.NewBook("test", 9, 300, site.config.BookMeta)
			book.SetError(errors.New("test"))
			book.Save(site.database)
			site.CommitDatabase()
			summary := site.database.Summary(site.Name)
			if summary.BookCount != 9 || summary.ErrorCount != 6 ||
				summary.WriterCount != 3 || summary.UniqueBookCount != 8 ||
				summary.MaxBookId != 9 || summary.LatestSuccessId != 3 ||
				summary.StatusCount[database.Error] != 6 ||
				summary.StatusCount[database.InProgress] != 1 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
			site.config.BookMeta.BaseUrl = server.URL + "/success/%v"
			err := site.addMissingRecords()
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.addMissingRecords return error: %v", err)
			}
			summary = site.database.Summary(site.Name)
			if summary.BookCount != 10 || summary.ErrorCount != 6 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 9 ||
				summary.MaxBookId != 9 || summary.LatestSuccessId != 8 ||
				summary.StatusCount[database.Error] != 6 ||
				summary.StatusCount[database.InProgress] != 2 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
		})

		t.Run("do nothing if nothing missed", func(t *testing.T) {
			err := site.addMissingRecords()
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.addMissingRecords return error: %v", err)
			}
			summary := site.database.Summary(site.Name)
			if summary.BookCount != 10 || summary.ErrorCount != 6 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 9 ||
				summary.MaxBookId != 9 || summary.LatestSuccessId != 8 ||
				summary.StatusCount[database.Error] != 6 ||
				summary.StatusCount[database.InProgress] != 2 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
		})
	})

	t.Run("func updateBooksByStorage", func(t *testing.T) {
		site.config.BookMeta.BaseUrl = server.URL + ""
		site.updateBooksByStorage()
		site.CommitDatabase()

		t.Run("success update book with storage", func(t *testing.T) {
			book := books.LoadBook(site.database, "test", 1, 100, site.config.BookMeta)
			if book.GetStatus() != database.Download {
				t.Fatalf("site.updateBooksByStorage does not update book with storage: %v", book.GetStatus())
			}
		})

		t.Run("success update book without storage", func(t *testing.T) {
			book := books.LoadBook(site.database, "test", 3, 102, site.config.BookMeta)
			if book.GetStatus() != database.End {
				t.Fatalf("site.updateBooksByStorage does not update book without storage: %v", book.GetStatus())
			}
		})
	})

	t.Run("func Fix", func(t *testing.T) {
		book := books.LoadBook(site.database, "test", 1, 100, site.config.BookMeta)
		book.SetStatus(database.InProgress)
		book.Save(site.database)
		site.CommitDatabase()

		t.Run("success for full site", func(t *testing.T) {
			err := operation(site, &flags.Flags{})
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site Fix return error for full site - error: %v", err)
			}
			book := books.LoadBook(site.database, "test", 1, 100, site.config.BookMeta)
			if book.GetStatus() != database.Download {
				t.Fatalf("site.Fix does not fix the record")
			}
		})

		t.Run("fail for invalid arguements", func(t *testing.T) {
			flagId := 123

			err := operation(site, &flags.Flags{ Id: &flagId })
			if err == nil {
				t.Fatalf("site Fix not return error for invalid arguments")
			}
		})

		t.Run("skip if arguments provide mismatch site name", func(t *testing.T) {
			flagSite := "others"

			err := operation(site, &flags.Flags{ Site: &flagSite })
			if err != nil {
				t.Fatalf("site Fix return error for not matching site name- error: %v", err)
			}
		})
	})
}