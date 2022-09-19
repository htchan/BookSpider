package router

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/site"
)

func Test_GeneralInfoAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		sites     map[string]*site.Site
		url       string
		expectRes string
	}{
		{
			name: "works",
			sites: map[string]*site.Site{
				"test1": site.MockSite("test1", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{}, nil),
				"test2": site.MockSite("test2", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{}, nil),
			},
			url:       "https://localhost/data",
			expectRes: `{"test1":{"BookCount":0,"WriterCount":0,"ErrorCount":0,"UniqueBookCount":0,"MaxBookID":0,"LatestSuccessID":0,"DownloadCount":0,"StatusCount":null},"test2":{"BookCount":0,"WriterCount":0,"ErrorCount":0,"UniqueBookCount":0,"MaxBookID":0,"LatestSuccessID":0,"DownloadCount":0,"StatusCount":null}}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}

			res := httptest.NewRecorder()
			GeneralInfoAPIHandler(test.sites).ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_SiteInfoAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		st        *site.Site
		url       string
		expectRes string
	}{
		{
			name:      "works",
			st:        site.MockSite("test1", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{}, nil),
			url:       "https://localhost/data",
			expectRes: `{"BookCount":0,"WriterCount":0,"ErrorCount":0,"UniqueBookCount":0,"MaxBookID":0,"LatestSuccessID":0,"DownloadCount":0,"StatusCount":null}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), "site", test.st)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			SiteInfoAPIHandler(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_BookSearchAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		st            *site.Site
		url           string
		title, writer string
		limit, offset int
		expectRes     string
	}{
		{
			name:      "works",
			st:        site.MockSite("test", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{}, nil),
			url:       "https://localhost/data",
			title:     "title",
			writer:    "writer",
			limit:     10,
			offset:    0,
			expectRes: `{"books":[]}`,
		},
		{
			name:      "error",
			st:        site.MockSite("test", mock.MockRepostory{Err: errors.New("error")}, config.BookConfig{}, config.SiteConfig{}, nil),
			url:       "https://localhost/data",
			title:     "title",
			writer:    "writer",
			limit:     10,
			offset:    0,
			expectRes: `{"error":"error"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), "site", test.st)
			ctx = context.WithValue(ctx, "title", test.title)
			ctx = context.WithValue(ctx, "writer", test.writer)
			ctx = context.WithValue(ctx, "limit", test.limit)
			ctx = context.WithValue(ctx, "offset", test.offset)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookSearchAPIHandler(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_BookRandomAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		st            *site.Site
		url           string
		limit, offset int
		expectRes     string
	}{
		{
			name:      "works",
			st:        site.MockSite("test", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{}, nil),
			url:       "https://localhost/data",
			limit:     10,
			offset:    0,
			expectRes: `{"books":[]}`,
		},
		{
			name:      "error",
			st:        site.MockSite("test", mock.MockRepostory{Err: errors.New("error")}, config.BookConfig{}, config.SiteConfig{}, nil),
			url:       "https://localhost/data",
			limit:     10,
			offset:    0,
			expectRes: `{"error":"error"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), "site", test.st)
			ctx = context.WithValue(ctx, "limit", test.limit)
			ctx = context.WithValue(ctx, "offset", test.offset)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookRandomAPIHandler(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_BookInfoAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		bk        *model.Book
		expectRes string
	}{
		{
			name: "works",
			url:  "https://localhost/data",
			bk: &model.Book{
				Site: "test", ID: 1, HashCode: 100,
				Title: "title", Writer: model.Writer{ID: 2, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
				IsDownloaded: true, Error: errors.New("error"),
			},
			expectRes: `{"site":"test","id":1,"hash_code":"2s","title":"title","writer":"writer","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"error":"error"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), "book", test.bk)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookInfoAPIHandler(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_BookDownloadAPIHandler(t *testing.T) {
	t.Parallel()

	os.Mkdir("storage", os.ModePerm)
	os.WriteFile("storage/1.txt", []byte("data"), os.ModePerm)

	t.Cleanup(func() {
		os.RemoveAll("storage")
	})

	tests := []struct {
		name      string
		url       string
		st        *site.Site
		bk        *model.Book
		expectRes string
	}{
		{
			name:      "works",
			url:       "https://localhost/data",
			st:        site.MockSite("test", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{Storage: "storage"}, nil),
			bk:        &model.Book{Site: "test", ID: 1, HashCode: 0, Status: model.End, IsDownloaded: true},
			expectRes: `data`,
		},
		{
			name:      "bk is not download",
			url:       "https://localhost/data",
			st:        site.MockSite("test", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{Storage: "storage"}, nil),
			bk:        &model.Book{Site: "test", ID: 1, HashCode: 0},
			expectRes: `{"error":"book is not download"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), "site", test.st)
			ctx = context.WithValue(ctx, "book", test.bk)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookDownloadAPIHandler(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_DBStatAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		sites     map[string]*site.Site
		expectRes string
	}{
		{
			name: "works",
			url:  "https://localhost/data",
			sites: map[string]*site.Site{
				"test1": site.MockSite("test1", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{}, nil),
				"test2": site.MockSite("test2", mock.MockRepostory{}, config.BookConfig{}, config.SiteConfig{}, nil),
			},
			expectRes: `{"stats":[{"MaxOpenConnections":0,"OpenConnections":0,"InUse":0,"Idle":0,"WaitCount":0,"WaitDuration":0,"MaxIdleClosed":0,"MaxIdleTimeClosed":0,"MaxLifetimeClosed":0},{"MaxOpenConnections":0,"OpenConnections":0,"InUse":0,"Idle":0,"WaitCount":0,"WaitDuration":0,"MaxIdleClosed":0,"MaxIdleTimeClosed":0,"MaxLifetimeClosed":0}]}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}

			res := httptest.NewRecorder()
			DBStatsAPIHandler(test.sites).ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}
