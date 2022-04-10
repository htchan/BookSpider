package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"

	"github.com/htchan/BookSpider/internal/utils"
)

type ServerConfig struct {
	AvailableApi []string `yaml:"api"`
}

func LoadServerConfigs(serverConfigLocation string) (config *ServerConfig) {
	defer utils.Recover(func() { config = nil })
	config = new(ServerConfig)

	data, err := ioutil.ReadFile(serverConfigLocation)
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, config))
	return
}