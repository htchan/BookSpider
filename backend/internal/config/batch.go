package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type BatchConfig struct {
	MaxCommonThreads int      `yaml:"maxCommonThreads"`
	EnabledSites     []string `yaml:"enabledSites"`
}

func LoadBatchConfig(configDirectory string) (BatchConfig, error) {
	data, err := os.ReadFile(filepath.Join(configDirectory, "batch_config.yaml"))
	if err != nil {
		return BatchConfig{}, err
	}
	var batchConfig BatchConfig
	err = yaml.Unmarshal(data, &batchConfig)
	if err != nil {
		return batchConfig, err
	}

	return batchConfig, err
}
