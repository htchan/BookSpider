package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/htchan/BookSpider/internal/config/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewClientPool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *ClientPool
	}{
		{
			name: "happy flow",
			want: &ClientPool{
				lock:       new(sync.Mutex),
				cond:       sync.NewCond(new(sync.Mutex)),
				clients:    []*http.Client{},
				failureMap: make(map[string]int),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewClientPool(config.ClientPoolConfig{})
			assert.Equal(t, tt.want.clients, got.clients)
			assert.Equal(t, tt.want.failureMap, got.failureMap)
			assert.Equal(t, tt.want.clientInPool, got.clientInPool)
		})
	}
}

func TestClientPool_GetClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		getClientPool  func() *ClientPool
		want           *http.Client
		wantClientPool *ClientPool
	}{
		{
			name: "happy flow",
			getClientPool: func() *ClientPool {
				pool := NewClientPool(config.ClientPoolConfig{})
				pool.AddClients(&http.Client{})

				return pool
			},
			want: &http.Client{},
			wantClientPool: &ClientPool{
				lock:         new(sync.Mutex),
				cond:         sync.NewCond(new(sync.Mutex)),
				clients:      []*http.Client{},
				failureMap:   map[string]int{},
				clientInPool: []string{""},
			},
		},
		{
			name: "empty pool - add client later",
			getClientPool: func() *ClientPool {
				pool := NewClientPool(config.ClientPoolConfig{})
				go func() {
					time.Sleep(10 * time.Millisecond)
					pool.AddClients(&http.Client{})
				}()

				return pool
			},
			want: &http.Client{},
			wantClientPool: &ClientPool{
				lock:         new(sync.Mutex),
				cond:         sync.NewCond(new(sync.Mutex)),
				clients:      []*http.Client{},
				failureMap:   map[string]int{},
				clientInPool: []string{""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cp := tt.getClientPool()
			got := cp.GetClient(nil)
			assert.Equal(t, tt.want, got)
			cp.lock.Lock()
			assert.Equal(t, tt.wantClientPool.clients, cp.clients)
			assert.Equal(t, tt.wantClientPool.failureMap, cp.failureMap)
			assert.Equal(t, tt.wantClientPool.clientInPool, cp.clientInPool)
			cp.lock.Unlock()
		})
	}
}

func TestClientPool_AddClients(t *testing.T) {
	t.Parallel()

	client1 := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Path: "proxy_1"})}}
	client2 := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(&url.URL{Path: "proxy_2"})}}

	tests := []struct {
		name             string
		clientsToAdd     []*http.Client
		wantClients      []*http.Client
		wantClientInPool []string
	}{
		{
			name:             "happy flow",
			clientsToAdd:     []*http.Client{client1, client2},
			wantClients:      []*http.Client{client1, client2},
			wantClientInPool: []string{"proxy_1", "proxy_2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pool := NewClientPool(config.ClientPoolConfig{})
			pool.AddClients(tt.clientsToAdd...)
			assert.Equal(t, tt.wantClients, pool.clients)
			assert.Equal(t, tt.wantClientInPool, pool.clientInPool)
		})
	}
}

func TestClientPool_addAvailableProxyClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		proxyAddr        string
		clientAvailable  func(*http.Client) bool
		wantClientLength int
		wantClientInPool []string
	}{
		{
			name:      "happy flow - available client",
			proxyAddr: "http://proxy_1",
			clientAvailable: func(cli *http.Client) bool {
				return true
			},
			wantClientLength: 1,
			wantClientInPool: []string{"http://proxy_1"},
		},
		{
			name:      "happy flow - unavailable client",
			proxyAddr: "http://proxy_2",
			clientAvailable: func(cli *http.Client) bool {
				return false
			},
			wantClientLength: 0,
			wantClientInPool: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pool := NewClientPool(config.ClientPoolConfig{})
			pool.addAvailableProxyClient(tt.proxyAddr, tt.clientAvailable)
			assert.Equal(t, tt.wantClientLength, len(pool.clients))
			assert.Equal(t, tt.wantClientInPool, pool.clientInPool)
		})
	}
}

