package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	APIConfig      APIConfig             `yaml:"api"`
	BatchConfig    BatchConfig           `yaml:"batch"`
	SiteConfigs    map[string]SiteConfig `yaml:"sites" validate:"dive"`
	DatabaseConfig DatabaseConfig        `yaml:"database"`
	ConfigLocation string                `env:"CONFIG_LOCATION,required" validate:"file"`
}

type APIConfig struct {
	APIRoutePrefix     string   `env:"API_ROUTE_PREFIX" envDefault:"/api/novel" validate:"startswith=/,endsnotwith=/"`
	LiteRoutePrefix    string   `env:"LITE_ROUTE_PREFIX" envDefault:"/lite/novel" validate:"startswith=/,endsnotwith=/"`
	AvailableSiteNames []string `env:"API_AVAILABLE_SITES,required" validate:"min=1,dive,min=1"`
}

type BatchConfig struct {
	MaxWorkingThreads  int      `env:"MAX_WORKING_THREADS" validate:"min=1"`
	AvailableSiteNames []string `env:"BATCH_AVAILABLE_SITES,required" validate:"min=1,dive,min=1"`
}

type DatabaseConfig struct {
	Host     string `env:"PSQL_HOST,required" validate:"min=1"`
	Port     string `env:"PSQL_PORT,required" validate:"min=1"`
	User     string `env:"PSQL_USER,required" validate:"min=1"`
	Password string `env:"PSQL_PASSWORD,required" validate:"min=1"`
	Name     string `env:"PSQL_NAME,required" validate:"min=1"`
}

func LoadConfig() (*Config, error) {
	var conf Config

	loadConfigFuncs := []func() error{
		func() error { return env.Parse(&conf) },
		func() error { return env.Parse(&conf.APIConfig) },
		func() error { return env.Parse(&conf.BatchConfig) },
		func() error { return env.Parse(&conf.DatabaseConfig) },
		func() error {
			data, err := os.ReadFile(conf.ConfigLocation)
			if err != nil {
				return err
			}

			err = yaml.Unmarshal(data, &conf.SiteConfigs)
			if err != nil {
				return err
			}

			return nil
		},
	}

	for _, f := range loadConfigFuncs {
		if err := f(); err != nil {
			return nil, fmt.Errorf("parse config: %w", err)
		}
	}

	return &conf, nil
}

func (conf *Config) Validate() error {
	return validator.New().Struct(conf)
}
