package vendor

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/htchan/goclient"
	circuitbreaker "github.com/htchan/goclient/middlewares/circuit_breaker"
	ratelimit "github.com/htchan/goclient/middlewares/rate_limit"
	"github.com/htchan/goclient/middlewares/retry"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

// ClientConfig holds the configuration for creating a BaseClient with middleware.
type ClientConfig struct {
	Name           string
	DecodeMethod   DecodeMethod
	RequestTimeout time.Duration

	RateLimitQueueSize int
	RateLimitInterval  time.Duration

	CircuitBreakerFailureThreshold int
	CircuitBreakerSuccessThreshold int
	CircuitBreakerRecoverDuration  time.Duration
	CircuitBreakerOpenQueueRatio   float64

	RetryMaxRetries   int
	RetryBaseInterval time.Duration
	RetryIntervalType string
}

func isServerError(req *http.Request, resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return resp != nil && resp.StatusCode >= 500
}

func retryOnError(req *http.Request, resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp != nil && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return true
	}
	return false
}

func newRetryIntervalCalculator(intervalType string, baseInterval time.Duration) retry.RetryIntervalCalculator {
	switch intervalType {
	case "linear":
		return retry.LinearRetryInterval(baseInterval)
	case "exponential":
		return retry.ExponentialRetryInterval(baseInterval)
	default:
		return retry.StaticRetryInterval(baseInterval)
	}
}

// NewBaseClientFromConfig creates a BaseClient with full middleware stack (rate limiting,
// circuit breaker, retry) and returns both the client and the vendor semaphore.
func NewBaseClientFromConfig(conf ClientConfig) (*BaseClient, *semaphore.Weighted) {
	queueSize := int64(conf.RateLimitQueueSize)
	queue := ratelimit.NewQueue(conf.RateLimitQueueSize)
	vendorSema := semaphore.NewWeighted(queueSize)

	var semaHeld atomic.Int64
	var acquireCancel context.CancelFunc

	breaker := circuitbreaker.NewCircuitBreaker(
		conf.CircuitBreakerFailureThreshold,
		conf.CircuitBreakerSuccessThreshold,
		conf.CircuitBreakerRecoverDuration,
		isServerError,
		circuitbreaker.WithOnStateChange(func(from, to circuitbreaker.State) {
			if acquireCancel != nil {
				acquireCancel()
				acquireCancel = nil
			}

			switch to {
			case circuitbreaker.StateOpen:
				newSize := int(float64(conf.RateLimitQueueSize) * conf.CircuitBreakerOpenQueueRatio)
				queue.Resize(newSize)

				var acquireCtx context.Context
				acquireCtx, acquireCancel = context.WithCancel(context.Background())

				go func() {
					for range queueSize {
						if err := vendorSema.Acquire(acquireCtx, 1); err != nil {
							return
						}
						semaHeld.Add(1)
					}
				}()
			case circuitbreaker.StateHalfOpen:
				allowedSlots := max(int64(float64(queueSize)*conf.CircuitBreakerOpenQueueRatio), 1)

				held := semaHeld.Load()
				toRelease := held - (queueSize - allowedSlots)
				if toRelease > 0 && toRelease <= held {
					vendorSema.Release(toRelease)
					semaHeld.Add(-toRelease)
				}
			case circuitbreaker.StateClosed:
				held := semaHeld.Load()
				if held > 0 {
					vendorSema.Release(held)
					semaHeld.Add(-held)
				}
				queue.Resize(conf.RateLimitQueueSize)
			}

			log.Info().
				Str("vendor", conf.Name).
				Str("from", from.String()).
				Str("to", to.String()).
				Int64("sema_held", semaHeld.Load()).
				Int64("sema_size", queueSize).
				Msg("circuit breaker state changed")
		}),
	)

	cli := goclient.NewClient(
		goclient.WithRequester((&http.Client{Timeout: conf.RequestTimeout}).Do),
		goclient.WithMiddlewares(
			retry.NewRetryMiddleware(
				conf.RetryMaxRetries,
				retryOnError,
				newRetryIntervalCalculator(
					conf.RetryIntervalType,
					conf.RetryBaseInterval,
				),
			),
			circuitbreaker.NewCircuitBreakerMiddleware(breaker),
			ratelimit.NewRateLimitMiddleware(queue, conf.RateLimitInterval),
		),
	)

	return NewBaseClient(cli, conf.DecodeMethod), vendorSema
}
