package config

import (
	"errors"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type WorkerConfig struct {
	WorkerCount    int                     `yaml:"worker_count" json:"worker_count"`
	AvailableSites []string                `yaml:"available_sites" json:"available_sites"`
	Clients        map[string]ClientConfig `yaml:"clients" json:"client"`

	Database DatabaseConfig
	Trace    TraceConfig
	Common   CommonConfig
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	conf := &WorkerConfig{}
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

func (conf *WorkerConfig) Validate() error {
	return validator.New().Struct(conf)
}
