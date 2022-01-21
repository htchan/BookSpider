package utils

import (
	"time"
)

func GenerateHash() int {
	return int(time.Now().Unix())
}