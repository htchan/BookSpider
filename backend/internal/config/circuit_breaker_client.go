package config

import (
	"gopkg.in/yaml.v2"
	"os"
	
	"path/filepath"

	"github.com/htchan/BookSpider/internal/utils"
)

type CircuitBreakerConfig struct {
	MaxFailCount int `yaml:"maxFailCount"`
	MaxFailMultiplier float64 `yaml:"MaxFailMultiplier"`
	CircuitBreakingSleep int `yaml:"circuitBreakingSleep"`
	IntervalSleep int `yaml:"internalSleep"`
	Timeout int `yaml:"timeout"`
	Retry503 int `yaml:"retry503"`
	RetryErr int `yaml:"retryErr"`
}

func LoadClientConfigs(configDirectory string) (config map[string]CircuitBreakerConfig) {
	defer utils.Recover(func() { config = nil })

	data, err := os.ReadFile(filepath.Join(configDirectory, "client_config.yaml"))
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, &config))

	return
}
