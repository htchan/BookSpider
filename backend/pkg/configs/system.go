package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/htchan/BookSpider/internal/utils"
)

type SystemConfig struct {
	AvailableSiteNames []string `yaml:"enabled_sites"`
	AvailableSiteConfigs map[string]*SiteConfig
}

func LoadSystemConfigs(configDirectory string) (config *SystemConfig) {
	defer utils.Recover(func() { config = nil })

	data, err := ioutil.ReadFile(filepath.Join(configDirectory, "system_config.yaml"))
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &config))
	config.AvailableSiteConfigs = make(map[string]*SiteConfig)
	siteConfig := LoadSiteConfigs(configDirectory)
	for _, siteName := range config.AvailableSiteNames {
		config.AvailableSiteConfigs[siteName] = siteConfig[siteName]
	}
	return
}
