package circuitbreaker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
	"github.com/htchan/BookSpider/internal/client/v2/simple"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	type args struct {
		conf       *CircuitBreakerClientConfig
		bookClient client.BookClient
	}

	simpleClient := simple.NewClient(&simple.SimpleClientConfig{})

	var atomicStatusClosed atomic.Value
	atomicStatusClosed.Store(StatusClosed)

	tests := []struct {
		name string
		args args
		want *CircuitBreakerClient
	}{
		{
			name: "happy flow",
			args: args{
				conf: &CircuitBreakerClientConfig{
					OpenThreshold:         2,
					AcquireTimeout:        1 * time.Second,
					MaxConcurrencyThreads: 10,
					RecoverThreads:        []int64{1, 2, 3},
					OpenDuration:          1 * time.Second,
					RecoverDuration:       1 * time.Second,
					CheckConfigs: []CheckConfig{
						{
							Type:  CheckTypeStatusCodes,
							Value: []interface{}{502},
						},
					},
				},
				bookClient: simpleClient,
			},
			want: &CircuitBreakerClient{
				config: &CircuitBreakerClientConfig{
					OpenThreshold:         2,
					AcquireTimeout:        1 * time.Second,
					MaxConcurrencyThreads: 10,
					RecoverThreads:        []int64{1, 2, 3},
					OpenDuration:          1 * time.Second,
					RecoverDuration:       1 * time.Second,
					CheckConfigs: []CheckConfig{
						{
							Type:  CheckTypeStatusCodes,
							Value: []interface{}{502},
						},
					},
				},
				client:   simpleClient,
				status:   atomicStatusClosed,
				weighted: semaphore.NewWeighted(10),
				failChecks: []FailCheck{
					newStatusFailCheck([]int{502}),
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewClient(test.args.conf, test.args.bookClient)

			assert.Equal(t, test.want.config, got.config)
			assert.Equal(t, test.want.client, got.client)
			assert.Equal(t, test.want.failCount, got.failCount)
			assert.Equal(t, test.want.status, got.status)
			assert.Equal(t, test.want.halfOpenLevel, got.halfOpenLevel)
		})
	}
}

func TestCircuitBreakerClient_requestWeights(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		prepareClient func() *CircuitBreakerClient
		want          int64
	}{
		{
			name: "client status is closed",
			prepareClient: func() *CircuitBreakerClient {
				return NewClient(&CircuitBreakerClientConfig{MaxConcurrencyThreads: 10}, nil)
			},
			want: 1,
		},
		{
			name: "client status is half open/max concurrency threads is dividable by recover thread",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2, 3},
					},
					nil,
				)
				cli.status.Store(StatusHalfOpen)
				cli.halfOpenLevel.Store(1)

				return cli
			},
			want: 5,
		},
		{
			name: "client status is half open/max concurrency threads is not dividable by recover thread",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2, 3},
					},
					nil,
				)
				cli.status.Store(StatusHalfOpen)
				cli.halfOpenLevel.Store(2)

				return cli
			},
			want: 3,
		},
		{
			name: "client status is open",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(&CircuitBreakerClientConfig{MaxConcurrencyThreads: 10}, nil)
				cli.status.Store(StatusOpen)

				return cli
			},
			want: 11,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cli := test.prepareClient()
			got := cli.requestWeights()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestCircuitBreakerClient_acquire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		prepareClient func() *CircuitBreakerClient
		want          int64
		wantDuration  time.Duration
	}{
		{
			name: "acquire with closed client",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						AcquireTimeout:        25 * time.Millisecond,
					},
					nil,
				)
				cli.status.Store(StatusClosed)

				return cli
			},
			want:         1,
			wantDuration: 0 * time.Millisecond,
		},
		{
			name: "acquire with half open client",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1},
						AcquireTimeout:        25 * time.Millisecond,
					},
					nil,
				)
				cli.status.Store(StatusHalfOpen)

				return cli
			},
			want:         10,
			wantDuration: 0 * time.Millisecond,
		},
		{
			name: "acquire with open and closed client",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						AcquireTimeout:        25 * time.Millisecond,
					},
					nil,
				)
				cli.status.Store(StatusOpen)

				go func() {
					time.Sleep(70 * time.Millisecond)
					cli.status.Store(StatusClosed)
				}()

				return cli
			},
			want:         1,
			wantDuration: 75 * time.Millisecond,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cli := test.prepareClient()

			startTime := time.Now()
			got := cli.acquire(context.Background())
			assert.Equal(t, test.wantDuration, time.Since(startTime).Truncate(5*time.Millisecond))
			assert.Equal(t, test.want, got)

			acquireFailAmount := cli.config.MaxConcurrencyThreads - got + 1
			failure := cli.weighted.TryAcquire(acquireFailAmount)
			if !assert.False(t, failure) {
				cli.weighted.Release(acquireFailAmount)
			}

			// ensure the acquire amount is correct
			acquireSuccessAmount := cli.config.MaxConcurrencyThreads - got
			success := cli.weighted.TryAcquire(acquireSuccessAmount)
			if assert.True(t, success) {
				cli.weighted.Release(acquireSuccessAmount)
			}

			// release acquired semaphore
			cli.weighted.Release(got)
		})
	}
}

