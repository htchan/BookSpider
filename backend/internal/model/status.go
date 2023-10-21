package model

import "strings"

type StatusCode int

// TODO: update variable name to StatusXXX
const (
	Error = iota
	InProgress
	End
)

// TODO: update variable name to StatusErrorKey
const (
	ErrorKey      = "ERROR"
	InProgressKey = "INPROGRESS"
	EndKey        = "END"
)

var StatusCodeMap = map[string]StatusCode{
	ErrorKey:      Error,
	InProgressKey: InProgress,
	EndKey:        End,
}

func StatusFromString(str string) StatusCode {
	result, ok := StatusCodeMap[strings.ToUpper(str)]
	if !ok {
		return Error
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
