package client

import "context"

//go:generate mockgen -source=./$GOFILE -destination=../mock/bookclient.go -package=mock
type BookClient interface {
	Get(ctx context.Context, url string) (string, error)
}
