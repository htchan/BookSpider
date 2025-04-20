package circuitbreaker

import (
	"errors"

	client "github.com/htchan/BookSpider/internal/client/v2"
)

type FailCheck func(res string, err error) bool

func newFailCheck(conf CheckConfig) FailCheck {
	switch conf.Type {
	default:
		values := conf.Value.([]interface{})

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
			for _, code := range targetStatusCodes {
				if code == statusCodeError.StatusCode {
					return true
				}
			}
		}

		return false
	}
}
