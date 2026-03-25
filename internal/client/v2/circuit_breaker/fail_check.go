package circuitbreaker

import (
	"errors"
	"slices"

	client "github.com/htchan/BookSpider/internal/client/v2"
)

type FailCheck func(res string, err error) bool

func newFailCheck(conf CheckConfig) FailCheck {
	switch conf.Type {
	default:
		values := conf.Value.([]any)

		statusCodes := make([]int, len(values))
		for i := range values {
			statusCodes[i] = values[i].(int)
		}

		return newStatusFailCheck(statusCodes)
	}
}

func newStatusFailCheck(targetStatusCodes []int) FailCheck {
	return func(res string, err error) bool {
		var statusCodeError client.StatusCodeError
		if errors.As(err, &statusCodeError) {
			if slices.Contains(targetStatusCodes, statusCodeError.StatusCode) {
				return true
			}
		}

		return false
	}
}
