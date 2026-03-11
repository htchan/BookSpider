package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/htchan/BookSpider/internal/config/v1"
	"github.com/htchan/goclient"
	pool "github.com/htchan/goclient/requester/client_pool"
)

var _ (pool.ClientPool) = (*ClientPool)(nil)

var _ pool.RequestRecorder

type ClientPool struct {
	clients      []*http.Client
	cond         *sync.Cond
	lock         *sync.Mutex
	failureMap   map[string]int
	clientInPool []string
	conf         config.ClientPoolConfig
}

func NewClientPool(conf config.ClientPoolConfig) *ClientPool {
	mutex := new(sync.Mutex)
	return &ClientPool{
		lock:       mutex,
		cond:       sync.NewCond(mutex),
		clients:    []*http.Client{},
		failureMap: make(map[string]int),
		conf:       conf,
	}
}

func (cp *ClientPool) GetClient(req *http.Request) *http.Client {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	if len(cp.clients) == 0 {
		cp.cond.Wait()
	}

	cli := cp.clients[0]
	cp.clients = cp.clients[1:]

	return cli
}

func (cp *ClientPool) AddClients(clients ...*http.Client) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	for _, cli := range clients {
		addr := getAddrFromClient(cli)
		if !slices.Contains(cp.clientInPool, addr) {
			cp.clientInPool = append(cp.clientInPool, addr)
		}
		cp.clients = append(cp.clients, cli)
		cp.cond.Signal()
	}
}

func (cp *ClientPool) addAvailableProxyClient(proxyAddr string, clientAvailable func(*http.Client) bool) {
	if slices.Contains(cp.clientInPool, proxyAddr) {
		return
	}

	proxyUrl, urlErr := url.Parse(proxyAddr)
	if urlErr != nil {
		return
	}

	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: cp.conf.ClientTimeout,
	}

	if clientAvailable(cli) {
		cp.AddClients(cli)
	}
}

func (cp *ClientPool) proxyList(protocol string, sourceURL string) []string {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return []string{}
	}

	bodyStr := string(bodyBytes)
	lines := strings.Split(bodyStr, "\n")
	proxyAddrs := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			proxyAddrs = append(proxyAddrs, fmt.Sprintf("%s://%s", protocol, line))
		}
	}

	return proxyAddrs
}

func (cp *ClientPool) refreshClients(ctx context.Context, clientAvailable func(*http.Client) bool) {
	proxyAddrs := make([]string, 0)
	if cp.conf.Socks5ProxySourceURL != "" {
		proxyAddrs = append(proxyAddrs, cp.proxyList("socks5", cp.conf.Socks5ProxySourceURL)...)
	}
	if cp.conf.Socks4ProxySourceURL != "" {
		proxyAddrs = append(proxyAddrs, cp.proxyList("socks4", cp.conf.Socks4ProxySourceURL)...)
	}
	if cp.conf.HTTPProxySourceURL != "" {
		proxyAddrs = append(proxyAddrs, cp.proxyList("http", cp.conf.HTTPProxySourceURL)...)
	}

	wg := new(sync.WaitGroup)
	for _, addr := range proxyAddrs {
		wg.Go(func() {
			cp.addAvailableProxyClient(addr, clientAvailable)
		})
		time.Sleep(100 * time.Millisecond)
	}

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	// early return if ctx.Done was called
	select {
	case <-ctx.Done():
		return
	case <-waitCh:
	}
}

func (cp *ClientPool) BackgroundRefreshClients(ctx context.Context, clientAvailable func(*http.Client) bool) {
	// refresh client at the beginning
	cp.refreshClients(ctx, clientAvailable)

	clock := time.NewTicker(cp.conf.RefreshInterval)
	defer clock.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-clock.C:
			cp.refreshClients(ctx, clientAvailable)
		}
	}
}

func (cp *ClientPool) RequestRecorder(pool pool.ClientPool, cli *http.Client, req *http.Request, resp *http.Response, err error) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	isRespSuccess := err == nil && ((resp.StatusCode >= 200 && resp.StatusCode <= 299) || resp.StatusCode == 400 || resp.StatusCode == 404)
	failureCount, ok := cp.failureMap[getAddrFromClient(cli)]
	if !ok {
		failureCount = 0
	}

	cooldownInterval := cp.conf.SuccessCooldownInterval
	if isRespSuccess {
		cp.failureMap[getAddrFromClient(cli)] = 0
	} else if failureCount >= cp.conf.DropClientFailureThreshold {
		// Drop client from pool
		index := slices.Index(cp.clientInPool, getAddrFromClient(cli))
		cp.clientInPool = slices.Delete(cp.clientInPool, index, index+1)
		delete(cp.failureMap, getAddrFromClient(cli))

		return
	} else {
		cp.failureMap[getAddrFromClient(cli)]++
		cooldownInterval = cp.conf.FailureCooldownInterval
	}

	go func() {
		time.Sleep(cooldownInterval)
		pool.AddClients(cli)
	}()
}

func RaiseErrorForNon2xxMiddleware(requester goclient.Requester) goclient.Requester {
	return func(req *http.Request) (*http.Response, error) {
		resp, err := requester(req)
		if err != nil {
			return resp, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, &ErrStatusCode{StatusCode: resp.StatusCode}
		}

		return resp, nil
	}
}

func getAddrFromClient(cli *http.Client) string {
	transport, ok := cli.Transport.(*http.Transport)
	if !ok {
		// Not the expected transport type
		return ""
	}

	// The Proxy field is a function; if using http.ProxyURL(addrUrl), you won't get the URL directly
	proxyURL, urlErr := transport.Proxy(nil)
	if urlErr != nil {
		return ""
	}

	return proxyURL.String()

}
