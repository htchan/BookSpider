package retry

import (
	"context"
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
	"github.com/rs/zerolog/log"
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
			if err != nil || shouldRetry {
				log.Debug().
					Str("url", url).Int("count", i).
					Bool("should_retry", shouldRetry).
					Int("retry_weight", retryWeight).
					Str("pause_duration", pauseDuration.String()).
					Msg("request result")
			}

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
