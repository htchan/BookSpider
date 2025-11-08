package uukanshu

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/htchan/BookSpider/internal/client/v1"
	"github.com/htchan/goclient"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Skipf("not implemented")
}

func TestClient_GetBookInfo(t *testing.T) {
	vendorProtocol = "http"
	vendorHost = strings.TrimLeft(serv.URL, "http://")
	bookURLTemplate = vendorProtocol + "://" + vendorHost + "/book_info_%v.html"
	t.Cleanup(func() {
		vendorProtocol = "https"
		vendorHost = "www.uukanshu.com"
		bookURLTemplate = vendorProtocol + "://" + vendorHost + "/b/%v/"
	})

	tests := []struct {
		name    string
		cli     uukanshuClient
		bookID  string
		want    *client.BookInfo
		wantErr string
	}{
		{
			name:   "happy flow",
			bookID: "success",
			cli: uukanshuClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want: &client.BookInfo{
				Title:         "book name",
				Author:        "author",
				Type:          "type",
				UpdateChapter: "chapter name",
				UpdateDate:    time.Date(time.Now().Year(), 5, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: "",
		},
		{
			name:   "unhappy flow - book not found",
			bookID: "not_found",
			cli: uukanshuClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "title not found\nwriter not found\ntype not found\ndate not found\nchapter not found\nbook fields not found",
		},
		{
			name:   "unhappy flow - timeout",
			bookID: "timeout",
			cli: uukanshuClient{
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
		vendorHost = "www.uukanshu.com"
		chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/b/%v/"
	})

	tests := []struct {
		name    string
		cli     uukanshuClient
		bookID  string
		want    client.ChapterEntryList
		wantErr string
	}{
		{
			name:   "happy flow",
			bookID: "success",
			cli: uukanshuClient{
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
			cli: uukanshuClient{
				cli:     goclient.NewClient(),
				decoder: client.NewDecoder(client.DecodeMethodUTF8),
			},
			want:    nil,
			wantErr: "empty chapter list",
		},
		{
			name:   "unhappy flow - timeout",
			bookID: "timeout",
			cli: uukanshuClient{
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
		vendorHost = "www.uukanshu.com"
		chapterListURLTemplate = vendorProtocol + "://" + vendorHost + "/b/%v/"
	})

	tests := []struct {
		name    string
		cli     uukanshuClient
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
			cli: uukanshuClient{
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
			cli: uukanshuClient{
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
			cli: uukanshuClient{
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
		vendorHost = "www.uukanshu.com"
	})

	tests := []struct {
		name     string
		cli      uukanshuClient
		setConst func()
		want     bool
	}{
		{
			name: "happy flow",
			cli: uukanshuClient{
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
			cli: uukanshuClient{
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
			cli: uukanshuClient{
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
