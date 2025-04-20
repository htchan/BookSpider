package circuitbreaker

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_validate_CheckConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  CheckConfig
		valid bool
	}{
		{
			name: "valid",
			conf: CheckConfig{
				Type:  CheckTypeStatusCodes,
				Value: nil,
			},
			valid: true,
		},
		{
			name: "invalid type",
			conf: CheckConfig{
				Type:  "unknown",
				Value: nil,
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := validator.New().Struct(test.conf)
			assert.Equal(t, test.valid, err == nil)
		})
	}
}

func Test_validate_CircuitBreakerConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  CircuitBreakerClientConfig
		valid bool
	}{
		{
			name: "valid",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1},
				OpenDuration:          time.Second,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: true,
		},
		{
			name: "invalid open threshold",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         9,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1},
				OpenDuration:          time.Second,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: false,
		},
		{
			name: "invalid acquire timeout",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        99 * time.Millisecond,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1},
				OpenDuration:          time.Second,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: false,
		},
		{
			name: "invalid max concurrency threads",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 1,
				RecoverThreads:        []int64{1},
				OpenDuration:          time.Second,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: false,
		},
		{
			name: "invalid recover threads/empty slice",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{},
				OpenDuration:          time.Second,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: false,
		},
		{
			name: "invalid recover threads/item minimum violation",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1, 0, 2},
				OpenDuration:          time.Second,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: false,
		},
		{
			name: "invalid open duration",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1},
				OpenDuration:          999 * time.Millisecond,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: false,
		},
		{
			name: "invalid recover duration",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1},
				OpenDuration:          time.Second,
				RecoverDuration:       999 * time.Millisecond,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: nil},
				},
			},
			valid: false,
		},
		{
			name: "invalid check configs",
			conf: CircuitBreakerClientConfig{
				OpenThreshold:         10,
				AcquireTimeout:        time.Second,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1},
				OpenDuration:          time.Second,
				RecoverDuration:       time.Second,
				CheckConfigs: []CheckConfig{
					{Type: "unknown", Value: nil},
				},
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := validator.New().Struct(test.conf)
			assert.Equal(t, test.valid, err == nil)
		})
	}
}
