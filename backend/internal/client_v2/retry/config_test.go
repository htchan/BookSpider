package retry

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_validate_RetryCondition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		cond  RetryCondition
		valid bool
	}{
		{
			name: "pass validation",
			cond: RetryCondition{
				Type:              RetryConditionTypeBodyContains,
				Value:             "some body",
				Weight:            10,
				PauseInterval:     1 * time.Second,
				PauseIntervalType: PauseIntervalTypeConst,
			},
			valid: true,
		},
		{
			name: "invalid type",
			cond: RetryCondition{
				Type:              "unknown",
				Value:             "some body",
				Weight:            10,
				PauseInterval:     1 * time.Second,
				PauseIntervalType: PauseIntervalTypeLinear,
			},
			valid: false,
		},
		{
			name: "invalid weight",
			cond: RetryCondition{
				Type:              RetryConditionTypeStatusCode,
				Weight:            0,
				PauseInterval:     1 * time.Second,
				PauseIntervalType: PauseIntervalTypeExponential,
			},
			valid: false,
		},
		{
			name: "invalid pause interval",
			cond: RetryCondition{
				Type:              RetryConditionTypeTimeout,
				Value:             "some body",
				Weight:            10,
				PauseInterval:     999 * time.Millisecond,
				PauseIntervalType: PauseIntervalTypeConst,
			},
			valid: false,
		},
		{
			name: "invalid pause interval type",
			cond: RetryCondition{
				Type:              RetryConditionTypeBodyContains,
				Value:             "some body",
				Weight:            10,
				PauseInterval:     1 * time.Second,
				PauseIntervalType: "invalid",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.cond)
			assert.Equal(t, test.valid, err == nil)
		})
	}
}

func Test_validate_RetryClientConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  RetryClientConfig
		valid bool
	}{
		{
			name: "valid",
			conf: RetryClientConfig{
				MaxRetryWeight: 10,
				RetryConditions: []RetryCondition{
					{
						Type:              RetryConditionTypeBodyContains,
						Value:             "some body",
						Weight:            10,
						PauseInterval:     1 * time.Second,
						PauseIntervalType: PauseIntervalTypeConst,
					},
				},
			},
			valid: true,
		},
		{
			name: "invalid retry condition",
			conf: RetryClientConfig{
				MaxRetryWeight: 0,
				RetryConditions: []RetryCondition{
					{
						Type:              RetryConditionTypeBodyContains,
						Value:             "some body",
						Weight:            10,
						PauseInterval:     1 * time.Second,
						PauseIntervalType: PauseIntervalTypeConst,
					},
				},
			},
			valid: false,
		},
		{
			name: "invalid max retry weight",
			conf: RetryClientConfig{
				MaxRetryWeight: 10,
				RetryConditions: []RetryCondition{
					{
						Type:              RetryConditionTypeBodyContains,
						Value:             "some body",
						Weight:            0,
						PauseInterval:     1 * time.Second,
						PauseIntervalType: PauseIntervalTypeConst,
					},
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
