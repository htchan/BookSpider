package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	RouteAPIKey  = "api"
	RouteLiteKey = "lite"
)

type BackendConfig struct {
	EnabledRoutes    []string `yaml:"enabledRoutes"`
	EnabledSiteNames []string `yaml:"enabledSites"`
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

	return backendConfig, err
}

func (conf BackendConfig) ContainsRoute(routeKey string) bool {
	for _, key := range conf.EnabledRoutes {
		if key == routeKey {
			return true
		}
	}
	return false
}
