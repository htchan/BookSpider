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

func TestNewClient(t *testing.T) {
	tests := []struct {
		name             string
		conf             config.ClientConfig
		wantDecodeMethod client.DecodeMethod
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
			wantDecodeMethod: client.DecodeMethodUTF8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cli := NewClient(ctx, tt.conf)
			assert.NotNil(t, cli)
			assert.Implements(t, (*client.Client)(nil), cli)

			// Verify it's the correct type
			hjwzwCli, ok := cli.(*hjwzwClient)
			assert.True(t, ok)
			assert.NotNil(t, hjwzwCli.cli)
			assert.NotNil(t, hjwzwCli.decoder)

			// Test decoder can decode
			testStr := "test string"
			decoded, err := hjwzwCli.decoder.Decode(testStr)
			assert.NoError(t, err)
			assert.NotEmpty(t, decoded)
		})
	}
}

func TestClient_GetBookInfo(t *testing.T) {
	vendorProtocol = "http"
	vendorHost = strings.TrimLeft(serv.URL, "http://")
	bookURLTemplate = vendorProtocol + "://" + vendorHost + "/book_info_%v.html"
	t.Cleanup(func() {
		vendorProtocol = "https"
		vendorHost = "tw.hjwzw.com"
		bookURLTemplate = vendorProtocol + "://" + vendorHost + "/Book/%v/"
	})

	tests := []struct {
		name    string
		cli     hjwzwClient
		bookID  string
		want    *client.BookInfo
		wantErr string
	}{
		{
			name:   "happy flow",
			bookID: "success",
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want: &client.BookInfo{
				Title:         "book name",
				Author:        "author",
				Type:          "type",
				UpdateChapter: "chapter name",
				UpdateDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: "",
		},
		{
			name:   "unhappy flow - book not found",
			bookID: "not_found",
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "title not found\nwriter not found\ntype not found\ndate not found\nchapter not found\nbook fields not found",
		},
		{
			name:   "unhappy flow - timeout",
			bookID: "timeout",
			cli: hjwzwClient{
				cli: goclient.NewClient(
					goclient.WithRequester(func(r *http.Request) (*http.Response, error) {
						cli := &http.Client{Timeout: time.Millisecond}
						return cli.Do(r)
					}),
				),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetBookInfo(context.Background(), tt.bookID)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_GetBookChapterList(t *testing.T) {
	vendorProtocol = "http"
	vendorHost = strings.TrimLeft(serv.URL, "http://")
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/chapter_list_%v.html"
	t.Cleanup(func() {
		vendorProtocol = "https"
		vendorHost = "tw.hjwzw.com"
		chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/Book/Chapter/%v/"
	})

	tests := []struct {
		name    string
		cli     hjwzwClient
		bookID  string
		want    client.ChapterEntryList
		wantErr string
	}{
		{
			name:   "happy flow",
			bookID: "success",
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want: client.ChapterEntryList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "chapter url 2", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: "chapter name 3"},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantErr: "",
		},
		{
			name:   "unhappy flow - book not found",
			bookID: "not_found",
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "empty chapter list",
		},
		{
			name:   "unhappy flow - timeout",
			bookID: "timeout",
			cli: hjwzwClient{
				cli: goclient.NewClient(
					goclient.WithRequester(func(r *http.Request) (*http.Response, error) {
						cli := &http.Client{Timeout: time.Millisecond}
						return cli.Do(r)
					}),
				),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetBookChapterList(context.Background(), tt.bookID)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_GetChapterContent(t *testing.T) {
	vendorProtocol = "http"
	vendorHost = strings.TrimLeft(serv.URL, "http://")
	chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/chapter_list_%v.html"
	t.Cleanup(func() {
		vendorProtocol = "https"
		vendorHost = "tw.hjwzw.com"
		chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/Book/Chapter/%v/"
	})

	tests := []struct {
		name    string
		cli     hjwzwClient
		chapter client.ChapterEntry
		want    *client.ChapterContent
		wantErr string
	}{
		{
			name: "happy flow",
			chapter: client.ChapterEntry{
				URL:   serv.URL + "/chapter/success",
				Title: "chapter name 1",
			},
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want: &client.ChapterContent{
				Title: "chapter name",
				Body:  "chapter content",
			},
			wantErr: "",
		},
		{
			name: "unhappy flow - book not found",
			chapter: client.ChapterEntry{
				URL:   serv.URL + "/chapter/not_found",
				Title: "chapter name 1",
			},
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "chapter title not found\nchapter content not found\nbook fields not found",
		},
		{
			name: "unhappy flow - timeout",
			chapter: client.ChapterEntry{
				URL:   serv.URL + "/chapter/timeout",
				Title: "chapter name 1",
			},
			cli: hjwzwClient{
				cli: goclient.NewClient(
					goclient.WithRequester(func(r *http.Request) (*http.Response, error) {
						cli := &http.Client{Timeout: time.Millisecond}
						return cli.Do(r)
					}),
				),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetChapterContent(context.Background(), tt.chapter)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_Available(t *testing.T) {
	vendorProtocol = "http"
	t.Cleanup(func() {
		vendorProtocol = "https"
		vendorHost = "tw.hjwzw.com"
	})

	tests := []struct {
		name     string
		cli      hjwzwClient
		setConst func()
		want     bool
	}{
		{
			name: "happy flow",
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			setConst: func() {
				vendorHost = strings.TrimLeft(serv.URL, "http://") + "/available/success"
			},
			want: true,
		},
		{
			name: "unhappy flow - book not found",
			cli: hjwzwClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			setConst: func() {
				vendorHost = strings.TrimLeft(serv.URL, "http://") + "/available/not_found"
			},
			want: false,
		},
		{
			name: "unhappy flow - timeout",
			cli: hjwzwClient{
				cli: goclient.NewClient(
					goclient.WithRequester(func(r *http.Request) (*http.Response, error) {
						cli := &http.Client{Timeout: time.Millisecond}
						return cli.Do(r)
					}),
				),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			setConst: func() {
				vendorHost = strings.TrimLeft(serv.URL, "http://") + "/available/timeout"
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setConst()
			got := tt.cli.Available(context.Background())
			assert.Equal(t, tt.want, got)
		})
	}
}
