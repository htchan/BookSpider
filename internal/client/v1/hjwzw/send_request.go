package hjwzw

import (
	"context"
	"io"
	"net/http"

	"github.com/htchan/BookSpider/internal/client/v1"
	"github.com/htchan/BookSpider/internal/config/v1"
	"github.com/htchan/goclient"
	"github.com/htchan/goclient/middlewares/retry"
	pool "github.com/htchan/goclient/requester/client_pool"
)

func clientAvailable(cli *http.Client) bool {
	resp, err := cli.Get(availabilityURL())
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func newClient(ctx context.Context, conf config.ClientConfig) *goclient.Client {
	clientPool := client.NewClientPool(conf.Pool)

	go clientPool.BackgroundRefreshClients(ctx, clientAvailable)

	return goclient.NewClient(
		goclient.WithMiddlewares(
			retry.NewRetryMiddleware(
				conf.Retry.MaxRetryCount,
				retry.RetryForError,
				retry.LinearRetryInterval(conf.Retry.RetryInterval),
			),
			client.RaiseErrorForNon2xxMiddleware,
		),
		goclient.WithRequester(
			pool.NewClientPoolRequester(
				clientPool,
				clientPool.RequestRecorder,
			),
		),
	)
}

func (c *hjwzwClient) get(ctx context.Context, url string) (string, error) {
	req, reqErr := http.NewRequestWithContext(ctx, "GET", url, nil)
	if reqErr != nil {
		return "", reqErr
	}

	resp, respErr := c.cli.Do(req)
	if respErr != nil {
		return "", respErr
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}

	decodedStr, decodeErr := c.decoder.Decode(string(bodyBytes))
	if decodeErr != nil {
		return "", decodeErr
	}

	return decodedStr, nil
}
