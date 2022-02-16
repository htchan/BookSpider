package sites

import (
	"testing"
	"io"
	"os"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/pkg/flags"
	"fmt"
)

func initCheckTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./check_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupCheckTest() {
	os.Remove("./check_test.db")
}

var checkConfig = configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml").SiteConfigs["test"]

func Test_Sites_Site_Check(t *testing.T) {
	checkConfig.DatabaseLocation = "./check_test.db"
	site := NewSite("test", checkConfig)
	site.OpenDatabase()
	defer site.CloseDatabase()

	t.Run("func Check", func(t *testing.T) {
		t.Run("success for full site", func(t *testing.T) {
			book := books.LoadBook(site.database, "test", 1, 100, site.config.BookMeta)
			book.SetUpdateChapter("后记abcdef")
			book.Save(site.database)
			site.CommitDatabase()

			f := &flags.Flags{}
			err := site.Check(f)
			fmt.Println("database", site.database)
			site.CommitDatabase()
			if err != nil {
				t.Fatalf("site Check return error for full site - error: %v", err)
			}
			book = books.LoadBook(site.database, "test", 1, 100, site.config.BookMeta)
			if book.GetStatus() != database.End {
				t.Fatalf("site.Check does not update the record status to end")
			}
			summary := site.database.Summary(site.Name)
			if summary.BookCount != 6 || summary.ErrorCount != 3 ||
				summary.WriterCount != 3 || summary.UniqueBookCount != 5 ||
				summary.MaxBookId != 5 || summary.LatestSuccessId != 3 ||
				summary.StatusCount[database.Error] != 3 ||
				summary.StatusCount[database.InProgress] != 0 ||
				summary.StatusCount[database.End] != 2 ||
				summary.StatusCount[database.Download] != 1 {
					t.Fatalf("before book update generate wrong summary: %v", summary)
				}
		})

		t.Run("fail for invalid arguements", func(t *testing.T) {
			flagId := 123
			f := &flags.Flags{
				Id: &flagId,
			}

			err := site.Check(f)
			if err == nil {
				t.Fatalf("site Check not return error for invalid arguments")
			}
		})

		t.Run("skip if arguments provide mismatch site name", func(t *testing.T) {
			flagSite := "others"
			f := &flags.Flags{
				Site: &flagSite,
			}

			err := site.Check(f)
			if err != nil {
				t.Fatalf("site Check return error for not matching site name- error: %v", err)
			}
		})
	})
}