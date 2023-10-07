package simple

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	client "github.com/htchan/BookSpider/internal/client_v2"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		conf *SimpleClientConfig
		want *SimpleClient
	}{
		{
			name: "happy path",
			conf: &SimpleClientConfig{
				RequestTimeout: 1 * time.Second,
				DecodeMethod:   client.DecodeMethodBig5,
			},
			want: &SimpleClient{
				client: http.Client{
					Timeout: 1 * time.Second,
				},
				decoder: client.NewDecoder(client.DecodeMethodBig5),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewClient(test.conf)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestSimpleClient_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		client        *SimpleClient
		serverHandler http.HandlerFunc
		getContext    func() context.Context
		url           string
		want          string
		wantError     error
	}{
		{
			name: "happy path/empty decode method",
			client: NewClient(&SimpleClientConfig{
				RequestTimeout: 1 * time.Second,
				DecodeMethod:   client.DecodeMethodUTF8,
			}),
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("hello"))
			},
			getContext: func() context.Context { return context.Background() },
			url:        "/test",
			want:       "hello",
			wantError:  nil,
		},
		{
			name: "happy path/big5 decode method",
			client: NewClient(&SimpleClientConfig{
				RequestTimeout: 1 * time.Second,
				DecodeMethod:   client.DecodeMethodBig5,
			}),
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				hexByte, _ := hex.DecodeString("a440")
				w.Write(hexByte)
			},
			getContext: func() context.Context { return context.Background() },
			url:        "/test",
			want:       "一",
			wantError:  nil,
		},
		{
			name: "return status code error",
			client: NewClient(&SimpleClientConfig{
				RequestTimeout: 1 * time.Second,
				DecodeMethod:   client.DecodeMethodUTF8,
			}),
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			getContext: func() context.Context { return context.Background() },
			url:        "/test",
			want:       "",
			wantError:  client.StatusCodeError{StatusCode: http.StatusBadRequest},
		},
		{
			name: "return timeout error",
			client: NewClient(&SimpleClientConfig{
				RequestTimeout: 1 * time.Millisecond,
				DecodeMethod:   client.DecodeMethodUTF8,
			}),
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				w.WriteHeader(http.StatusBadRequest)
			},
			getContext: func() context.Context { return context.Background() },
			url:        "/test",
			want:       "",
			wantError:  client.ErrTimeout,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(test.serverHandler)
			defer server.Close()

			got, err := test.client.Get(test.getContext(), server.URL+test.url)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
