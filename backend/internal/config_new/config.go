package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	APIConfig       APIConfig             `yaml:"api"`
	BatchConfig     BatchConfig           `yaml:"batch"`
	SiteConfigs     map[string]SiteConfig `yaml:"sites" validate:"dive"`
	DatabaseConfig  DatabaseConfig        `yaml:"database"`
	ConfigDirectory string                `env:"CONFIG_DIRECTORY,required" validate:"dir"`
}

type APIConfig struct {
	APIRoutePrefix     string   `env:"NOVEL_SPIDER_API_ROUTE_PREFIX,required" validate:"startswith=/,endsnotwith=/"`
	LiteRoutePrefix    string   `env:"NOVEL_SPIDER_LITE_ROUTE_PREFIX,required" validate:"startswith=/,endsnotwith=/"`
	AvailableSiteNames []string `env:"API_AVAILABLE_SITES,required" validate:"min=1,dive,min=1"`
}

type BatchConfig struct {
	MaxWorkingThreads  int      `env:"MAX_WORKING_THREADS" validate:"min=1"`
	AvailableSiteNames []string `env:"BATCH_AVAILABLE_SITES,required" validate:"min=1,dive,min=1"`
}

type DatabaseConfig struct {
	Host            string        `env:"PSQL_HOST,required" validate:"min=1"`
	Port            string        `env:"PSQL_PORT,required" validate:"min=1"`
	User            string        `env:"PSQL_USER,required" validate:"min=1"`
	Password        string        `env:"PSQL_PASSWORD,required" validate:"min=1"`
	Name            string        `env:"PSQL_NAME,required" validate:"min=1"`
	MaxOpenConns    int           `env:"PSQL_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `env:"PSQL_MAX_IDLE_CONNS"`
	ConnMaxIdleTime time.Duration `env:"PSQL_CONN_MAX_IDLE_TIME"`
}

func LoadConfig() (*Config, error) {
	var conf Config

	loadConfigFuncs := []func() error{
		func() error { return env.Parse(&conf) },
		func() error { return env.Parse(&conf.APIConfig) },
		func() error { return env.Parse(&conf.BatchConfig) },
		func() error { return env.Parse(&conf.DatabaseConfig) },
		func() error {
			var referenceData []byte

			filepath.Walk(conf.ConfigDirectory, func(path string, file fs.FileInfo, _ error) error {
				if filepath.Ext(path) == ".yaml" && file.Name() != "main.yaml" {
					data, err := os.ReadFile(path)
					if err == nil {
						referenceData = append(referenceData, data...)
					}
				}
				return nil
			})

			configData, err := os.ReadFile(conf.ConfigDirectory + "/main.yaml")
			if err != nil {
				return err
			}

			fullConfig := struct {
				Sites map[string]SiteConfig `yaml:"sites"`
			}{}

			err = yaml.Unmarshal(append(referenceData, configData...), &fullConfig)
			if err != nil {
				return err
			}

			conf.SiteConfigs = fullConfig.Sites

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
