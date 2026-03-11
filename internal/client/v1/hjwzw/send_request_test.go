package hjwzw

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/htchan/BookSpider/internal/client/v1"
	"github.com/htchan/BookSpider/internal/config/v1"
	"github.com/htchan/goclient"
	"github.com/stretchr/testify/assert"
)

func Test_clientAvailable(t *testing.T) {
	tests := []struct {
		name string
		cli  *http.Client
		url  string
		want bool
	}{
		{
			name: "happy flow",
			cli:  http.DefaultClient,
			url:  strings.TrimLeft(serv.URL, "http://") + "/success",
			want: true,
		},
		{
			name: "unhappy flow - status not 200",
			cli:  http.DefaultClient,
			url:  strings.TrimLeft(serv.URL, "http://") + "/forbidden",
			want: false,
		},
		{
			name: "unhappy flow - timeout",
			cli: &http.Client{
				Timeout: 10 * time.Millisecond,
			},
			url:  strings.TrimLeft(serv.URL, "http://") + "/timeout",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vendorProtocol = "http"
			vendorHost = tt.url
			got := clientAvailable(tt.cli)
			assert.Equal(t, tt.want, got)

			t.Cleanup(func() {
				vendorProtocol = "https"
				vendorHost = "tw.hjwzw.com"
			})
		})
	}
}

func Test_newClient(t *testing.T) {
	tests := []struct {
		name string
		conf config.ClientConfig
	}{
		{
			name: "happy flow",
			conf: config.ClientConfig{
				Pool: config.ClientPoolConfig{
					RefreshInterval: time.Minute,
				},
				Retry: config.RetryConfig{
					MaxRetryCount:       3,
					LinearRetryInterval: time.Second,
				},
				DecodeMethod: "utf8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cli := newClient(ctx, tt.conf)
			assert.NotNil(t, cli)
		})
	}
}

func TestClient_get(t *testing.T) {
	vendorProtocol = "http"
	vendorHost = strings.TrimLeft(serv.URL, "http://")
	t.Cleanup(func() {
		vendorProtocol = "https"
		vendorHost = "tw.hjwzw.com"
	})

	tests := []struct {
		name    string
		cli     hjwzwClient
		url     string
		want    string
		wantErr string
	}{
		{
			name: "happy flow",
			cli: hjwzwClient{
				cli:     goclient.NewClient(goclient.WithMiddlewares(client.RaiseErrorForNon2xxMiddleware)),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			url:     serv.URL + "/available/success",
			want:    "黃金屋",
			wantErr: "",
		},
		{
			name: "unhappy flow - not found",
			cli: hjwzwClient{
				cli:     goclient.NewClient(goclient.WithMiddlewares(client.RaiseErrorForNon2xxMiddleware)),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			url:     serv.URL + "/not_found",
			want:    "",
			wantErr: "invalid status code: 404",
		},
		{
			name: "unhappy flow - timeout",
			cli: hjwzwClient{
				cli: goclient.NewClient(
					goclient.WithMiddlewares(client.RaiseErrorForNon2xxMiddleware),
					goclient.WithRequester(func(r *http.Request) (*http.Response, error) {
						cli := &http.Client{Timeout: time.Millisecond}
						return cli.Do(r)
					}),
				),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			url:     serv.URL + "/timeout",
			want:    "",
			wantErr: "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := tt.cli.get(ctx, tt.url)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
