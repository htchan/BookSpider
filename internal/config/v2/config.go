package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"

	"github.com/go-playground/validator/v10"
)

type APIConfig struct {
	APIRoutePrefix     string                `env:"NOVEL_SPIDER_API_ROUTE_PREFIX,required" validate:"startswith=/,endsnotwith=/"`
	LiteRoutePrefix    string                `env:"NOVEL_SPIDER_LITE_ROUTE_PREFIX,required" validate:"startswith=/,endsnotwith=/"`
	TraceConfig        TraceConfig           `validate:"dive"`
	AvailableSiteNames []string              `env:"API_AVAILABLE_SITES,required" validate:"min=1,dive,min=1"`
	SiteConfigs        map[string]SiteConfig `yaml:"sites" validate:"dive"`
	DatabaseConfig     DatabaseConfig        `yaml:"database"`
	ConfigDirectory    string                `env:"CONFIG_DIRECTORY,required" validate:"dir"`
}

type WorkerConfig struct {
	MaxWorkingThreads  int                   `env:"MAX_WORKING_THREADS" validate:"min=1"`
	AvailableSiteNames []string              `env:"BATCH_AVAILABLE_SITES,required" validate:"min=1,dive,min=1"`
	TraceConfig        TraceConfig           `validate:"dive"`
	SiteConfigs        map[string]SiteConfig `yaml:"sites" validate:"dive"`
	DatabaseConfig     DatabaseConfig        `yaml:"database"`
	ScheduleConfig     ScheduleConfig        `yaml:"schedule"`
	ConfigDirectory    string                `env:"CONFIG_DIRECTORY,required" validate:"dir"`
}
type TraceConfig struct {
	OtelURL         string `env:"OTEL_URL,required" validate:"url"`
	OtelServiceName string `env:"OTEL_SERVICE_NAME,required" validate:"min=1"`
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

func LoadAPIConfig() (*APIConfig, error) {
	var conf APIConfig

	loadConfigFuncs := []func() error{
		func() error { return env.Parse(&conf) },
		func() error { return env.Parse(&conf.DatabaseConfig) },
		func() error { return env.Parse(&conf.TraceConfig) },
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

func (conf *APIConfig) Validate() error {
	return validator.New().Struct(conf)
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	var conf WorkerConfig

	loadConfigFuncs := []func() error{
		func() error { return env.Parse(&conf) },
		func() error { return env.Parse(&conf.DatabaseConfig) },
		func() error { return env.Parse(&conf.TraceConfig) },
		func() error { return env.Parse(&conf.ScheduleConfig) },
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

func (conf *WorkerConfig) Validate() error {
	validStruct := validator.New().Struct(conf)
	if validStruct != nil {
		return validStruct
	}

	if conf.ScheduleConfig.IntervalMonth+conf.ScheduleConfig.IntervalMonth <= 0 {
		return errors.New("interval is zero")
	}

	return nil
}
