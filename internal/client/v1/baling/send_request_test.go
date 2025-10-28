package baling

import (
	"net/http"
	"strings"
	"testing"
	"time"

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
				vendorHost = "www.80zw.info"
			})
		})
	}
}

func Test_newClient(t *testing.T) {
	t.Skipf("not implemented")
}

func TestClient_get(t *testing.T) {
	t.Skipf("not implemented")
}
