package client

import "context"

//go:generate mockgen -destination=../mock/client/v2/book_client.go -package=mockclient . BookClient
type BookClient interface {
	Get(ctx context.Context, url string) (string, error)
}
