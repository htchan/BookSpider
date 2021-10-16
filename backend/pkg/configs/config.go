package configs

import (
	"gopkg.in/yaml.v2"
	"encoding/json"
	"io/ioutil"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/sites"
)

type Config struct {
	Sites map[string]map[string]string `yaml:"sites"`
	MaxThreads int `yaml:"maxThreads"`
	MaxExploreError int `yaml:"maxExploreError"`
	ConstSleep int `yaml: "constSleep"`
	Backend struct {
		Api []string `yaml:"api"`
		StageFile string `yaml:"stageFile"`
		LogFile string `yaml:"logFile"`
	}
}

func LoadSiteYaml(siteName string, config map[string]string) sites.Site {
	data, err := ioutil.ReadFile(config["configFileLocation"])
	utils.CheckError(err);
	var info map[string]string;
	utils.CheckError(yaml.Unmarshal(data, &info))
	site, err := sites.LoadSite(siteName, config, info)
	utils.CheckError(err)
	return *site
}

func LoadSitesYaml(config Config) map[string]sites.Site {
	sites := make(map[string]sites.Site)
	for siteName, siteConfig := range config.Sites {
		sites[siteName] = LoadSiteYaml(siteName, siteConfig)
	}
	return sites
}

func LoadConfigYaml(configFileLocation string) Config {
	var config Config
	data, err := ioutil.ReadFile(configFileLocation)
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &config))
	return config
}

func LoadSiteJson(siteName string, config map[string]string) sites.Site {
	data, err := ioutil.ReadFile(config["configFileLocation"])
	utils.CheckError(err);
	var info map[string]string;
	utils.CheckError(json.Unmarshal(data, &info))
	site, err := sites.LoadSite(siteName, config, info)
	utils.CheckError(err)
	return *site
}

func LoadSitesJson(config Config) map[string]sites.Site {
	sites := make(map[string]sites.Site)
	for siteName, config := range config.Sites {
		sites[siteName] = LoadSiteJson(siteName, config)
	}
	return sites
}

func LoadConfigJson(configFileLocation string) Config {
	var config Config
	data, err := ioutil.ReadFile(configFileLocation)
	utils.CheckError(err)
	utils.CheckError(json.Unmarshal(data, &config))
	return config
}