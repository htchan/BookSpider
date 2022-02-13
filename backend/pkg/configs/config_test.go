package configs

import (
	"testing"
	"os"
)

func Test_Config_Config(t *testing.T) {
	t.Run("func LoadConfigYaml", func(t *testing.T) {
		t.Run("has correct data", func(t *testing.T) {
			t.Parallel()
			config := LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml")

			if config.Backend.StageFile != "/log/stage.txt" ||
				config.Backend.LogFile != "/log/controller.log" ||
				len(config.Backend.Api) != 8 {
					t.Fatalf("Load wrong backend config %v", config.Backend)
				}
			
			if site, ok := config.SiteConfigs["test"]; !ok ||
				site.DecoderString != "big5" ||
				site.BookConfigLocation != "/test-data/book-config.yml" ||
				site.DatabaseEngine != "sqlite3" ||
				site.DatabaseLocation != "./test.db" ||
				site.StorageDirectory != "/test-data/storage/" ||
				site.DownloadThreadsCount != 5 ||
				site.ThreadsCount != 1000 || site.ConstSleep != 1000 ||
				site.MaxExploreError != 1000 ||
				site.Decoder == nil || site.BookMeta == nil {
					t.Fatalf("Load wrong site config %v", site.DownloadThreadsCount)
				}
			
			if bookConfig := config.SiteConfigs["test"].BookMeta;
				bookConfig.CONST_SLEEP != 1000 ||
				bookConfig.Decoder == nil ||
				bookConfig.StorageDirectory != "/test-data/storage/" {
					t.Fatalf("Load wrong book config %v", bookConfig)
				}
		})
	})
}