func TestClientPool_socks5ProxyList(t *testing.T) {
	t.Parallel()

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.Write([]byte("1.2.3.4\n5.6.7.8"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))

	t.Cleanup(func() {
		serv.Close()
	})

	tests := []struct {
		name      string
		protocol  string
		sourceURL string
		want      []string
	}{
		{
			name:      "happy flow",
			protocol:  "socks5",
			sourceURL: serv.URL + "/success",
			want:      []string{"socks5://1.2.3.4", "socks5://5.6.7.8"},
		},
		{
			name:      "unhappy flow - bad response",
			protocol:  "socks5",
			sourceURL: serv.URL + "/bad_response",
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pool := NewClientPool(config.ClientPoolConfig{})
			got := pool.proxyList(tt.protocol, tt.sourceURL)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientPool_BackgroundRefreshClients(t *testing.T) {
	t.Skipf("not implemented")
}

func TestClientPool_RequestRecorder(t *testing.T) {
	t.Parallel()

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.Write([]byte("success"))
		} else if r.URL.Path == "/delay" {
			time.Sleep(50 * time.Millisecond)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}))

	t.Cleanup(func() {
		serv.Close()
	})

	tests := []struct {
		name              string
		getPool           func() *ClientPool
		sendReq           func(*testing.T, *ClientPool) (*http.Client, *http.Request, *http.Response, error)
		wantClientsLength int
		wantClientInPool  []string
		wantFailureMap    map[string]int
	}{
		{
			name: "happy flow - successful request",
			getPool: func() *ClientPool {
				pool := NewClientPool(config.ClientPoolConfig{})
				pool.AddClients(&http.Client{})

				return pool
			},
			sendReq: func(t *testing.T, pool *ClientPool) (*http.Client, *http.Request, *http.Response, error) {
				cli := pool.GetClient(nil)

				req, err := http.NewRequest("GET", serv.URL+"/success", nil)
				assert.NoError(t, err)

				resp, err := cli.Do(req)
				return cli, req, resp, err
			},
			wantClientsLength: 1,
			wantClientInPool:  []string{""},
			wantFailureMap:    map[string]int{"": 0},
		},
		{
			name: "happy flow - failure request",
			getPool: func() *ClientPool {
				pool := NewClientPool(config.ClientPoolConfig{DropClientFailureThreshold: 2})
				pool.AddClients(&http.Client{})

				return pool
			},
			sendReq: func(t *testing.T, pool *ClientPool) (*http.Client, *http.Request, *http.Response, error) {
				cli := pool.GetClient(nil)

				req, err := http.NewRequest("GET", serv.URL+"/failure", nil)
				assert.NoError(t, err)

				resp, err := cli.Do(req)
				return cli, req, resp, err
			},
			wantClientsLength: 1,
			wantClientInPool:  []string{""},
			wantFailureMap:    map[string]int{"": 1},
		},
		{
			name: "happy flow - failure request reaching threshold",
			getPool: func() *ClientPool {
				pool := NewClientPool(config.ClientPoolConfig{DropClientFailureThreshold: 1})
				pool.failureMap[""] = 1
				pool.AddClients(&http.Client{})

				return pool
			},
			sendReq: func(t *testing.T, pool *ClientPool) (*http.Client, *http.Request, *http.Response, error) {
				cli := pool.GetClient(nil)

				req, err := http.NewRequest("GET", serv.URL+"/failure", nil)
				assert.NoError(t, err)

				resp, err := cli.Do(req)
				return cli, req, resp, err
			},
			wantClientsLength: 0,
			wantClientInPool:  []string{},
			wantFailureMap:    map[string]int{},
		},
		{
			name: "happy flow - got timeout",
			getPool: func() *ClientPool {
				pool := NewClientPool(config.ClientPoolConfig{DropClientFailureThreshold: 1})
				pool.failureMap[""] = 1
				pool.AddClients(&http.Client{Timeout: time.Millisecond})

				return pool
			},
			sendReq: func(t *testing.T, pool *ClientPool) (*http.Client, *http.Request, *http.Response, error) {
				cli := pool.GetClient(nil)

				req, err := http.NewRequest("GET", serv.URL+"/delay", nil)
				assert.NoError(t, err)

				resp, err := cli.Do(req)
				return cli, req, resp, err
			},
			wantClientsLength: 0,
			wantClientInPool:  []string{},
			wantFailureMap:    map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pool := tt.getPool()
			cli, req, resp, err := tt.sendReq(t, pool)
			pool.RequestRecorder(pool, cli, req, resp, err)
			time.Sleep(time.Millisecond)
			pool.lock.Lock()
			assert.Equal(t, tt.wantClientsLength, len(pool.clients))
			assert.Equal(t, tt.wantClientInPool, pool.clientInPool)
			assert.Equal(t, tt.wantFailureMap, pool.failureMap)
			pool.lock.Unlock()
		})
	}
}

func TestRaiseErrorForNon2xxMiddleware(t *testing.T) {
	t.Parallel()

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.Write([]byte("success"))
		} else if r.URL.Path == "/delay" {
			time.Sleep(50 * time.Millisecond)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}))

	t.Cleanup(func() {
		serv.Close()
	})

	tests := []struct {
		name         string
		url          string
		wantErrExist bool
	}{
		{
			name:         "happy flow - 2xx response",
			url:          serv.URL + "/success",
			wantErrExist: false,
		},
		{
			name:         "unhappy flow - non-2xx response",
			url:          serv.URL + "/failure",
			wantErrExist: true,
		},
		{
			name:         "unhappy flow - timeout",
			url:          serv.URL + "/delay",
			wantErrExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli := &http.Client{Timeout: time.Millisecond * 10}
			middleware := RaiseErrorForNon2xxMiddleware(cli.Do)
			req, err := http.NewRequest("GET", tt.url, nil)
			assert.NoError(t, err)

			_, err = middleware(req)
			assert.Equal(t, tt.wantErrExist, err != nil)
		})
	}
}

func Test_getAddrFromClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		getClient func(t *testing.T) *http.Client
		want      string
	}{
		{
			name: "happy flow - with proxy",
			getClient: func(t *testing.T) *http.Client {
				proxyURL, err := url.Parse("sock4:google.com")
				assert.NoError(t, err)
				return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
			},
			want: "sock4:google.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cli := tt.getClient(t)
			got := getAddrFromClient(cli)
			assert.Equal(t, tt.want, got)
		})
	}
}
