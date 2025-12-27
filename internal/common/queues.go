package common

import (
	"github.com/htchan/BookSpider/internal/config/v1"
	"github.com/nats-io/nats.go"
)

func ConnectNatsQueue(conf *config.NatsConfig) (*nats.Conn, error) {
	nc, err := nats.Connect(conf.URL)
	if err != nil {
		return nil, err
	}

	return nc, err
}