func TestCircuitBreakerClient_recover(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		prepareClient         func() *CircuitBreakerClient
		wantOpenDuration      time.Duration
		wantHalfOpenDurations []time.Duration
		wantClosedDuration    time.Duration
	}{
		{
			name: "happy flow",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2},
						RecoverDuration:       50 * time.Millisecond,
						OpenDuration:          20 * time.Millisecond,
					},
					nil,
				)
				cli.status.Store(StatusOpen)

				return cli
			},
			wantOpenDuration: 30 * time.Millisecond, // the open duration
			wantHalfOpenDurations: []time.Duration{
				50 * time.Millisecond, // the recover duration of first half open level
			},
			wantClosedDuration: 50 * time.Millisecond, // the recover duration of second (last) half open level
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cli := test.prepareClient()

			go cli.recover()

			time.Sleep(test.wantOpenDuration)
			assert.Equal(t, StatusHalfOpen, cli.status.Load())
			assert.Equal(t, int32(0), cli.halfOpenLevel.Load())

			for i, dur := range test.wantHalfOpenDurations {
				time.Sleep(dur)
				assert.Equal(t, StatusHalfOpen, cli.status.Load())
				assert.Equal(t, int32(i+1), cli.halfOpenLevel.Load(), "index: %d", i)
			}

			time.Sleep(test.wantClosedDuration)
			assert.Equal(t, StatusClosed, cli.status.Load())
			assert.Equal(t, int32(0), cli.halfOpenLevel.Load())
		})
	}
}

func TestCircuitBreakerClient_handleCircuitOpen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		prepareClient             func() *CircuitBreakerClient
		wantStatus                CircuitBreakerStatus
		wantHalfOpenLevel         int32
		wantFailCount             uint32
		notCloseDuration          time.Duration
		closeDuration             time.Duration
		wantFailCountAfterRecover uint32
	}{
		{
			name: "happy flow",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2},
						RecoverDuration:       50 * time.Millisecond,
						OpenDuration:          20 * time.Millisecond,
					},
					nil,
				)

				return cli
			},
			wantStatus:                StatusOpen,
			wantHalfOpenLevel:         int32(0),
			wantFailCount:             uint32(0),
			notCloseDuration:          119 * time.Millisecond,
			closeDuration:             130 * time.Millisecond,
			wantFailCountAfterRecover: 0,
		},
		{
			name: "it fail in middle",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2},
						RecoverDuration:       50 * time.Millisecond,
						OpenDuration:          20 * time.Millisecond,
					},
					nil,
				)

				go func() {
					time.Sleep(60 * time.Millisecond)
					cli.failCount.Add(1)
				}()

				return cli
			},
			wantStatus:                StatusOpen,
			wantHalfOpenLevel:         int32(0),
			wantFailCount:             uint32(0),
			notCloseDuration:          189 * time.Millisecond,
			closeDuration:             200 * time.Millisecond,
			wantFailCountAfterRecover: 0,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cli := test.prepareClient()
			cli.handleCircuitOpen()
			assert.Equal(t, test.wantStatus, cli.status.Load())
			assert.Equal(t, test.wantHalfOpenLevel, cli.halfOpenLevel.Load())
			assert.Equal(t, test.wantFailCount, cli.failCount.Load())

			time.Sleep(test.notCloseDuration)
			assert.NotEqual(t, StatusClosed, cli.status.Load())
			time.Sleep(test.closeDuration - test.notCloseDuration)
			assert.Equal(t, StatusClosed, cli.status.Load())
		})
	}
}

