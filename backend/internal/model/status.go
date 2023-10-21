package model

import "strings"

type StatusCode int

const (
	StatusError = iota
	StatusInProgress
	StatusEnd
)

const (
	StatusErrorKey      = "ERROR"
	StatusInProgressKey = "INPROGRESS"
	StatusEndKey        = "END"
)

var StatusCodeMap = map[string]StatusCode{
	StatusErrorKey:      StatusError,
	StatusInProgressKey: StatusInProgress,
	StatusEndKey:        StatusEnd,
}

func StatusFromString(str string) StatusCode {
	result, ok := StatusCodeMap[strings.ToUpper(str)]
	if !ok {
		return StatusError
	}
	return result
}

func (status StatusCode) String() string {
	for key, value := range StatusCodeMap {
		if value == status {
			return key
		}
	}
	return StatusErrorKey
}
