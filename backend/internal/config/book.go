package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"

	"github.com/htchan/BookSpider/internal/utils"
)

type URLConfig struct {
	Base string `yaml:"base"`
	Download string `yaml:"download"`
	ChapterPrefix string `yaml:"chapterPrefix"`
}

type BookConfig struct {
	URL URLConfig `yaml:"url"`
	MaxChaptersError int `yaml:"maxChapterError"`
	Storage string
	SourceKey string
}

type RequestFunc func(string) (string, error)

func LoadBookConfigs(configDirectory string) (config map[string]BookConfig) {
	defer utils.Recover(func() { config = nil })

	data, err := os.ReadFile(filepath.Join(configDirectory, "source_config.yaml"))
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &config))
	for key, value := range config {
		value.SourceKey = key
		config[key] = value
	}
	return
}