func TestCircuitBreakClient_reqchOpenThreshold(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		prepareClient             func() *CircuitBreakerClient
		wantStatus                CircuitBreakerStatus
		wantHalfOpenLevel         int32
		wantFailCount             uint32
		notCloseDuration          time.Duration
		closeDuration             time.Duration
		wantFailCountAfterRecover uint32
	}{
		{
			name: "not reach open threshols",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2},
						RecoverDuration:       50 * time.Millisecond,
						OpenDuration:          20 * time.Millisecond,
						OpenThreshold:         10,
					},
					nil,
				)
				cli.failCount.Add(9)

				return cli
			},
			wantStatus:        StatusClosed,
			wantHalfOpenLevel: 0,
			wantFailCount:     9,
		},
		{
			name: "set status to open when it reach fail threshold",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2},
						RecoverDuration:       50 * time.Millisecond,
						OpenDuration:          20 * time.Millisecond,
					},
					nil,
				)
				cli.failCount.Add(10)

				return cli
			},
			wantStatus:                StatusOpen,
			wantHalfOpenLevel:         int32(0),
			wantFailCount:             uint32(0),
			notCloseDuration:          119 * time.Millisecond,
			closeDuration:             130 * time.Millisecond,
			wantFailCountAfterRecover: 0,
		},
		{
			name: "it fail again in middle",
			prepareClient: func() *CircuitBreakerClient {
				cli := NewClient(
					&CircuitBreakerClientConfig{
						MaxConcurrencyThreads: 10,
						RecoverThreads:        []int64{1, 2},
						RecoverDuration:       50 * time.Millisecond,
						OpenDuration:          20 * time.Millisecond,
					},
					nil,
				)
				cli.failCount.Add(11)

				go func() {
					time.Sleep(100 * time.Millisecond)
					cli.failCount.Add(1)
				}()

				return cli
			},
			wantStatus:                StatusOpen,
			wantHalfOpenLevel:         int32(0),
			wantFailCount:             uint32(0),
			notCloseDuration:          239 * time.Millisecond,
			closeDuration:             250 * time.Millisecond,
			wantFailCountAfterRecover: 0,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cli := test.prepareClient()
			cli.reachOpenThreshold()
			assert.Equal(t, test.wantStatus, cli.status.Load())
			assert.Equal(t, test.wantHalfOpenLevel, cli.halfOpenLevel.Load())
			assert.Equal(t, test.wantFailCount, cli.failCount.Load())
			if test.wantStatus == StatusClosed {
				return
			}

			time.Sleep(test.notCloseDuration)
			assert.NotEqual(t, StatusClosed, cli.status.Load())
			time.Sleep(test.closeDuration - test.notCloseDuration)
			assert.Equal(t, StatusClosed, cli.status.Load())
			assert.Equal(t, test.wantFailCount, cli.failCount.Load())
		})
	}
}

