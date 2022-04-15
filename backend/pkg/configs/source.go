package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
	"path/filepath"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"fmt"

	"github.com/htchan/BookSpider/internal/utils"
)

type SourceConfig struct {
	BaseUrl string `yaml:"baseUrl"`
	DownloadUrl string `yaml:"downloadUrl"`
	ChapterUrl string `yaml:"chapterUrl"`
	ChapterUrlPattern string `yaml:"chapterUrlPattern"`
	DecoderString string `yaml:"decode"`
	Threads int `yaml:"threads"`
	DownloadBookThreads int `yaml:"downloadBookThreads"`
	ConstSleep int `yaml:"constSleep"`
	MaxExploreError int `yaml:"maxExploreError"`
	UseRequestInterval bool `yaml:"useRequestInterval"`
	SourceKey string
	Decoder *encoding.Decoder
}

func LoadSourceConfigs(configDirectory string) (config map[string]*SourceConfig) {
	defer utils.Recover(func() { config = nil })

	data, err := ioutil.ReadFile(filepath.Join(configDirectory, "source_config.yaml"))
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &config))
	for key, value := range config {
		if (value.DecoderString == "big5") {
			value.Decoder = traditionalchinese.Big5.NewDecoder()
			value.SourceKey = key
			config[key] = value
		}
	}
	return
}

func (config SourceConfig)Populate(id int) SourceConfig {
	config.BaseUrl = fmt.Sprintf(config.BaseUrl, id)
	config.DownloadUrl = fmt.Sprintf(config.DownloadUrl, id)
	return config
}