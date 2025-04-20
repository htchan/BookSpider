package simple

import (
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
)

type SimpleClientConfig struct {
	RequestTimeout time.Duration       `yaml:"request_timeout" validate:"min=1s"`
	DecodeMethod   client.DecodeMethod `yaml:"decode_method" validate:"oneof=gbk big5 utf8"`
}
