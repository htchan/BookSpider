package configs

import (
	"testing"
	"os"
)

func Test_SourceConfig(t *testing.T) {
	sourceConfigDirectory := os.Getenv("ASSETS_LOCATION") + "/configs"
	t.Run("func LoadSourceConfigs", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			result := LoadSourceConfigs(sourceConfigDirectory)
			if result == nil || len(result) != 6 {
				t.Errorf("result: %v", result)
			}

			sourceConfig := result["ck101-desktop"]
			if sourceConfig.BaseUrl != "https://www.ck101.org/book/%v.html" ||
			sourceConfig.DownloadUrl != "https://www.ck101.org/0/%v/" ||
			sourceConfig.ChapterUrl != "https://www.ck101.org%v" ||
			sourceConfig.ChapterUrlPattern != "/.*?\\.html" ||
			sourceConfig.DecoderString != "big5" ||
			sourceConfig.Threads != 1000 ||
			sourceConfig.DownloadBookThreads != 10 ||
			sourceConfig.ConstSleep != 1000 ||
			sourceConfig.MaxExploreError != 1000 ||
			sourceConfig.UseRequestInterval != false ||
			sourceConfig.Decoder == nil {
				t.Errorf("wrong content: %v", sourceConfig)
			}
		})

		t.Run("return nil config if file not exist", func(t *testing.T) {
			result := LoadSourceConfigs(sourceConfigDirectory + "abc")
			if result != nil {
				t.Errorf("result: %v", result)
			}
		})
	})

	t.Run("func Populate", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			config := LoadSourceConfigs(sourceConfigDirectory)
	
			sourceConfig := config["ck101-desktop"]

			result := sourceConfig.Populate(1)
			
			if result.BaseUrl != "https://www.ck101.org/book/1.html" ||
				result.DownloadUrl != "https://www.ck101.org/0/1/" {
					t.Errorf("wrong content: %v", result)
				}
		})
	})
}