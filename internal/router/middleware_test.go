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
	mockservice "github.com/htchan/BookSpider/internal/mock/service/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetSiteMiddleware(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		siteName  string
		wantSite  string
		expectRes string
	}{
		{
			name:      "set request context site if site found",
			siteName:  "test",
			wantSite:  "test",
			expectRes: "site found",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			handler := GetSiteMiddleware(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					site := r.Context().Value(ContextKeySiteName).(string)

					assert.Equal(t, test.wantSite, site)
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

func Test_GetBookMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		setupServ       func(ctrl *gomock.Controller) service.ReadDataService
		idHash          string
		expectBook      *model.Book
		expectBookGroup *model.BookGroup
		wantRes         string
	}{
		{
			name: "set request context book for existing id",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().BookGroup(gomock.Any(), "test", "1", "").Return(
					&model.Book{ID: 1},
					&model.BookGroup{{ID: 1}, {ID: 2}},
					nil,
				)

				return serv
			},
			idHash:          "1",
			expectBook:      &model.Book{ID: 1},
			expectBookGroup: &model.BookGroup{{ID: 1}, {ID: 2}},
			wantRes:         "ok",
		},
		{
			name: "set request context book for existing id-hash",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().BookGroup(gomock.Any(), "test", "1", "2s").Return(
					&model.Book{ID: 1, HashCode: 100},
					&model.BookGroup{{ID: 1, HashCode: 100}, {ID: 2}},
					nil,
				)

				return serv
			},
			idHash:          "1-2s",
			expectBook:      &model.Book{ID: 1, HashCode: 100},
			expectBookGroup: &model.BookGroup{{ID: 1, HashCode: 100}, {ID: 2}},
			wantRes:         "ok",
		},
		{
			name: "return error for not exist id",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().BookGroup(gomock.Any(), "test", "1", "").Return(
					nil,
					nil,
					errors.New("some error"),
				)

				return serv
			},
			idHash:          "1",
			expectBook:      nil,
			expectBookGroup: nil,
			wantRes:         `{"error":"book not found"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := GetBookMiddleware(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					bk := r.Context().Value(ContextKeyBook).(*model.Book)
					if !cmp.Equal(bk, test.expectBook) {
						t.Errorf("site diff: %v", cmp.Diff(bk, test.expectBook))
					}

					bkGroup := r.Context().Value(ContextKeyBookGroup).(*model.BookGroup)
					if !cmp.Equal(bkGroup, test.expectBookGroup) {
						t.Errorf("site diff: %v", cmp.Diff(bkGroup, test.expectBookGroup))
					}
					fmt.Fprintln(w, test.wantRes)
				},
			))

			req, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}

			serv := test.setupServ(ctrl)

			ctx := context.WithValue(req.Context(), ContextKeyReadDataServ, serv)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("idHash", test.idHash)
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = context.WithValue(ctx, ContextKeySiteName, "test")
			req = req.WithContext(ctx)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)

			assert.Equal(t, test.wantRes, strings.Trim(res.Body.String(), "\n"))
		})
	}
}

func Test_GetSearchParamsMiddleware(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name       string
		url        string
		wantTitle  string
		wantWriter string
		wantRes    string
	}{
		{
			name:       "empty title and empty writer",
			url:        "http://host/test",
			wantTitle:  "",
			wantWriter: "",
			wantRes:    "ok",
		},
		{
			name:       "title and not writer",
			url:        "http://host/test?title=title",
			wantTitle:  "title",
			wantWriter: "",
			wantRes:    "ok",
		},
		{
			name:       "not title and writer",
			url:        "http://host/test?writer=writer",
			wantTitle:  "",
			wantWriter: "writer",
			wantRes:    "ok",
		},
		{
			name:       "title and writer",
			url:        "http://host/test?title=title&writer=writer",
			wantTitle:  "title",
			wantWriter: "writer",
			wantRes:    "ok",
		},
		{
			name:       "some unrelated params",
			url:        "http://host/test?title=title&writer=writer&unknown=1",
			wantTitle:  "title",
			wantWriter: "writer",
			wantRes:    "ok",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := GetSearchParamsMiddleware(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					title := r.Context().Value(ContextKeyTitle).(string)
					if title != test.wantTitle {
						t.Errorf("title diff: %v", cmp.Diff(title, test.wantTitle))
					}

					writer := r.Context().Value(ContextKeyWriter).(string)
					if writer != test.wantWriter {
						t.Errorf("writer diff: %v", cmp.Diff(writer, test.wantWriter))
					}

					fmt.Fprintln(w, test.wantRes)
				},
			))
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.wantRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.wantRes)
			}
		})
	}
}

func Test_GetPageParamsMiddleware(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name       string
		url        string
		wantLimit  int
		wantOffset int
		wantRes    string
	}{
		{
			name:       "empty page and empty per page",
			url:        "http://host/test",
			wantLimit:  0,
			wantOffset: 0,
			wantRes:    "ok",
		},
		{
			name:       "page and empty per page",
			url:        "http://host/test?page=2",
			wantLimit:  0,
			wantOffset: 0,
			wantRes:    "ok",
		},
		{
			name:       "empty page and per page",
			url:        "http://host/test?per_page=10",
			wantLimit:  10,
			wantOffset: 0,
			wantRes:    "ok",
		},
		{
			name:       "page and per page",
			url:        "http://host/test?page=2&per_page=10",
			wantLimit:  10,
			wantOffset: 20,
			wantRes:    "ok",
		},
		{
			name:       "page and per page of unknown value",
			url:        "http://host/test?page=limit&per_page=offset",
			wantLimit:  0,
			wantOffset: 0,
			wantRes:    "ok",
		},
		{
			name:       "some unrelated params",
			url:        "http://host/test?page=2&per_page=10&unknown=1",
			wantLimit:  10,
			wantOffset: 20,
			wantRes:    "ok",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := GetPageParamsMiddleware(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					limit := r.Context().Value(ContextKeyLimit).(int)
					if limit != test.wantLimit {
						t.Errorf("limit diff: %v", cmp.Diff(limit, test.wantLimit))
					}

					offset := r.Context().Value(ContextKeyOffset).(int)
					if offset != test.wantOffset {
						t.Errorf("offset diff: %v", cmp.Diff(offset, test.wantOffset))
					}

					fmt.Fprintln(w, test.wantRes)
				},
			))
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)

			if strings.Trim(res.Body.String(), "\n") != test.wantRes {
				t.Error("got different response as expect")
				t.Error(res.Body.String())
				t.Error(test.wantRes)
			}
		})
	}
}
