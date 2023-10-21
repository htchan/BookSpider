package model

import "strings"

type StatusCode int

// TODO: update variable name to StatusXXX
const (
	StatusError = iota
	StatusInProgress
	StatusEnd
)

// TODO: update variable name to StatusErrorKey
const (
	ErrorKey      = "ERROR"
	InProgressKey = "INPROGRESS"
	EndKey        = "END"
)

var StatusCodeMap = map[string]StatusCode{
	ErrorKey:      StatusError,
	InProgressKey: StatusInProgress,
	EndKey:        StatusEnd,
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
	return ErrorKey
}
