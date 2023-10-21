package circuitbreaker

import "time"

type CircuitBreakerClientConfig struct {
	OpenThreshold         uint32        `yaml:"open_threshold" validate:"min=10"`
	AcquireTimeout        time.Duration `yaml:"acquire_timeout" validate:"min=100ms"`
	MaxConcurrencyThreads int64         `yaml:"max_concurrency_threads" validate:"min=2"`
	RecoverThreads        []int64       `yaml:"recover_threads" validate:"min=1,dive,min=1"`
	OpenDuration          time.Duration `yaml:"open_duration" validate:"min=1s"`
	RecoverDuration       time.Duration `yaml:"recover_duration" validate:"min=1s"`
	CheckConfigs          []CheckConfig `yaml:"check_configs" validate:"dive"`
}

type CheckType string

const (
	CheckTypeStatusCodes CheckType = "status-codes"
)

type CheckConfig struct {
	Type  CheckType `yaml:"type" validate:"oneof=status-codes"`
	Value any       `yaml:"value"`
}
