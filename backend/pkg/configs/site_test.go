package configs

import (
	"testing"
	"os"
)

func Test_SiteConfig(t *testing.T) {
	siteConfigDirectory := os.Getenv("ASSETS_LOCATION") + "/configs"
	t.Run("func LoadSiteConfigs", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result := LoadSiteConfigs(siteConfigDirectory)
			if result == nil || len(result) != 5 {
				t.Fatalf("result: %v", result)
			}
			siteConfig := result["ck101"]
			if siteConfig.SourceKey != "ck101-desktop" ||
				siteConfig.DatabaseEngine != "sqlite3" ||
				siteConfig.DatabaseLocation != "/database/ck101.db" ||
				siteConfig.StorageDirectory != "/books/ck101/" ||
				siteConfig.BackupDirectory != "/backup/" ||
				siteConfig.SourceConfig.BaseUrl == "" {
					t.Fatalf("wrong content: %v", siteConfig)
				}
		})

		t.Run("return nil config if file not exist", func(t *testing.T) {
			result := LoadSiteConfigs(siteConfigDirectory + "abc")
			if result != nil {
				t.Fatalf("result: %v", result)
			}
		})
	})
}