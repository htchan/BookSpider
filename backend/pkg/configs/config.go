package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"os"

	"github.com/htchan/BookSpider/internal/utils"
)

type SiteConfig struct {
    BookMeta *BookConfig
	DownloadBookThreads int `yaml:"downloadBookThreads"`
	Threads int `yaml:"threads"`
	MaxExploreError int `yaml:"maxExploreError"`
	ConstSleep int `yaml:"constSleep"`
	Decoder *encoding.Decoder
	DatabaseEngine string `yaml:"databaseEngine"`
    DatabaseLocation string `yaml:"databaseLocation"`
	StorageDirectory string `yaml:"downloadDirectory"`
	DecoderString string `yaml:"decode"`
	BookConfigLocation string `yaml:"configFileLocation"`
	BackupDirectory string `yaml:"backupDirectory"`
	UseRequestInterval bool `yaml:"useRequestInterval"`
}

type Config struct {
    SiteConfigs map[string]*SiteConfig `yaml:"sites"`
	MaxThreads int `yaml:"maxThreads"`
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
		config.SiteConfigs[key].BookMeta = LoadBookConfigYaml(os.Getenv("ASSETS_LOCATION") + value.BookConfigLocation)
		config.SiteConfigs[key].BookMeta.UseRequestInterval = value.UseRequestInterval
		config.SiteConfigs[key].BookMeta.CONST_SLEEP = value.ConstSleep
		config.SiteConfigs[key].BookMeta.StorageDirectory = value.StorageDirectory
		if (value.DecoderString == "big5") {
			config.SiteConfigs[key].Decoder = traditionalchinese.Big5.NewDecoder()
			config.SiteConfigs[key].BookMeta.Decoder = value.Decoder
		}
	}
	return
}

func LoadConfigJson(configLocation string) Config {
	return Config{}
}
