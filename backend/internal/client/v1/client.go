package client

//go:generate mockgen -destination=../../mock/client/client.go -package=mockclient . Client
type Client interface {
	Acquire() error
	Release()
	Get(url string) (string, error)
}
