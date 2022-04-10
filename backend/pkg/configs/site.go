package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"

	"github.com/htchan/BookSpider/internal/utils"
	// "fmt"
)

type SiteConfig struct {
	SourceKey string `yaml:"source_key"`
	DatabaseEngine string `yaml:"database_engine"`
	DatabaseLocation string `yaml:"database_location"`
	StorageDirectory string `yaml:"download_directory"`
	BackupDirectory string `yaml:"backup_directory"`
}

func LoadSiteConfigs(siteConfigLocation string) (config map[string]SiteConfig) {
	defer utils.Recover(func() { config = nil })

	data, err := ioutil.ReadFile(siteConfigLocation)
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &config))
	return
}
