package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIConfig      APIConfig             `yaml:"api"`
	BatchConfig    BatchConfig           `yaml:"batch"`
	SiteConfigs    map[string]SiteConfig `yaml:"sites"`
	DatabaseConfig DatabaseConfig        `yaml:"database"`
	ConfigLocation string                `env:"CONFIG_LOCATION,required"`
}

type APIConfig struct {
	APIRoutePrefix     string   `env:"API_ROUTE_PREFIX" envDefault:"/api/novel"`
	LiteRoutePrefix    string   `env:"LITE_ROUTE_PREFIX" envDefault:"/lite/novel"`
	AvailableSiteNames []string `env:"API_AVAILABLE_SITES,required"`
}

type BatchConfig struct {
	MaxWorkingThreads  int      `env:"MAX_WORKING_THREADS"`
	AvailableSiteNames []string `env:"BATCH_AVAILABLE_SITES,required"`
}

type DatabaseConfig struct {
	Host     string `env:"PSQL_HOST,required"`
	Port     string `env:"PSQL_PORT,required"`
	User     string `env:"PSQL_USER,required"`
	Password string `env:"PSQL_PASSWORD,required"`
	Name     string `env:"PSQL_NAME,required"`
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
