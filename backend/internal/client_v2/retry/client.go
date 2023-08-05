package retry

import (
	"context"
	"fmt"
	"time"

	client "github.com/htchan/BookSpider/internal/client_v2"
)

type RetryClient struct {
	c           client.BookClient
	RetryChecks []RetryCheck
	conf        *RetryClientConfig
}

var _ client.BookClient = (*RetryClient)(nil)

func NewClient(conf *RetryClientConfig, bookClient client.BookClient) *RetryClient {
	c := &RetryClient{
		c:    bookClient,
		conf: conf,
	}

	c.RetryChecks = make([]RetryCheck, 0, len(conf.RetryConditions))
	for _, retryCondition := range conf.RetryConditions {
		c.RetryChecks = append(c.RetryChecks, NewRetryCheck(retryCondition))
	}

	return c
}

func (c *RetryClient) Get(ctx context.Context, url string) (string, error) {
	var (
		retryWeight = 0
		body        string
		err         error
	)

	for i := 0; retryWeight < c.conf.MaxRetryWeight; i++ {
		body, err = c.c.Get(ctx, url)

		var (
			shouldRetry   bool
			weight        int
			pauseDuration time.Duration
		)
		for _, check := range c.RetryChecks {
			shouldRetry, weight, pauseDuration = check(i, body, err)
			fmt.Println(url, i, shouldRetry, weight, pauseDuration)
			if shouldRetry {
				retryWeight += weight
				time.Sleep(pauseDuration)

				break
			}
		}

		// shop retrying if no retry check report the result need to be retried
		if !shouldRetry {
			break
		}
	}

	return body, err
}
