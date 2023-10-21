package retry

import (
	"fmt"
	"testing"
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePauseDuration(t *testing.T) {
	t.Parallel()

	type args struct {
		attempt      int
		interval     time.Duration
		intervalType PauseIntervalType
	}

	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "calculate duration for const retry interval",
			args: args{
				attempt:      5,
				interval:     time.Second,
				intervalType: PauseIntervalTypeConst,
			},
			want: time.Second,
		},
		{
			name: "calculate duration for linear retry interval",
			args: args{
				attempt:      5,
				interval:     time.Second,
				intervalType: PauseIntervalTypeLinear,
			},
			want: time.Second * 6,
		},
		{
			name: "calculate duration for exponential retry interval",
			args: args{
				attempt:      5,
				interval:     time.Second,
				intervalType: PauseIntervalTypeExponential,
			},
			want: time.Second * 36,
		},
		{
			name: "calculate duration for invalid retry interval",
			args: args{
				attempt:      5,
				interval:     time.Second,
				intervalType: "invalid",
			},
			want: 0,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := CalculatePauseDuration(test.args.attempt, test.args.interval, test.args.intervalType)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNewRetryCheck(t *testing.T) {
	t.Parallel()

	retryCheckForTimeout := NewRetryCheck(RetryCondition{
		Type:              RetryConditionTypeTimeout,
		Weight:            10,
		PauseInterval:     1 * time.Second,
		PauseIntervalType: PauseIntervalTypeLinear,
	})

	retryCheckForStatusCode := NewRetryCheck(RetryCondition{
		Type:              RetryConditionTypeStatusCode,
		Value:             []interface{}{100, 300, 400},
		Weight:            20,
		PauseInterval:     2 * time.Second,
		PauseIntervalType: PauseIntervalTypeConst,
	})

	retryCheckForBodyContains := NewRetryCheck(RetryCondition{
		Type:              RetryConditionTypeBodyContains,
		Value:             "d.*a",
		Weight:            30,
		PauseInterval:     3 * time.Second,
		PauseIntervalType: PauseIntervalTypeExponential,
	})

	type args struct {
		attempt int
		body    string
		err     error
	}

	type result struct {
		shouldRetry   bool
		weight        int
		pauseDuration time.Duration
	}

	tests := []struct {
		name       string
		args       args
		retryCheck RetryCheck
		want       result
	}{
		{
			name: "retry for error/timeout",
			args: args{
				attempt: 3,
				body:    "body",
				err:     fmt.Errorf("stub error: %w", client.ErrTimeout),
			},
			retryCheck: retryCheckForTimeout,
			want: result{
				shouldRetry:   true,
				weight:        10,
				pauseDuration: 4 * time.Second,
			},
		},
		{
			name: "retry for error/unexpected error",
			args: args{
				attempt: 3,
				body:    "body",
				err:     fmt.Errorf("stub error"),
			},
			retryCheck: retryCheckForTimeout,
			want: result{
				shouldRetry:   false,
				weight:        0,
				pauseDuration: 0,
			},
		},
		{
			name: "retry for error/no error",
			args: args{
				attempt: 3,
				body:    "body",
				err:     nil,
			},
			retryCheck: retryCheckForTimeout,
			want: result{
				shouldRetry:   false,
				weight:        0,
				pauseDuration: 0,
			},
		},
		{
			name: "retry for status code/matched",
			args: args{
				attempt: 3,
				body:    "body",
				err:     client.StatusCodeError{StatusCode: 400},
			},
			retryCheck: retryCheckForStatusCode,
			want: result{
				shouldRetry:   true,
				weight:        20,
				pauseDuration: 2 * time.Second,
			},
		},
		{
			name: "retry for status code/unmatched",
			args: args{
				attempt: 3,
				body:    "body",
				err:     client.StatusCodeError{StatusCode: 500},
			},
			retryCheck: retryCheckForStatusCode,
			want: result{
				shouldRetry:   false,
				weight:        0,
				pauseDuration: 0,
			},
		},
		{
			name: "retry for status code/no error",
			args: args{
				attempt: 3,
				body:    "body",
				err:     nil,
			},
			retryCheck: retryCheckForStatusCode,
			want: result{
				shouldRetry:   false,
				weight:        0,
				pauseDuration: 0,
			},
		},
		{
			name: "retry for body contains/matched",
			args: args{
				attempt: 3,
				body:    "some data some data some data",
				err:     nil,
			},
			retryCheck: retryCheckForBodyContains,
			want: result{
				shouldRetry:   true,
				weight:        30,
				pauseDuration: 48 * time.Second,
			},
		},
		{
			name: "retry for body contains/not matched",
			args: args{
				attempt: 3,
				body:    "here is no matching string",
				err:     nil,
			},
			retryCheck: retryCheckForBodyContains,
			want: result{
				shouldRetry:   false,
				weight:        0,
				pauseDuration: 0,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			shouldRetry, retryWeight, pauseDuration := test.retryCheck(test.args.attempt, test.args.body, test.args.err)
			assert.Equal(t, test.want.shouldRetry, shouldRetry)
			assert.Equal(t, test.want.weight, retryWeight)
			assert.Equal(t, test.want.pauseDuration, pauseDuration)
		})
	}
}
