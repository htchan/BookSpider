package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/site"
)

func Test_GetSite(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		sites      map[string]*site.Site
		siteName   string
		expectSite *site.Site
		expectRes  string
	}{
		{
			name: "set request context site if site found",
			sites: map[string]*site.Site{
				"test": {Name: "hello"},
			},
			siteName:   "test",
			expectSite: &site.Site{Name: "hello"},
			expectRes:  "site found",
		},
		{
			name:       "return error if site not found",
			sites:      map[string]*site.Site{},
			siteName:   "unknown",
			expectSite: nil,
			expectRes:  `{"error": "site not found"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handlerFunc := GetSite(test.sites)
			handler := handlerFunc(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					st := r.Context().Value("site").(*site.Site)
					if !cmp.Equal(st, test.expectSite) {
						t.Errorf("site diff: %v", cmp.Diff(st, test.expectSite))
					}
					fmt.Fprintln(w, test.expectRes)
				},
			))
			req, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := req.Context()
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("siteName", test.siteName)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_GetBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		st         *site.Site
		idHash     string
		expectBook *model.Book
		expectRes  string
	}{
		{
			name:       "set request context book for existing id",
			st:         site.MockSite("test", mock.MockRepostory{}, &config.BookConfig{}, &config.SiteConfig{}, nil),
			idHash:     "1",
			expectBook: &model.Book{ID: 1},
			expectRes:  "ok",
		},
		{
			name:       "set request context book for existing id-hash",
			st:         site.MockSite("test", mock.MockRepostory{}, &config.BookConfig{}, &config.SiteConfig{}, nil),
			idHash:     "1-2s",
			expectBook: &model.Book{ID: 1, HashCode: 100},
			expectRes:  "ok",
		},
		{
			name:       "return error for not exist id",
			st:         site.MockSite("test", mock.MockRepostory{Err: errors.New("")}, &config.BookConfig{}, &config.SiteConfig{}, nil),
			idHash:     "1",
			expectBook: &model.Book{ID: 1},
			expectRes:  `{"error": "book not found"}`,
		},
		{
			name:       "return error for not exist id",
			st:         site.MockSite("test", mock.MockRepostory{Err: errors.New("")}, &config.BookConfig{}, &config.SiteConfig{}, nil),
			idHash:     "1-2s",
			expectBook: &model.Book{ID: 1, HashCode: 100},
			expectRes:  `{"error": "book not found"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := GetBook(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					bk := r.Context().Value("book").(*model.Book)
					if !cmp.Equal(bk, test.expectBook) {
						t.Errorf("site diff: %v", cmp.Diff(bk, test.expectBook))
					}
					fmt.Fprintln(w, test.expectRes)
				},
			))
			req, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), "site", test.st)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("idHash", test.idHash)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_GetSearchParams(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name         string
		url          string
		expectTitle  string
		expectWriter string
		expectRes    string
	}{
		{
			name:         "empty title and empty writer",
			url:          "http://host/test",
			expectTitle:  "",
			expectWriter: "",
			expectRes:    "ok",
		},
		{
			name:         "title and not writer",
			url:          "http://host/test?title=title",
			expectTitle:  "title",
			expectWriter: "",
			expectRes:    "ok",
		},
		{
			name:         "not title and writer",
			url:          "http://host/test?writer=writer",
			expectTitle:  "",
			expectWriter: "writer",
			expectRes:    "ok",
		},
		{
			name:         "title and writer",
			url:          "http://host/test?title=title&writer=writer",
			expectTitle:  "title",
			expectWriter: "writer",
			expectRes:    "ok",
		},
		{
			name:         "some unrelated params",
			url:          "http://host/test?title=title&writer=writer&unknown=1",
			expectTitle:  "title",
			expectWriter: "writer",
			expectRes:    "ok",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := GetSearchParams(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					title := r.Context().Value("title").(string)
					if title != test.expectTitle {
						t.Errorf("title diff: %v", cmp.Diff(title, test.expectTitle))
					}

					writer := r.Context().Value("writer").(string)
					if writer != test.expectWriter {
						t.Errorf("writer diff: %v", cmp.Diff(writer, test.expectWriter))
					}

					fmt.Fprintln(w, test.expectRes)
				},
			))
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}

func Test_GetPageParams(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name         string
		url          string
		expectLimit  int
		expectOffset int
		expectRes    string
	}{
		{
			name:         "empty page and empty per page",
			url:          "http://host/test",
			expectLimit:  0,
			expectOffset: 0,
			expectRes:    "ok",
		},
		{
			name:         "page and empty per page",
			url:          "http://host/test?page=2",
			expectLimit:  0,
			expectOffset: 0,
			expectRes:    "ok",
		},
		{
			name:         "empty page and per page",
			url:          "http://host/test?per_page=10",
			expectLimit:  10,
			expectOffset: 0,
			expectRes:    "ok",
		},
		{
			name:         "page and per page",
			url:          "http://host/test?page=2&per_page=10",
			expectLimit:  10,
			expectOffset: 20,
			expectRes:    "ok",
		},
		{
			name:         "page and per page of unknown value",
			url:          "http://host/test?page=limit&per_page=offset",
			expectLimit:  0,
			expectOffset: 0,
			expectRes:    "ok",
		},
		{
			name:         "some unrelated params",
			url:          "http://host/test?page=2&per_page=10&unknown=1",
			expectLimit:  10,
			expectOffset: 20,
			expectRes:    "ok",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := GetPageParams(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					limit := r.Context().Value("limit").(int)
					if limit != test.expectLimit {
						t.Errorf("limit diff: %v", cmp.Diff(limit, test.expectLimit))
					}

					offset := r.Context().Value("offset").(int)
					if offset != test.expectOffset {
						t.Errorf("offset diff: %v", cmp.Diff(offset, test.expectOffset))
					}

					fmt.Fprintln(w, test.expectRes)
				},
			))
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.expectRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.expectRes)
			}
		})
	}
}
