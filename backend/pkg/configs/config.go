package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"

	"github.com/htchan/BookSpider/internal/utils"
)

type SiteConfig struct {
    BookMeta *BookConfig
	ThreadsCount int `yaml:"threadsCount"`
	MAX_THREAD_COUNT int `yaml:"maxExploreError"`
	ConstSleep int `yaml:"constSleep"`
	Decoder *encoding.Decoder
    DatabaseLocation string `yaml:"databaseLocation"`
	StorageDirectory string `yaml:"downloadDirectory"`
	DecoderString string `yaml:"decode"`
	BookConfigLocation string `yaml:"configFileLocation"`
}

type Config struct {
    SiteConfigs map[string]*SiteConfig `yaml:"sites"`
	Backend struct {
		Api []string `yaml:"api"`
		StageFile string `yaml:"stageFile"`
		LogFile string `yaml:"logFile"`
	}
}

func LoadConfigYaml(configLocation string) (config *Config) {
	defer utils.Recover(func() { config = nil })
	config = new(Config)

	data, err := ioutil.ReadFile(configLocation)
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, config))
	for key, value := range config.SiteConfigs {
		config.SiteConfigs[key].BookMeta = LoadBookConfigYaml(value.BookConfigLocation)
		config.SiteConfigs[key].BookMeta.CONST_SLEEP = config.SiteConfigs[key].ConstSleep
		if (value.DecoderString == "big5") {
			config.SiteConfigs[key].Decoder = traditionalchinese.Big5.NewDecoder()
			config.SiteConfigs[key].BookMeta.Decoder = config.SiteConfigs[key].Decoder
		}
	}
	return
}

func LoadConfigJson(configLocation string) Config {
	return Config{}
}
