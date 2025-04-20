package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type URLConfig struct {
	Base          string `yaml:"base"`
	Download      string `yaml:"download"`
	ChapterPrefix string `yaml:"chapterPrefix"`
}

type BookConfig struct {
	URLConfig        URLConfig `yaml:"url"`
	MaxChaptersError int       `yaml:"maxChapterError"`
	UpdateDateLayout string    `yaml:"updateDateLayout"`
	UnwantContent    []string  `yaml:"unwantContent"`
}

func LoadBookConfigs(configDirectory string) (map[string]*BookConfig, error) {
	data, err := os.ReadFile(filepath.Join(configDirectory, "book_configs.yaml"))
	if err != nil {
		return nil, err
	}

	var bookConfigs map[string]*BookConfig
	err = yaml.Unmarshal(data, &bookConfigs)
	if err != nil {
		return nil, err
	}
	return bookConfigs, nil
}
