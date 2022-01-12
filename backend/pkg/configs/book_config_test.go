package configs

import (
	"testing"
)

func Test_Config_BookConfig_LoadBookConfigYaml(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := LoadBookConfigYaml("../../assets/test-data/book-config.yml")

		if config.BaseUrl != "https://base-url/%v" ||
			config.DownloadUrl != "https://download-url/%v" ||
			config.ChapterUrl != "https://chapter-url/%v" ||
			config.ChapterUrlPattern != "chapter-url-pattern" ||
			config.TitleRegex != "title-regex" ||
			config.WriterRegex != "writer-regex" ||
			config.TypeRegex != "type-regex" ||
			config.LastUpdateRegex != "last-update-regex" ||
			config.LastChapterRegex != "last-chapter-regex" ||
			config.ChapterUrlRegex != "chapter-url-regex" ||
			config.ChapterTitleRegex != "chapter-title-regex" ||
			config.ChapterContentRegex != "chapter-content-regex" {
				t.Fatalf("LoadBookConfigYaml return wrong book config: %v", config)
			}
	})

	t.Run("fail if location not exist", func(t *testing.T) {
		config := LoadBookConfigYaml("../../assets/test-data/not-exist-condig.yml")

		if config != nil {
			t.Fatalf("LoadBookConfigYaml for not exist file return non null value: %v", config)
		}
	})
}

func Test_Config_BookConfig_Populate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := LoadBookConfigYaml("../../assets/test-data/book-config.yml")

		populatedConfig := config.Populate(1)

		if config.BaseUrl != "https://base-url/%v" ||
			config.DownloadUrl != "https://download-url/%v" ||
			config.ChapterUrl != "https://chapter-url/%v" {
				t.Fatalf("config populate update existing config")
			}
		if populatedConfig.BaseUrl != "https://base-url/1" ||
			populatedConfig.DownloadUrl != "https://download-url/1" ||
			populatedConfig.ChapterUrl != "https://chapter-url/1" {
				t.Fatalf("config populate fail to update url")
			}
	})
}