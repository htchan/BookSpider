package simple

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"

	client "github.com/htchan/BookSpider/internal/client/v2"
)

type SimpleClient struct {
	decoder client.Decoder
	client  http.Client
}

var _ client.BookClient = (*SimpleClient)(nil)

func NewClient(conf *SimpleClientConfig) *SimpleClient {
	return &SimpleClient{
		decoder: client.NewDecoder(conf.DecodeMethod),
		client:  http.Client{Timeout: conf.RequestTimeout},
	}
}

func (c *SimpleClient) Get(ctx context.Context, url string) (string, error) {
	res, reqErr := c.client.Get(url)
	if reqErr != nil {
		var timeoutError net.Error
		if errors.As(reqErr, &timeoutError); timeoutError.Timeout() {
			return "", client.ErrTimeout
		}

		return "", reqErr
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", client.StatusCodeError{StatusCode: res.StatusCode}
	}

	html, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "", readErr
	}

	result, decodeErr := c.decoder.Decode(string(html))
	if decodeErr != nil {
		return "", decodeErr
	}

	return result, nil
}
