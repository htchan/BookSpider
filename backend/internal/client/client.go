package client

type Client interface {
	Acquire() error
	Release()
	Get(url string) (string, error)
}
