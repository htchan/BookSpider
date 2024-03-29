package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"path/filepath"
)

type DecoderConfig struct {
	Method string `yaml:"method"`
}

type CircuitBreakerClientConfig struct {
	MaxFailCount         int           `yaml:"maxFailCount"`
	MaxFailMultiplier    float64       `yaml:"maxFailMultiplier"`
	CircuitBreakingSleep int           `yaml:"circuitBreakingSleep"`
	IntervalSleep        int           `yaml:"intervalSleep"`
	Timeout              int           `yaml:"timeout"`
	RetryUnavailable     int           `yaml:"retryUnavailable"`
	RetryErr             int           `yaml:"retryErr"`
	DecoderConfig        DecoderConfig `yaml:"decoder"`
	MaxThreads           int           `yaml:"maxThreads"`
}

func LoadClientConfigs(configDirectory string) (map[string]*CircuitBreakerClientConfig, error) {
	data, err := os.ReadFile(filepath.Join(configDirectory, "client_configs.yaml"))
	if err != nil {
		return nil, err
	}

	var clientConfigs map[string]*CircuitBreakerClientConfig
	err = yaml.Unmarshal(data, &clientConfigs)

	return clientConfigs, err
}
