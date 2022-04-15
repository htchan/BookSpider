package sites

import (
	"testing"
	"io"
	"os"
	// "errors"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/pkg/flags"
)

func initDownloadTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./download_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupDownloadTest() {
	os.Remove("./download_test.db")
}

var downloadConfig = configs.LoadSiteConfigs(os.Getenv("ASSETS_LOCATION") + "/test-data/configs")["test"]

func Test_Sites_Site_Download(t *testing.T) {
	downloadConfig.DatabaseLocation = "./download_test.db"
	site := NewSite("test", downloadConfig)
	site.OpenDatabase()
	defer site.CloseDatabase()
	var operation SiteOperation
	operation = Download

	server := mock.DownloadServer()
	defer server.Close()
	t.Run("func Download", func(t *testing.T) {
		site.config.SourceConfig.DownloadUrl = server.URL + "/content/success/%v"
		site.config.SourceConfig.ChapterUrl = server.URL + "/chapter/success%v"
		t.Run("success for specific book", func(t *testing.T) {
			flagSite, flagId, flagThreads := "test", 3, 0
			f := &flags.Flags{
				Site: &flagSite,
				Id: &flagId,
				MaxThreads: &flagThreads,
			}
			err := operation(site, f)
			utils.CheckError(site.CommitDatabase())
			if err != nil {
				t.Fatalf("site Download return error for specific book: %v", err)
			}
			book := books.LoadBook(site.database, "test", 3, -1, site.config.SourceConfig)
			if book == nil || book.GetStatus() != database.Download || !book.HasContent(site.config.StorageDirectory) {
				t.Fatalf("site Download does not turn book status to download and save content: %v", book)
			}
		})

		t.Run("success for all site", func(t *testing.T) {
			book := books.LoadBook(site.database, site.Name, 2, -1, site.config.SourceConfig)
			book.SetStatus(database.End)
			book.Save(site.database)
			site.CommitDatabase()

			err := operation(site, &flags.Flags{})
			utils.CheckError(site.CommitDatabase())
			if err != nil {
				t.Fatalf("site Download return error for all site: %v", err)
			}
			book = books.LoadBook(site.database, "test", 2, -1, site.config.SourceConfig)
			if book == nil || book.GetStatus() != database.Download || !book.HasContent(site.config.StorageDirectory) {
				t.Fatalf("site Download does not turn book status to download and save content: %v", book)
			}
		})

		t.Run("fail if only id is provided", func(t *testing.T) {
			flagId := 123

			err := operation(site, &flags.Flags{ Id: &flagId })
			if err == nil {
				t.Fatalf("site Download not return error for invalid arguments")
			}
		})

		t.Run("skip if arguments provide mismatch site name", func(t *testing.T) {
			flagSite := "others"

			err := operation(site, &flags.Flags{ Site: &flagSite })
			if err != nil {
				t.Fatalf("site Download return error for not matching site name- error: %v", err)
			}
		})

		t.Run("skip if specific book not exist", func(t *testing.T) {
			flagSite, flagId, flagThreads := "test", 999, 999
			f := &flags.Flags{
				Site: &flagSite,
				Id: &flagId,
				MaxThreads: &flagThreads,
			}
			err := operation(site, f)
			if err == nil {
				t.Fatalf("site Download not return error for not exist book")
			}
		})
	})
}