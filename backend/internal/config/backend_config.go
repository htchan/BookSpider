package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

const (
	RouteAPIKey  = "api"
	RouteLiteKey = "lite"
)

type BackendConfig struct {
	EnabledRoutes    map[string]bool `yaml:"enabledRoutes"`
	EnabledSiteNames []string        `yaml:"enabledSites"`

	SiteConfigs map[string]*SiteConfig
}

func LoadBackendConfig(configDirectory string) (BackendConfig, error) {
	data, err := os.ReadFile(filepath.Join(configDirectory, "backend_config.yaml"))
	if err != nil {
		return BackendConfig{}, err
	}
	var backendConfig BackendConfig
	err = yaml.Unmarshal(data, &backendConfig)
	if err != nil {
		return backendConfig, err
	}

	backendConfig.SiteConfigs, err = LoadSiteConfigs(configDirectory)
	return backendConfig, err
}
