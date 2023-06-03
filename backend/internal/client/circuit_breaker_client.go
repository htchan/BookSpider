package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/htchan/BookSpider/internal/config"
	config_new "github.com/htchan/BookSpider/internal/config_new"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

var (
	ZeroLengthError = errors.New("zero length")
)

type CircuitBreakerClient struct {
	conf           config.CircuitBreakerClientConfig
	confV2         *config_new.CircuitBreakerClientConfig
	retryConf      map[string]int
	decoder        Decoder
	ctx            context.Context
	weighted       *semaphore.Weighted
	commonCtx      *context.Context
	commonWeighted *semaphore.Weighted
	failCount      int
	waitGroup      sync.WaitGroup
	client         *http.Client
}

func NewClient(conf config.CircuitBreakerClientConfig, commonWeighted *semaphore.Weighted, commonCtx *context.Context) CircuitBreakerClient {
	if commonWeighted == nil {
		commonWeighted = semaphore.NewWeighted(int64(conf.MaxThreads))
	}
	if commonCtx == nil {
		ctxObj := context.Background()
		commonCtx = &ctxObj
	}
	return CircuitBreakerClient{
		conf:           conf,
		client:         &http.Client{Timeout: time.Duration(conf.Timeout) * time.Second},
		decoder:        NewDecoder(conf.DecoderConfig),
		ctx:            context.Background(),
		weighted:       semaphore.NewWeighted(int64(conf.MaxThreads)),
		commonCtx:      commonCtx,
		commonWeighted: commonWeighted,
	}
}

func NewClientV2(conf *config_new.SiteConfig, commonWeighted *semaphore.Weighted, commonCtx *context.Context) *CircuitBreakerClient {
	if commonWeighted == nil {
		commonWeighted = semaphore.NewWeighted(int64(conf.MaxThreads))
	}
	if commonCtx == nil {
		ctxObj := context.Background()
		commonCtx = &ctxObj
	}
	return &CircuitBreakerClient{
		confV2:         &conf.CircuitBreakerConfig,
		retryConf:      conf.RetryConfig,
		client:         &http.Client{Timeout: conf.RequestTimeout},
		decoder:        NewDecoderV2(conf.DecodeMethod),
		ctx:            context.Background(),
		weighted:       semaphore.NewWeighted(int64(conf.MaxThreads)),
		commonCtx:      commonCtx,
		commonWeighted: commonWeighted,
	}
}

func (client *CircuitBreakerClient) Acquire() error {
	err := client.commonWeighted.Acquire(*client.commonCtx, 1)
	if err != nil {
		return err
	}
	return client.weighted.Acquire(client.ctx, 1)
}

func (client *CircuitBreakerClient) Release() {
	client.weighted.Release(1)
	client.commonWeighted.Release(1)
}

func (client *CircuitBreakerClient) SendRequest(url string) (string, error) {
	resp, err := client.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New(fmt.Sprintf("code %v", resp.StatusCode))
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(html) == 0 {
		return "", ZeroLengthError
	}
	result, err := client.decoder.Decode(string(html))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (client *CircuitBreakerClient) SendRequestWithCircuitBreaker(url string) (string, error) {
	client.waitGroup.Wait()
	response, err := client.SendRequest(url)
	maxFailCount, maxFailMultiplier := client.conf.MaxFailCount, client.conf.MaxFailMultiplier
	circuitBreakingSleep := time.Duration(client.conf.CircuitBreakingSleep) * time.Second
	if client.confV2 != nil {
		maxFailCount, maxFailMultiplier = client.confV2.MaxFailCount, client.confV2.MaxFailMultiplier
		circuitBreakingSleep = client.confV2.SleepInterval
	}
	if err != nil && err.Error() == "code 503" {
		client.failCount++
		if client.failCount > maxFailCount {
			client.waitGroup.Add(1)
			go func() {
				defer client.waitGroup.Done()
				time.Sleep(circuitBreakingSleep)
				if client.failCount > int(float64(maxFailCount)*maxFailMultiplier) {
					client.failCount = maxFailCount / 2
				}
			}()
		}
	} else {
		client.failCount = 0
	}
	return response, err
}

func (client *CircuitBreakerClient) Get(url string) (string, error) {
	var (
		html string
		err  error
	)
	intervalSleep := time.Duration(client.conf.IntervalSleep) * time.Second
	if client.confV2 != nil {
		intervalSleep = client.confV2.SleepInterval
	}
	for i := 0; true; i++ {
		html, err = client.SendRequestWithCircuitBreaker(url)
		if err != nil {
			if client.shouldRetryRequest(err, i) {
				log.Warn().Err(err).Str("url", url).Int("trial", i).Msg("send request failed")
				time.Sleep(time.Duration(i+1) * intervalSleep)
				continue
			} else {
				log.Error().Err(err).Str("url", url).Int("trial", i).Msg("send request failed")
				return html, err
			}
		}
		break
	}
	return html, err
}

func (client *CircuitBreakerClient) shouldRetryRequest(err error, trial int) bool {
	retryErr, retryUnavailable := client.conf.RetryErr, client.conf.RetryUnavailable

	if client.confV2 != nil {
		retryErr, retryUnavailable = client.retryConf["default"], client.retryConf["unavailable"]
	}

	if err.Error() == "code 503" || err.Error() == "code 502" {
		return trial < retryUnavailable
	}

	return trial < retryErr
}

func (client *CircuitBreakerClient) Close() error {
	client.waitGroup.Wait()
	return nil
}
