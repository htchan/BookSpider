package sites

import (
	"testing"
	"io"
	"os"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/ApiParser"
)

func initSiteTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupSiteTest() {
	os.Remove("./test.db")
}

var siteConfig = configs.LoadSiteConfigs(os.Getenv("ASSETS_LOCATION") + "/test-data/configs")["test"]

func Test_Sites_Constructor_NewSite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		site := NewSite("test", siteConfig)
		if site.Name != "test" || site.database != nil || site.config != siteConfig {
			t.Fatalf("site init with wrong data %v", site)
		}
		if !site.semaphore.TryAcquire(1000) || site.semaphore.TryAcquire(1) {
				t.Fatalf("site init with wrong threads count")
		}
	})
}

func Test_Sites_Site_OpenDatabase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		site := NewSite("test", siteConfig)
		err := site.OpenDatabase()
		defer site.CloseDatabase()
		if err != nil {
			t.Fatalf("site OpenDatabase return error: %v", err)
		}
	})

	t.Run("fail because of unknown engine", func(t *testing.T) {
		tempConfig := *siteConfig
		tempConfig.DatabaseEngine = "unknown"
		site := NewSite("test", &tempConfig)
		err := site.OpenDatabase()
		if err == nil {
			t.Fatalf("site OpenDatabase not return error")
		}
	})
}

func Test_Sites_Site_Map(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		site := NewSite("test", siteConfig)
		site.OpenDatabase()
		defer site.CloseDatabase()
		siteMap := site.Map()
		if siteMap["name"] != "test" || siteMap["bookCount"] != 6 ||
			siteMap["uniqueBookCount"] != 5 || siteMap["writerCount"] != 3 || 
			siteMap["errorCount"] != 3 ||
			siteMap["maxBookId"]!= 5 || siteMap["latestSuccessBookId"]!= 3 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.Error]!= 3 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.InProgress]!= 1 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.End]!= 1 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.Download]!= 1 {
				t.Fatalf("site.Map return wrong data: %v", siteMap)
		}
	})

	t.Run("fail", func(t *testing.T) {
		tempConfig := *siteConfig
		tempConfig.DatabaseEngine = "unknown"
		site := NewSite("test", &tempConfig)
		site.OpenDatabase()
		defer site.CloseDatabase()
		siteMap := site.Map()
		if siteMap["name"] != "test" || siteMap["bookCount"] != 0 ||
			siteMap["uniqueBookCount"] != 0 || siteMap["writerCount"] != 0 || 
			siteMap["errorCount"] != 0 ||
			siteMap["maxBookId"]!= 0 || siteMap["latestSuccessBookId"]!= 0 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.Error]!= 0 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.InProgress]!= 0 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.End]!= 0 ||
			siteMap["statusCount"].(map[database.StatusCode]int)[database.Download]!= 0 {
				t.Fatalf("site.Map return wrong data: %v", siteMap)
		}
	})
}

func TestMain(m *testing.M) {
	ApiParser.Setup(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser")
	
	initBackupTest()
	initCheckTest()
	initDownloadTest()
	initExploreTest()
	initFixTest()
	initQueryTest()
	initSiteTest()
	initUpdateTest()

	code := m.Run()

	cleanupUpdateTest()
	cleanupSiteTest()
	cleanupQueryTest()
	cleanupFixTest()
	cleanupExploreTest()
	cleanupDownloadTest()
	cleanupCheckTest()
	cleanupBackupTest()
	os.Exit(code)
}