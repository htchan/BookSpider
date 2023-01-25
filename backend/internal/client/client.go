package client

//go:generate mockgen -source=./$GOFILE -destination=../mock/$GOFILE -package=mock
type Client interface {
	Acquire() error
	Release()
	Get(url string) (string, error)
}
