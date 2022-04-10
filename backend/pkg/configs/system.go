package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"

	"github.com/htchan/BookSpider/internal/utils"
)

type SystemConfig struct {
	AvailableSites []string `yaml:"enabled_sites"`
}

func LoadSystemConfigs(systemConfigLocation string) (configs *SystemConfig) {
	defer utils.Recover(func() { configs = nil })

	data, err := ioutil.ReadFile(systemConfigLocation)
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &configs))
	return
}
