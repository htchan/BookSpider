package configs

import (
	"testing"
)

func Test_Config_Config_LoadConfigYaml(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := LoadConfigYaml("../../assets/test-data/config.yml")

		if config.Backend.StageFile != "/log/stage.txt" ||
			config.Backend.LogFile != "/log/controller.log" ||
			len(config.Backend.Api) != 8 {
				t.Fatalf("Load wrong backend config %v", config.Backend)
			}
		
		if site, ok := config.SiteConfigs["test"]; !ok ||
			site.DecoderString != "big5" ||
			site.BookConfigLocation != "../../assets/test-data/book-config.yml" ||
			site.DatabaseLocation != "./test.db" ||
			site.StorageDirectory != "../../assets/test-data/storage/" ||
			site.ThreadsCount != 1000 || site.ConstSleep != 1000 ||
			site.MAX_THREAD_COUNT != 1000 ||
			site.Decoder == nil || site.BookMeta == nil {
				t.Fatalf("Load wrong site config %v", site)
			}
		
		if bookConfig := config.SiteConfigs["test"].BookMeta;
			bookConfig.CONST_SLEEP != 1000 ||
			bookConfig.Decoder == nil ||
			bookConfig.StorageDirectory != "../../assets/test-data/storage/" {
				t.Fatalf("Load wrong book config %v", bookConfig)
			}
	})
}
