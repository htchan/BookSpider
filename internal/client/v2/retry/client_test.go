package retry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
	"github.com/htchan/BookSpider/internal/client/v2/simple"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		conf *RetryClientConfig
		cli  client.BookClient
		want *RetryClient
	}{
		{
			name: "happy path",
			conf: &RetryClientConfig{
				MaxRetryWeight: 10,
				RetryConditions: []RetryCondition{
					{
						Type:              RetryConditionTypeStatusCode,
						Value:             []interface{}{404},
						Weight:            1,
						PauseInterval:     time.Second,
						PauseIntervalType: PauseIntervalTypeConst,
					},
				},
			},
			cli: nil,
			want: &RetryClient{
				c: nil,
				RetryChecks: []RetryCheck{
					retryWhenStatusCodeRange([]int{404}, 1, time.Second, PauseIntervalTypeConst),
				},
				conf: &RetryClientConfig{
					MaxRetryWeight: 10,
					RetryConditions: []RetryCondition{
						{
							Type:              RetryConditionTypeStatusCode,
							Value:             []interface{}{404},
							Weight:            1,
							PauseInterval:     time.Second,
							PauseIntervalType: PauseIntervalTypeConst,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewClient(test.conf, test.cli)
			assert.Equal(t, test.want.c, got.c)
			assert.Equal(t, len(test.want.RetryChecks), len(got.RetryChecks))
			assert.Equal(t, test.want.conf, got.conf)
		})
	}
}

func TestRetryClient_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		getContext func() context.Context
		url        string
	}

	simpleClient := simple.NewClient(&simple.SimpleClientConfig{
		RequestTimeout: 7 * time.Millisecond,
		DecodeMethod:   "",
	})

	tests := []struct {
		name                  string
		client                *RetryClient
		serverHandler         http.HandlerFunc
		args                  args
		want                  string
		wantError             error
		expectedDurationTaken time.Duration
	}{
		{
			name: "happy flow",
			client: NewClient(
				&RetryClientConfig{
					RetryConditions: []RetryCondition{},
					MaxRetryWeight:  5,
				},
				simpleClient,
			),
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello"))
			},
			args: args{
				getContext: func() context.Context { return context.Background() },
				url:        "/test",
			},
			want:                  "hello",
			wantError:             nil,
			expectedDurationTaken: 0 * time.Millisecond,
		},
		{
			name: "exceed retry count by timeout",
			client: NewClient(
				&RetryClientConfig{
					RetryConditions: []RetryCondition{
						{
							Type:              RetryConditionTypeTimeout,
							Weight:            1,
							PauseInterval:     10 * time.Millisecond,
							PauseIntervalType: PauseIntervalTypeLinear,
						},
					},
					MaxRetryWeight: 5,
				},
				simpleClient,
			),
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(9 * time.Millisecond)
				w.Write([]byte("hello"))
			},
			args: args{
				getContext: func() context.Context { return context.Background() },
				url:        "/test",
			},
			want:                  "",
			wantError:             client.ErrTimeout,
			expectedDurationTaken: 190 * time.Millisecond,
		},
		{
			name: "success before using all retry count",
			client: NewClient(
				&RetryClientConfig{
					RetryConditions: []RetryCondition{
						{
							Type:              RetryConditionTypeTimeout,
							Weight:            1,
							PauseInterval:     10 * time.Millisecond,
							PauseIntervalType: PauseIntervalTypeLinear,
						},
					},
					MaxRetryWeight: 5,
				},
				simpleClient,
			),
			serverHandler: func() http.HandlerFunc {
				var i atomic.Int32
				return func(w http.ResponseWriter, r *http.Request) {
					if i.Load() < 4 {
						time.Sleep(9 * time.Millisecond)
						i.Add(1)
					}
					w.Write([]byte("hello"))
				}
			}(),
			args: args{
				getContext: func() context.Context { return context.Background() },
				url:        "/test3",
			},
			want:                  "hello",
			wantError:             nil,
			expectedDurationTaken: 130 * time.Millisecond,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(test.serverHandler)
			defer server.Close()

			start := time.Now()
			got, err := test.client.Get(test.args.getContext(), server.URL+test.args.url)
			timeTaken := time.Since(start).Truncate(10 * time.Millisecond)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
			assert.Equal(t, test.expectedDurationTaken, timeTaken)
		})
	}
}
