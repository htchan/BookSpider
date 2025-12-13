package config

import (
	"errors"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type BackendConfig struct {
	APIRoutePrefix  string   `yaml:"api_route_prefix" json:"api_route_prefix"`
	LiteRoutePrefix string   `yaml:"lite_route_prefix" json:"lite_route_prefix"`
	AvailableSites  []string `yaml:"available_sites" json:"available_sites"`

	Database DatabaseConfig
	Trace    TraceConfig
	Common   CommonConfig
}

func LoadBackendConfig() (*BackendConfig, error) {
	conf := &BackendConfig{}
	commonConfigErr := env.Parse(&conf.Common)
	dbConfigErr := env.Parse(&conf.Database)
	traceConfigErr := env.Parse(&conf.Trace)

	validataConfigErr := validator.New().Struct(conf)
	if commonConfigErr != nil || dbConfigErr != nil || traceConfigErr != nil || validataConfigErr != nil {
		return nil, errors.Join(commonConfigErr, dbConfigErr, traceConfigErr, validataConfigErr)
	}

	data, readFileErr := os.ReadFile(conf.Common.ConfigLocation)
	if readFileErr != nil {
		return nil, readFileErr
	}

	yamlErr := yaml.Unmarshal(data, &conf)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return conf, nil
}
