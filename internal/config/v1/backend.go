package config

import (
	"os"
	"path/filepath"
	"slices"

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
	return slices.Contains(conf.EnabledRoutes, routeKey)
}
