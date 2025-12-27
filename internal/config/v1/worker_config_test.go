package config

import (
	"strings"
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/stretchr/testify/assert"
)

func TestLoadWorkerConfig(t *testing.T) {

	envMap := map[string]string{
		"PSQL_HOST":               "psql.host",
		"PSQL_PORT":               "5432",
		"PSQL_USER":               "psql.user",
		"PSQL_PASSWORD":           "psql.password",
		"PSQL_NAME":               "psql.db",
		"PSQL_MAX_OPEN_CONNS":     "10",
		"PSQL_MAX_IDLE_CONNS":     "5",
		"PSQL_CONN_MAX_IDLE_TIME": "5s",

		"OTEL_URL":          "http://otel.host",
		"OTEL_SERVICE_NAME": "otel test name",

		"NATS_URL": "nats://localhost:4222",

		"STORAGE_PATH":    ".",
		"CONFIG_LOCATION": "../../../config/v1/worker.yml",
	}

	tests := []struct {
		name        string
		setupEnv    func(*testing.T)
		wantConfig  *WorkerConfig
		assertError func(*testing.T, error)
	}{
		{
			name: "happy flow",
			setupEnv: func(t *testing.T) {
				for key, val := range envMap {
					t.Setenv(key, val)
				}
			},
			wantConfig: &WorkerConfig{
				WorkerCount:    10,
				AvailableSites: []string{"site1", "site2"},
				Clients: map[string]ClientConfig{
					"test": {
						DecodeMethod: "big5",
						Retry: RetryConfig{
							MaxRetryCount:       20,
							LinearRetryInterval: 3 * time.Second,
						},
						Pool: ClientPoolConfig{
							RefreshInterval:            168 * time.Hour,
							DropClientFailureThreshold: 20,
							FailureCooldownInterval:    60 * time.Second,
							SuccessCooldownInterval:    10 * time.Second,
							ClientTimeout:              60 * time.Second,
							Socks5ProxySourceURL:       "https://test.com/socks5",
							Socks4ProxySourceURL:       "https://test.com/socks4",
							HTTPProxySourceURL:         "https://test.com/http",
						},
					},
				},
				Database: DatabaseConfig{
					Host:            "psql.host",
					Port:            "5432",
					User:            "psql.user",
					Password:        "psql.password",
					Name:            "psql.db",
					MaxOpenConns:    10,
					MaxIdleConns:    5,
					ConnMaxIdleTime: 5 * time.Second,
				},
				Trace: TraceConfig{
					OtelURL:         "http://otel.host",
					OtelServiceName: "otel test name",
				},
				Nats: NatsConfig{
					URL: "nats://localhost:4222",
				},
				Common: CommonConfig{
					StoragePath:    ".",
					ConfigLocation: "../../../config/v1/worker.yml",
				},
			},
			assertError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error flow/common config not exist",
			setupEnv: func(t *testing.T) {
				for key, val := range envMap {
					if key == "STORAGE_PATH" {
						continue
					}

					t.Setenv(key, val)
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "STORAGE_PATH"})
			},
		},
		{
			name: "error flow/db config not exist",
			setupEnv: func(t *testing.T) {
				for key, val := range envMap {
					if strings.HasPrefix(key, "PSQL") {
						continue
					}

					t.Setenv(key, val)
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_HOST"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_PORT"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_USER"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_PASSWORD"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_NAME"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_MAX_OPEN_CONNS"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_MAX_IDLE_CONNS"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "PSQL_CONN_MAX_IDLE_TIME"})
			},
		},
		{
			name: "error flow/trace config not exist",
			setupEnv: func(t *testing.T) {
				for key, val := range envMap {
					if strings.HasPrefix(key, "OTEL") {
						continue
					}

					t.Setenv(key, val)
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "OTEL_URL"})
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "OTEL_SERVICE_NAME"})
			},
		},
		{
			name: "error fllow/nats config not exist",
			setupEnv: func(t *testing.T) {
				for key, val := range envMap {
					if strings.HasPrefix(key, "NATS") {
						continue
					}

					t.Setenv(key, val)
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, env.EnvVarIsNotSetError{Key: "NATS_URL"})
			},
		},
		{
			name: "error flow/config file not exist",
			setupEnv: func(t *testing.T) {
				for key, val := range envMap {
					if key == "CONFIG_LOCATION" {
						t.Setenv(key, "not_exist.yml")
						continue
					}

					t.Setenv(key, val)
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "Key: 'WorkerConfig.Common.ConfigLocation' Error:Field validation for 'ConfigLocation' failed on the 'file' tag")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// this test need to change env which cannot be run in parallel
			test.setupEnv(t)

			config, err := LoadWorkerConfig()
			assert.Equal(t, test.wantConfig, config)
			test.assertError(t, err)
		})
	}
}
