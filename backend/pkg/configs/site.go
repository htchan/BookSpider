package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/htchan/BookSpider/internal/utils"
	// "fmt"
)

type SiteConfig struct {
	SourceKey string `yaml:"source_key"`
	DatabaseEngine string `yaml:"database_engine"`
	DatabaseLocation string `yaml:"database_location"`
	StorageDirectory string `yaml:"download_directory"`
	BackupDirectory string `yaml:"backup_directory"`

	SourceConfig SourceConfig
}

func LoadSiteConfigs(configDirectory string) (config map[string]SiteConfig) {
	defer utils.Recover(func() { config = nil })

	data, err := ioutil.ReadFile(filepath.Join(configDirectory, "site_config.yaml"))
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &config))
	
	sourceConfigs := LoadSourceConfigs(configDirectory)
	for key, value := range config {
		value.SourceConfig = sourceConfigs[value.SourceKey]
		config[key] = value
	}
	return
}
