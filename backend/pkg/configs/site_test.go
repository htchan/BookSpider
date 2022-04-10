package configs

import (
	"testing"
	"os"
)

func Test_SiteConfig(t *testing.T) {
	siteConfigLocation := os.Getenv("ASSETS_LOCATION") + "/configs/site_config.yaml"
	t.Run("func LoadSiteConfigs", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result := LoadSiteConfigs(siteConfigLocation)
			if result == nil || len(result) != 5 {
				t.Fatalf("result: %v", result)
			}
			siteConfig := result["ck101"]
			if siteConfig.SourceKey != "ck101-desktop" ||
				siteConfig.DatabaseEngine != "sqlite3" ||
				siteConfig.DatabaseLocation != "/database/ck101.db" ||
				siteConfig.StorageDirectory != "/books/ck101/" ||
				siteConfig.BackupDirectory != "/backup/" {
					t.Fatalf("wrong content: %v", siteConfig)
				}
		})

		t.Run("return nil config if file not exist", func(t *testing.T) {
			result := LoadSiteConfigs(siteConfigLocation + "abc")
			if result != nil {
				t.Fatalf("result: %v", result)
			}
		})
	})
}