func TestCircuitBreakerClient_Get(t *testing.T) {
	t.Parallel()

	prepareClient := func() *CircuitBreakerClient {
		return NewClient(
			&CircuitBreakerClientConfig{
				OpenThreshold:         6,
				AcquireTimeout:        10 * time.Millisecond,
				MaxConcurrencyThreads: 2,
				RecoverThreads:        []int64{1, 2},
				RecoverDuration:       100 * time.Millisecond,
				OpenDuration:          200 * time.Millisecond,
				CheckConfigs: []CheckConfig{
					{Type: CheckTypeStatusCodes, Value: []interface{}{500}},
				},
			},
			simple.NewClient(&simple.SimpleClientConfig{
				RequestTimeout: time.Second,
			}),
		)
	}

	type respFunc func(w http.ResponseWriter)
	// delaySuccessRespFunc := func(w http.ResponseWriter) { w.Write([]byte("hello")) }
	// delayFailRespFunc := func(w http.ResponseWriter) { w.WriteHeader(http.StatusInternalServerError) }
	delaySuccessRespFunc := func(w http.ResponseWriter) {
		time.Sleep(20 * time.Millisecond)
		w.Write([]byte("hello"))
	}
	delayFailRespFunc := func(w http.ResponseWriter) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusInternalServerError)
	}

	prepareServer := func(resps []respFunc) *httptest.Server {
		var i atomic.Int64
		return httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				resps[i.Load()](w)
				i.Add(1)
			},
		))
	}

	tests := []struct {
		name              string
		prepareClient     func() *CircuitBreakerClient
		resps             []respFunc
		sendNRequest      int
		url               string
		wantStatus        CircuitBreakerStatus
		wantHalfOpenLevel int32
		wantFailCount     uint32
		durationTaken     time.Duration
	}{
		{
			name:          "happy flow",
			prepareClient: prepareClient,
			resps: []respFunc{
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
			},
			sendNRequest:      20,
			url:               "/happy-flow",
			wantStatus:        StatusClosed,
			wantHalfOpenLevel: 0,
			wantFailCount:     0,
			durationTaken:     200 * time.Millisecond,
		},
		{
			name:              "fail but not reach threshold",
			prepareClient:     prepareClient,
			resps:             []respFunc{delayFailRespFunc, delayFailRespFunc, delayFailRespFunc},
			sendNRequest:      3,
			url:               "/happy-flow",
			wantStatus:        StatusClosed,
			wantHalfOpenLevel: 0,
			wantFailCount:     3,
			durationTaken:     0 * time.Millisecond,
		},
		{
			name:              "fail and reach threshold",
			prepareClient:     prepareClient,
			resps:             []respFunc{delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc},
			sendNRequest:      6,
			url:               "/happy-flow",
			wantStatus:        StatusOpen,
			wantHalfOpenLevel: 0,
			wantFailCount:     0,
			durationTaken:     50 * time.Millisecond,
		},
		{
			name:          "reach failed threshold and in half open status",
			prepareClient: prepareClient,
			resps: []respFunc{
				delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc,
				delayFailRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
			},
			sendNRequest:      9,
			url:               "/happy-flow",
			wantStatus:        StatusHalfOpen,
			wantHalfOpenLevel: 0,
			wantFailCount:     0,
			durationTaken:     250 * time.Millisecond,
		},
		{
			name:          "reach failed threshold and in 2 level of half open",
			prepareClient: prepareClient,
			resps: []respFunc{
				delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc,
				delayFailRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
			},
			sendNRequest:      13,
			url:               "/happy-flow",
			wantStatus:        StatusHalfOpen,
			wantHalfOpenLevel: 1,
			wantFailCount:     0,
			durationTaken:     350 * time.Millisecond,
		},
		{
			name:          "reach failed threshold and recover to close status",
			prepareClient: prepareClient,
			resps: []respFunc{
				delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc,
				delayFailRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
			},
			sendNRequest:      22,
			url:               "/happy-flow",
			wantStatus:        StatusClosed,
			wantHalfOpenLevel: 0,
			wantFailCount:     0,
			durationTaken:     450 * time.Millisecond,
		},
		{
			name:          "reach failed threshold and fail in recover",
			prepareClient: prepareClient,
			resps: []respFunc{
				delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc,
				delayFailRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delayFailRespFunc, delayFailRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
			},
			sendNRequest:      22,
			url:               "/happy-flow",
			wantStatus:        StatusOpen,
			wantHalfOpenLevel: 0,
			wantFailCount:     0,
			durationTaken:     450 * time.Millisecond,
		},
		{
			name:          "reach failed threshold, fail in recover and recover to close status finally",
			prepareClient: prepareClient,
			resps: []respFunc{
				delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc, delayFailRespFunc,
				delayFailRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delayFailRespFunc, delayFailRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
				delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc, delaySuccessRespFunc,
			},
			sendNRequest:      28,
			url:               "/happy-flow",
			wantStatus:        StatusClosed,
			wantHalfOpenLevel: 0,
			wantFailCount:     0,
			durationTaken:     750 * time.Millisecond,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			cli := test.prepareClient()
			srv := prepareServer(test.resps)
			defer srv.Close()

			start := time.Now()
			var wg sync.WaitGroup
			for i := 0; i < test.sendNRequest; i++ {
				wg.Add(1)
				go func() {
					cli.Get(context.Background(), srv.URL+test.url)
					wg.Done()
				}()
			}
			wg.Wait()
			timeTaken := time.Since(start).Truncate(50 * time.Millisecond)
			assert.Equal(t, test.durationTaken, timeTaken)
			assert.Equal(t, test.wantStatus, cli.status.Load())
			assert.Equal(t, test.wantHalfOpenLevel, cli.halfOpenLevel.Load())
			assert.Equal(t, test.wantFailCount, cli.failCount.Load())
		})
	}
}
