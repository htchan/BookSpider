package client

import (
	"testing"
	"time"
	"strings"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/config"
)

func TestCircuitBreakerClient_Init(t *testing.T) {
	t.Parallel()
	client := CircuitBreakerClient{
		CircuitBreakerConfig: config.CircuitBreakerConfig{Timeout: 10},
	}
	if client.client != nil {
		t.Errorf("the client is not nil in default")
		return
	}
	client.Init(0)

	if client.client == nil || client.client.Timeout != 10 * time.Second {
		t.Errorf("wrong client created: %v", client.client)
	}
}

func TestCircuitBreakerClient_AcquireRelease(t *testing.T) {
	t.Parallel()

	client := CircuitBreakerClient{}
	client.Init(1)

	t.Run("acquire block other acquire until release", func (t *testing.T) {
		t.Parallel()
		registerPoint := time.Now()
		go func() {
			defer client.Release()
			client.Acquire()
			time.Sleep(1 * time.Second)
		}()
		time.Sleep(10 * time.Millisecond)
		client.Acquire()
		client.Release()
		if time.Now().Before(registerPoint.Add(1 * time.Second)) {
			t.Errorf(
				"acquire takes %v millisecond to process",
				time.Now().UnixMilli() - registerPoint.UnixMilli(),
			)
		}
	})
}

func TestCircuitBreakerClient_SendRequest(t *testing.T) {
	t.Parallel()
	server := mock.MockCircuitBreakerServer(2)
	client := CircuitBreakerClient{
		CircuitBreakerConfig: config.CircuitBreakerConfig{Timeout: 1},
	}
	client.Init(0)

	t.Cleanup( func () {
		server.Close()
	})

	t.Run("get 200", func (t *testing.T) {
		t.Parallel()
		resp, err := client.SendRequest(server.URL + "/200")
		if err != nil || resp != "200" {
			t.Errorf("send request return error: %v; response: %v", err, resp)
		}
	})

	t.Run("get 400", func (t *testing.T) {
		t.Parallel()
		resp, err := client.SendRequest(server.URL + "/400")
		if err == nil || err.Error() != "code 400" {
			t.Errorf("send request return error: %v; response: %v", err, resp)
		}
	})

	t.Run("get 503", func (t *testing.T) {
		t.Parallel()
		resp, err := client.SendRequest(server.URL + "/503")
		if err == nil || err.Error() != "code 503" {
			t.Errorf("send request return error: %v; response: %v", err, resp)
		}
	})

	t.Run("get zero length", func (t *testing.T) {
		t.Parallel()
		resp, err := client.SendRequest(server.URL + "/empty")
		if err == nil || err.Error() != "zero length" {
			t.Errorf("send request return error: %v; response: %v", err, resp)
		}
	})

	t.Run("get timeout", func (t *testing.T) {
		t.Parallel()
		resp, err := client.SendRequest(server.URL + "/timeout")
		if err == nil || !strings.Contains(err.Error(), "timeout") {
			t.Errorf("send request return error: %v; response: %v", err, resp)
		}
	})

	t.Run("get connection reject", func (t *testing.T) {
		t.Parallel()
		resp, err := client.SendRequest("http://127.0.0.1:39999")
		if err == nil || !strings.Contains(err.Error(), "connection refused") {
			t.Errorf("send request return error: %v; response: %v", err, resp)
		}
	})
}

func TestCircuitBreakerClient_SendRequestWithCircuitBreaker(t *testing.T) {
	t.Parallel()
	server := mock.MockCircuitBreakerServer(2)
	t.Cleanup( func () {
		server.Close()
	})

	client := CircuitBreakerClient{
		CircuitBreakerConfig: config.CircuitBreakerConfig{
			Timeout: 1,
			MaxFailCount: 1,
			CircuitBreakingSleep: 2,
		},
	}
	client.Init(0)

	t.Run("stop send request if exceed limit", func (t *testing.T) {
		tempClient := client
		t.Parallel()
		registerPoint := time.Now()
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/503")
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/503")
		if time.Now().After(registerPoint.Add(1 * time.Second)) {
			t.Errorf("send request with circuit breaker break request sending after 2 request was made")
		}
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/503")
		if time.Now().Before(registerPoint.Add(2 * time.Second)) {
			t.Errorf("time different (unix milli): %v", time.Now().UnixMilli() - registerPoint.UnixMilli())
			t.Errorf("send request with circuit breaker does not break request sending after 2 request was made")
		}
	})

	t.Run("reset fail count if it receive success", func (t *testing.T) {
		tempClient := client
		t.Parallel()
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/503")
		if tempClient.failCount == 0 {
			t.Errorf("send request with circuit breaker does not update fail count when if get 503")
		}
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/200")
		if tempClient.failCount != 0 {
			t.Errorf("send request with circuit breaker does not reset fail count when if get 503")
		}
	})

	t.Run("reset to half of Max if it exceed max fail count * max fail multiplier", func (t *testing.T) {
		tempClient := client
		tempClient.MaxFailCount = 2
		tempClient.MaxFailMultiplier = 1
		tempClient.CircuitBreakingSleep = 0
		t.Parallel()
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/503")
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/503")
		if tempClient.failCount != 2 {
			t.Errorf("send request with circuit breaker does not update fail count when if get 503")
		}
		tempClient.SendRequestWithCircuitBreaker(server.URL + "/503")
		time.Sleep(1 * time.Millisecond)
		if tempClient.failCount != 1 {
			t.Error(tempClient.failCount)
			t.Errorf("send request with circuit breaker does not half fail count when if get 503")
		}
	})
}

func TestCircuitBreakerClient_Get(t *testing.T) {
	t.Parallel()
	
	server := mock.MockCircuitBreakerServer(1)

	t.Cleanup(func () {
		server.Close()
	})

	client := CircuitBreakerClient{
		CircuitBreakerConfig: config.CircuitBreakerConfig{
			Retry503: 2,
			RetryErr: 3,
			IntervalSleep: 1,
		},
	}
	client.Init(0)

	t.Run("return html if receive 200", func (t *testing.T) {
		t.Parallel()
		html, err := client.Get(server.URL + "/200")
		if err != nil || html != "200" {
			t.Errorf("get return html: %v, err: %v", html, err)
		}
	})

	t.Run("retry until reach Retry503 limit if getting 503 status code", func (t *testing.T) {
		t.Parallel()
		registerPoint := time.Now()
		html, err := client.Get(server.URL + "/503")
		if err == nil || html != "" {
			t.Errorf("get return html: %v, err: %v", html, err)
		}
		if time.Now().Before(registerPoint.Add(3 * time.Second)) {
			t.Errorf(
				"get takes %v millisecond",
				time.Now().UnixMilli() - registerPoint.UnixMilli(),
			)
		}
	})

	t.Run("retry until reach RetryErr limit if getting non 503 status code", func (t *testing.T) {
		t.Parallel()
		registerPoint := time.Now()
		html, err := client.Get(server.URL + "/400")
		if err == nil || html != "" {
			t.Errorf("get return html: %v, err: %v", html, err)
		}
		if time.Now().Before(registerPoint.Add(6 * time.Second)) {
			t.Errorf(
				"get takes %v millisecond",
				time.Now().UnixMilli() - registerPoint.UnixMilli(),
			)
		}
	})
}