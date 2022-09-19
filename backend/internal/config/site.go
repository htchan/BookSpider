package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ConcurrencyConfig struct {
	DownloadThreads int `yaml:"DownloadThreads"`
}

type DatabaseConfig struct {
	Engine    string `yaml:"engine"`
	Timeout   int    `yaml:"timeout"`
	IdleConns int    `yaml:"idleConns"`
	OpenConns int    `yaml:"openConns"`
}

type AvailabilityConfig struct {
	URL   string `yaml:"url"`
	Check string `yaml:"check"`
}

type SiteConfig struct {
	BookKey            string             `yaml:"bookKey"`
	ClientKey          string             `yaml:"clientKey"`
	Storage            string             `yaml:"storage"`
	BackupDirectory    string             `yaml:"backupDirectory"`
	MaxExploreError    int                `yaml:"maxExploreError"`
	ConcurrencyConfig  ConcurrencyConfig  `yaml:"concurrency"`
	DatabaseConfig     DatabaseConfig     `yaml:"database"`
	AvailabilityConfig AvailabilityConfig `yaml:"availability"`
}

func LoadSiteConfigs(configDirectory string) (map[string]*SiteConfig, error) {
	data, err := os.ReadFile(filepath.Join(configDirectory, "site_configs.yaml"))
	if err != nil {
		return nil, err
	}

	var siteConfigs map[string]*SiteConfig
	err = yaml.Unmarshal(data, &siteConfigs)
	if err != nil {
		return nil, err
	}
	return siteConfigs, nil
}
