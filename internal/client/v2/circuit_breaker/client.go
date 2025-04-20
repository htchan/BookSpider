package circuitbreaker

import (
	"context"
	"sync/atomic"
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
	"golang.org/x/sync/semaphore"
)

const WeightedCount = 10

type CircuitBreakerStatus string

const (
	StatusOpen     CircuitBreakerStatus = "open"
	StatusHalfOpen CircuitBreakerStatus = "half-open"
	StatusClosed   CircuitBreakerStatus = "closed"
)

type CircuitBreakerClient struct {
	config *CircuitBreakerClientConfig
	client client.BookClient
	// circuit breaker
	failCount     atomic.Uint32
	weighted      *semaphore.Weighted
	status        atomic.Value
	halfOpenLevel atomic.Int32
	failChecks    []FailCheck
}

var _ client.BookClient = (*CircuitBreakerClient)(nil)

func NewClient(
	conf *CircuitBreakerClientConfig,
	bookClient client.BookClient,
) *CircuitBreakerClient {
	c := &CircuitBreakerClient{
		config:   conf,
		client:   bookClient,
		weighted: semaphore.NewWeighted(conf.MaxConcurrencyThreads),
	}
	c.status.Store(StatusClosed)
	c.halfOpenLevel.Store(0)
	c.failCount.Store(0)
	for _, checkConf := range conf.CheckConfigs {
		c.failChecks = append(c.failChecks, newFailCheck(checkConf))
	}

	return c
}

func (c *CircuitBreakerClient) requestWeights() int64 {
	var weight int64

	switch c.status.Load() {
	case StatusOpen:
		weight = c.config.MaxConcurrencyThreads + 1
	case StatusHalfOpen:
		halfOpenLevel := c.halfOpenLevel.Load()
		weight = c.config.MaxConcurrencyThreads / c.config.RecoverThreads[halfOpenLevel]
	default:
		weight = 1
	}

	if weight == 0 {
		return 1
	}

	return weight
}

func (c *CircuitBreakerClient) acquire(ctx context.Context) int64 {
	for {
		weight := c.requestWeights()

		err := func() error {
			ctxTimeout, cancel := context.WithTimeout(ctx, c.config.AcquireTimeout)
			defer cancel()
			return c.weighted.Acquire(ctxTimeout, weight)
		}()
		if err == nil {
			return weight
		}
	}
}

func (c *CircuitBreakerClient) recover() {
	// keep status open for a duration
	time.Sleep(c.config.OpenDuration)
	// assumed no request is on the flight
	ok := c.status.CompareAndSwap(StatusOpen, StatusHalfOpen)
	if !ok {
		return
	}
	c.halfOpenLevel.Store(0)
	// clear the fail count
	c.failCount.Store(0)

	for i := range c.config.RecoverThreads {
		time.Sleep(c.config.RecoverDuration)

		// set status to open again if there is any failure during the recover
		if c.failCount.Load() > 0 {
			c.handleCircuitOpen()
			return
		}

		ok := c.halfOpenLevel.CompareAndSwap(int32(i), int32(i+1))
		if !ok {
			return
		}
	}

	ok = c.status.CompareAndSwap(StatusHalfOpen, StatusClosed)
	if ok {
		c.halfOpenLevel.Store(0)
	}
}

func (c *CircuitBreakerClient) handleCircuitOpen() {
	// set status to open
	c.status.Store(StatusOpen)
	c.halfOpenLevel.Store(0)
	// deploy go routine to delay recover
	go c.recover()
	c.failCount.Store(0)
}

func (c *CircuitBreakerClient) reachOpenThreshold() {
	failCount := c.failCount.Load()
	if failCount >= c.config.OpenThreshold {
		c.handleCircuitOpen()
	}
}

func (c *CircuitBreakerClient) Get(ctx context.Context, url string) (string, error) {
	acquireAmount := c.acquire(ctx)
	defer func() {
		c.weighted.Release(acquireAmount)
	}()

	res, reqErr := c.client.Get(ctx, url)
	isRequestFail := false
	for _, check := range c.failChecks {
		if check(res, reqErr) {
			isRequestFail = true
			break
		}
	}
	if isRequestFail {
		c.failCount.Add(1)
		c.reachOpenThreshold()
	} else if c.status.Load() != StatusHalfOpen {
		c.failCount.Store(0)
	}

	return res, reqErr
}
