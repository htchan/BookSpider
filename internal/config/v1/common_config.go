package config

import "time"

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

type TraceConfig struct {
	OtelURL         string `env:"OTEL_URL,required" validate:"url"`
	OtelServiceName string `env:"OTEL_SERVICE_NAME,required" validate:"min=1"`
}

// TODO: we will use NATS to receive the task request
type NatsConfig struct{}

type CommonConfig struct {
	StoragePath    string `env:"STORAGE_PATH,required" validate:"dir"`
	ConfigLocation string `env:"CONFIG_LOCATION,required" validate:"file"`
}
