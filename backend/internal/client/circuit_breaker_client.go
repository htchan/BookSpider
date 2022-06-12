package client

import (
	"time"
	"net/http"
	"sync"
	"errors"
	"fmt"
	"io"
	"context"
	"golang.org/x/sync/semaphore"
	"github.com/htchan/BookSpider/internal/config"
)

type CircuitBreakerClient struct {
	config.CircuitBreakerConfig
	decoder Decoder
	ctx context.Context
	weighted *semaphore.Weighted
	failCount int
	waitGroup sync.WaitGroup
	client *http.Client
}

func (client *CircuitBreakerClient) Init(maxThreads int) {
	client.client = &http.Client{Timeout: time.Duration(client.Timeout) * time.Second}
	client.decoder.Load()
	client.ctx = context.Background()
	client.weighted = semaphore.NewWeighted(int64(maxThreads))
}

func (client CircuitBreakerClient) Acquire() error {
	return client.weighted.Acquire(client.ctx, 1)
}

func (client CircuitBreakerClient) Release() {
	client.weighted.Release(1)
}

func (client CircuitBreakerClient) SendRequest(url string) (string, error) {
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
		return "", errors.New("zero length")
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
	if err != nil && err.Error() == "code 503" {
		client.failCount++
		if client.failCount > client.MaxFailCount {
			client.waitGroup.Add(1)
			go func() {
				defer client.waitGroup.Done()
				time.Sleep(time.Duration(client.CircuitBreakingSleep) * time.Second)
				if client.failCount > int(float64(client.MaxFailCount) * client.MaxFailMultiplier) {
					client.failCount = client.MaxFailCount / 2
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
		err error
	)
	for i := 0; true; i++ {
		html, err = client.SendRequestWithCircuitBreaker(url)
		if err != nil {
			if err.Error() == "code 503" && i >= client.Retry503 {
				return html, err
			} else if i >= client.RetryErr {
				return html, err
			}
			time.Sleep(time.Duration((i + 1) * client.IntervalSleep) * time.Second)
			continue
		}
		break
	}
	return html, err
}