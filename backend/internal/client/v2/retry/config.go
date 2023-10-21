package retry

import (
	"time"
)

type RetryConditionType string

const (
	RetryConditionTypeStatusCode   RetryConditionType = "status-codes"
	RetryConditionTypeTimeout      RetryConditionType = "error"
	RetryConditionTypeBodyContains RetryConditionType = "body-contains"
)

type RetryCondition struct {
	Type              RetryConditionType `yaml:"type" validate:"oneof=status-codes error body-contains"`
	Value             any                `yaml:"value"`
	Weight            int                `yaml:"weight" validate:"min=1"`
	PauseInterval     time.Duration      `yaml:"pause_interval" validate:"min=1s"`
	PauseIntervalType PauseIntervalType  `yaml:"pause_interval_type" validate:"oneof=const linear exponential"`
}

type PauseIntervalType string

const (
	PauseIntervalTypeConst       PauseIntervalType = "const"
	PauseIntervalTypeLinear      PauseIntervalType = "linear"
	PauseIntervalTypeExponential PauseIntervalType = "exponential"
)

type RetryClientConfig struct {
	RetryConditions []RetryCondition `yaml:"retry_conditions" validate:"dive"`
	MaxRetryWeight  int              `yaml:"max_retry_weight" validate:"min=1"`
}
