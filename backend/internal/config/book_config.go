package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
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
	Storage          string
	SourceKey        string
}

type RequestFunc func(string) (string, error)

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
	for key, value := range bookConfigs {
		value.SourceKey = key
		bookConfigs[key] = value
	}
	return bookConfigs, nil
}
