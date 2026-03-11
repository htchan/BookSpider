package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/htchan/BookSpider/internal/client/v1"
	"github.com/htchan/BookSpider/internal/client/v1/baling"
	"github.com/htchan/BookSpider/internal/client/v1/bestory"
	"github.com/htchan/BookSpider/internal/client/v1/ck101"
	"github.com/htchan/BookSpider/internal/client/v1/hjwzw"
	"github.com/htchan/BookSpider/internal/client/v1/uukanshu"
	"github.com/htchan/BookSpider/internal/client/v1/xbiquge"
	"github.com/htchan/BookSpider/internal/client/v1/xqishu"
	"github.com/htchan/BookSpider/internal/config/v1"
)

var UnknownClientError = errors.New("unknown client site")

func loadClient(ctx context.Context, site string, conf config.ClientConfig) (client.Client, error) {
	switch site {
	case "baling":
		return baling.NewClient(ctx, conf), nil
	case "bestory":
		return bestory.NewClient(ctx, conf), nil
	case "ck101":
		return ck101.NewClient(ctx, conf), nil
	case "hjwzw":
		return hjwzw.NewClient(ctx, conf), nil
	case "uukanshu":
		return uukanshu.NewClient(ctx, conf), nil
	case "xbiquge":
		return xbiquge.NewClient(ctx, conf), nil
	case "xqishu":
		return xqishu.NewClient(ctx, conf), nil
	default:
		return nil, UnknownClientError
	}
}

func LoadClients(ctx context.Context, configs map[string]config.ClientConfig) (map[string]client.Client, error) {
	clients := make(map[string]client.Client)

	var err error
	for key, conf := range configs {
		client, clientErr := loadClient(ctx, key, conf)
		if clientErr != nil {
			err = errors.Join(err, fmt.Errorf("%v: %s", clientErr, key))
		}
		clients[key] = client
	}

	return clients, err
}
