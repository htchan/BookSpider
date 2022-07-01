package config

import (
	"gopkg.in/yaml.v2"
	"os"

	"path/filepath"
)

type SiteConfig struct {
	SourceKey       string `yaml:"sourceKey"`
	ClientKey       string `yaml:"clientKey"`
	Storage         string `yaml:"storage"`
	BackupDirectory string `yaml:"backupDirectory"`

	ConcurrencyConfig ConcurrencyConfig `yaml:"concurrency"`
	DatabaseConfig    DatabaseConfig    `yaml:"database"`

	MaxExploreError int `yaml:"maxExploreError"`

	AvailabilityConfig AvailabilityConfig `yaml:"availability"`

	*BookConfig
	*CircuitBreakerClientConfig
}

type ConcurrencyConfig struct {
	FetchThreads    int `yaml:"fetchThreads"`
	DownloadThreads int `yaml:"downloadThreads"	`
}

type DatabaseConfig struct {
	Engine    string `yaml:"engine"`
	Conn      string `yaml:"conn"`
	Timeout   int    `yaml:"timeout"`
	IdleConns int    `yaml:"idleConns"`
	OpenConns int    `yaml:"openConns"`
}

type AvailabilityConfig struct {
	URL   string `yaml:"url"`
	Check string `yaml:"check"`
}

func LoadSiteConfigs(configDirectory string) (map[string]*SiteConfig, error) {
	data, err := os.ReadFile(filepath.Join(configDirectory, "site_configs.yaml"))
	if err != nil {
		return nil, err
	}

	var siteConfigs map[string]*SiteConfig
	err = yaml.Unmarshal(data, &siteConfigs)
	if err != nil {
		return siteConfigs, err
	}

	bookConfigs, err := LoadBookConfigs(configDirectory)
	if err != nil {
		return siteConfigs, err
	}
	clientConfigs, err := LoadClientConfigs(configDirectory)
	if err != nil {
		return siteConfigs, err
	}

	for key, value := range siteConfigs {
		// value.BookConfig = new(BookConfig)
		value.BookConfig = bookConfigs[value.SourceKey]
		value.CircuitBreakerClientConfig = clientConfigs[value.ClientKey]
		value.BookConfig.Storage = value.Storage

		siteConfigs[key] = value
	}
	return siteConfigs, nil
}
