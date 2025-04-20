package retry

import (
	"errors"
	"regexp"
	"syscall"
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
)

type RetryCheck func(int, string, error) (shouldRetry bool, retryWeight int, pauseDuration time.Duration)

func CalculatePauseDuration(i int, retryInterval time.Duration, intervalType PauseIntervalType) time.Duration {
	i += 1
	switch intervalType {
	case PauseIntervalTypeConst:
		return retryInterval
	case PauseIntervalTypeLinear:
		return time.Duration(i) * retryInterval
	case PauseIntervalTypeExponential:
		return time.Duration(i*i) * retryInterval
	default:
		return 0
	}
}

func NewRetryCheck(condition RetryCondition) RetryCheck {
	switch condition.Type {
	case RetryConditionTypeStatusCode:
		values := condition.Value.([]interface{})
		statusCodes := make([]int, len(values))
		for i := range values {
			statusCodes[i] = values[i].(int)
		}

		return retryWhenStatusCodeRange(statusCodes, condition.Weight, condition.PauseInterval, condition.PauseIntervalType)
	case RetryConditionTypeTimeout:
		return retryWhenError(client.ErrTimeout, condition.Weight, condition.PauseInterval, condition.PauseIntervalType)
	case RetryConditionTypeConnectionReset:
		return retryWhenError(syscall.ECONNRESET, condition.Weight, condition.PauseInterval, condition.PauseIntervalType)
	case RetryConditionTypeBodyContains:
		target := condition.Value.(string)
		return retryWhenBodyContains(target, condition.Weight, condition.PauseInterval, condition.PauseIntervalType)
	default:
		return func(i int, s string, err error) (shouldRetry bool, retryWeight int, pauseDuration time.Duration) {
			return false, 0, 0
		}
	}
}

func retryWhenError(targetErr error, retryWeight int, retryDuration time.Duration, intervalType PauseIntervalType) RetryCheck {
	return func(i int, s string, err error) (bool, int, time.Duration) {
		if errors.Is(err, targetErr) {
			return true, retryWeight, CalculatePauseDuration(i, retryDuration, intervalType)
		}

		return false, 0, 0
	}
}

func retryWhenStatusCodeRange(statusCodes []int, retryWeight int, retryDuration time.Duration, intervalType PauseIntervalType) RetryCheck {
	return func(i int, s string, err error) (bool, int, time.Duration) {
		var statusCodeError client.StatusCodeError
		if errors.As(err, &statusCodeError) {
			for _, statusCode := range statusCodes {
				if statusCodeError.StatusCode == statusCode {
					return true, retryWeight, CalculatePauseDuration(i, retryDuration, intervalType)
				}
			}
		}

		return false, 0, 0
	}
}

func retryWhenBodyContains(target string, retryWeight int, retryDuration time.Duration, intervalType PauseIntervalType) RetryCheck {
	return func(i int, s string, err error) (bool, int, time.Duration) {
		regex := regexp.MustCompile(target)
		if regex.FindString(s) != "" {
			return true, retryWeight, CalculatePauseDuration(i, retryDuration, intervalType)
		}

		return false, 0, 0
	}
}
