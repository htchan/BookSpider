package client

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"

	"github.com/htchan/goclient"
)

type Client struct {
	decoder Decoder
	cli     *goclient.Client
}

var _ BookClient = (*Client)(nil)

func NewClient(cli *goclient.Client, decodeMethod DecodeMethod) *Client {
	return &Client{
		decoder: NewDecoder(decodeMethod),
		cli:     cli,
	}
}

func (c *Client) Get(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		var timeoutError net.Error
		if errors.As(err, &timeoutError) && timeoutError.Timeout() {
			return "", ErrTimeout
		}

		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", StatusCodeError{StatusCode: resp.StatusCode}
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result, err := c.decoder.Decode(string(html))
	if err != nil {
		return "", err
	}

	return result, nil
}
