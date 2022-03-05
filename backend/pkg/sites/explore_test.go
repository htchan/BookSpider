package sites

import (
	"testing"
	"io"
	"os"
	"golang.org/x/sync/semaphore"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/pkg/flags"
)

func initExploreTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./explore_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupExploreTest() {
	os.Remove("./explore_test.db")
}

var exploreConfig = configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml").SiteConfigs["test"]

func Test_Sites_Site_Explore(t *testing.T) {
	exploreConfig.DatabaseLocation = "./explore_test.db"
	exploreConfig.BookMeta.TitleRegex = "(title-.*?) "
	exploreConfig.BookMeta.WriterRegex = "(writer-.*?) "
	exploreConfig.BookMeta.TypeRegex = "(type-.*?) "
	exploreConfig.BookMeta.LastUpdateRegex = " (last-update-.*?) "
	exploreConfig.BookMeta.LastChapterRegex = "(last-chapter-.*?)$"
	site := NewSite("test", exploreConfig)
	site.OpenDatabase()
	defer site.CloseDatabase()

	server := mock.UpdateServer()
	defer server.Close()

	var operation SiteOperation
	operation = Explore
	
	t.Run("func exploreOldBook", func(t *testing.T) {
		t.Run("update if book exist in db and updated in web", func(t *testing.T) {
			count := 1
			site.config.BookMeta.BaseUrl = server.URL + "/success/%v"
			err := site.exploreOldBook(2, &count)
			site.CommitDatabase()
			if count != 0 {
				t.Fatalf("site.exploreOldBook not reset count: %v", count)
			}
			if err != nil {
				t.Fatalf("site.exploreOldBook return error: %v", err)
			}
			book := books.LoadBook(site.database, "test", 2, 101, site.config.BookMeta)
			if book.GetStatus() == database.Error || book.GetTitle() != "title-regex" || 
				book.GetWriter() != "writer-regex" || book.GetType() != "type-regex" ||
				book.GetUpdateDate() != "last-update-regex" || book.GetError() != nil {
				t.Fatalf("site.updateOldBook fail update book: %v, err: %v", book, err)
			}
		})

		t.Run("add count if book exist in db but not updated in web", func(t *testing.T) {
			count := 0
			site.config.BookMeta.BaseUrl = server.URL + "/partial_fail/%v"
			err := site.exploreOldBook(4, &count)
			site.CommitDatabase()
			if count != 1 {
				t.Fatalf("site.exploreOldBook not update count: %v", count)
			}
			if err != nil {
				t.Fatalf("site.exploreOldBook return error: %v", err)
			}
			book := books.LoadBook(site.database, "test", 4, -1, site.config.BookMeta)
			if book.GetStatus() != database.Error || book.GetTitle() != "" || 
				book.GetWriter() != "" || book.GetType() != "" ||
				book.GetUpdateDate() != "" || book.GetError() == nil {
				t.Fatalf("site.updateOldBook success update book: %v, err: %v", book, err)
			}
		})

		t.Run("return error if book not found", func(t *testing.T) {
			count := 0
			err := site.exploreOldBook(999, &count)
			if err == nil || err.Error() != "load book test-999 fail"{
				t.Fatalf("site.exploreOldBook return %v", err)
			}
		})

		t.Run("return error if book is not in error status", func(t *testing.T) {
			count := 0
			err := site.exploreOldBook(1, &count)
			if err == nil || err.Error() != "load book test-1 return status 1"{
				t.Fatalf("site.exploreOldBook return %v", err)
			}
		})
	})

	t.Run("func exploreNewBook", func(t *testing.T) {
		t.Run("success if book not exist in db and exist in web", func(t *testing.T) {
			count := 1
			site.config.BookMeta.BaseUrl = server.URL + "/success/%v"
			err := site.exploreNewBook(6, &count)
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.updateNewBook return error: %v", err)
			}
			if count != 0 {
				t.Fatalf("site.updateNewBook not reset count: %v", count)
			}
			book := books.LoadBook(site.database, "test", 6, -1, site.config.BookMeta)
			if book.GetStatus() == database.Error || book.GetTitle() != "title-regex" ||
				book.GetWriter() != "writer-regex" || book.GetError() != nil {
					t.Fatalf("sites.updateNewBook does not create duplicated book for existing book: %v", book.GetWriter())
				}
		})

		t.Run("add count and save book if book not exist in db and web", func(t *testing.T) {
			count := 0
			site.config.BookMeta.BaseUrl = server.URL + "/partial_fail/%v"
			err := site.exploreNewBook(7, &count)
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.updateNewBook return error: %v", err)
			}
			if count != 1 {
				t.Fatalf("site.updateNewBook not update count: %v", count)
			}
			book := books.LoadBook(site.database, "test", 7, -1, site.config.BookMeta)
			if book.GetStatus() != database.Error || book.GetTitle() != "" ||
				book.GetWriter() != "" || book.GetError() == nil {
					t.Fatalf("sites.updateNewBook does not create duplicated book for existing book: %v", book.GetWriter())
				}
		})

		t.Run("create duplicated book if book exist in db", func(t *testing.T) {
			count := 0
			site.config.BookMeta.BaseUrl = server.URL + "/partial_fail/%v"
			err := site.exploreNewBook(1, &count)
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.updateNewBook return error: %v", err)
			}
			if count != 1 {
				t.Fatalf("site.updateNewBook not update count: %v", count)
			}
			book := books.LoadBook(site.database, "test", 1, -1, site.config.BookMeta)
			if _, _, hashCode := book.GetInfo(); hashCode == 100 ||
				book.GetStatus() != database.Error || book.GetTitle() != "" ||
				book.GetWriter() != "" || book.GetError() == nil {
					t.Fatalf("sites.updateNewBook does not create duplicated book for existing book: %v", book)
				}
		})
	})

	t.Run("func explore", func(t *testing.T) {
		site.config.BookMeta.BaseUrl = server.URL + "/success/%v"
		site.config.BookMeta.CONST_SLEEP = 0
		site.config.ThreadsCount = 2
		site.semaphore = semaphore.NewWeighted(int64(site.config.ThreadsCount))
		site.config.MaxExploreError = 3
		t.Run("success for adding new books", func(t *testing.T) {
			summary := site.database.Summary(site.Name)
			if summary.BookCount != 9 || summary.ErrorCount != 4 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 7 ||
				summary.MaxBookId != 7 || summary.LatestSuccessId != 6 ||
				summary.StatusCount[database.Error] != 4 ||
				summary.StatusCount[database.InProgress] != 3 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
			err := site.explore()
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.explore return error: %v", err)
			}
			summary = site.database.Summary(site.Name)
			if summary.BookCount != 15 || summary.ErrorCount != 8 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 13 ||
				summary.MaxBookId != 13 || summary.LatestSuccessId != 8 ||
				summary.StatusCount[database.Error] != 8 ||
				summary.StatusCount[database.InProgress] != 5 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
		})
		t.Run("not add new book if it reach limit in exploring existing books", func(t *testing.T) {
			summary := site.database.Summary(site.Name)
			if summary.BookCount != 15 || summary.ErrorCount != 8 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 13 ||
				summary.MaxBookId != 13 || summary.LatestSuccessId != 8 ||
				summary.StatusCount[database.Error] != 8 ||
				summary.StatusCount[database.InProgress] != 5 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
			err := site.explore()
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site.explore return error: %v", err)
			}
			summary = site.database.Summary(site.Name)
			if summary.BookCount != 15 || summary.ErrorCount != 8 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 13 ||
				summary.MaxBookId != 13 || summary.LatestSuccessId != 8 ||
				summary.StatusCount[database.Error] != 8 ||
				summary.StatusCount[database.InProgress] != 5 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
		})
	})

	t.Run("func Explore", func(t *testing.T) {
		t.Run("success for full site", func(t *testing.T) {
			site.config.BookMeta.BaseUrl = server.URL + "/partial_fail/%v"

			err := operation(site, &flags.Flags{})
			if err != nil {
				t.Fatalf("site Explore return error for full site - error: %v", err)
			}

			summary := site.database.Summary(site.Name)
			if summary.BookCount != 15 || summary.ErrorCount != 8 ||
				summary.WriterCount != 4 || summary.UniqueBookCount != 13 ||
				summary.MaxBookId != 13 || summary.LatestSuccessId != 8 ||
				summary.StatusCount[database.Error] != 8 ||
				summary.StatusCount[database.InProgress] != 5 ||
				summary.StatusCount[database.End] != 1 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
		})

		t.Run("fail for invalid arguements", func(t *testing.T) {
			flagId := 123

			err := operation(site, &flags.Flags{ Id: &flagId })
			if err == nil {
				t.Fatalf("site Explore not return error for invalid arguments")
			}
		})

		t.Run("skip if arguments provide mismatch site name", func(t *testing.T) {
			flagSite := "others"

			err := operation(site, &flags.Flags{ Site: &flagSite })
			if err != nil {
				t.Fatalf("site Explore return error for not matching site name- error: %v", err)
			}
		})
	})
}