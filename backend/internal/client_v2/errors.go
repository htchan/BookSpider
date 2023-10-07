package client

import (
	"errors"
	"fmt"
)

type StatusCodeError struct {
	StatusCode int
}

var _ error = StatusCodeError{}

func (err StatusCodeError) Error() string {
	return fmt.Sprintf("status code is %d", err.StatusCode)
}

var (
	ErrTimeout = errors.New("request timeout")
)
