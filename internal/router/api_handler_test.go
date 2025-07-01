package router

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	mockservice "github.com/htchan/BookSpider/internal/mock/service/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GeneralInfoAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupServs func(ctrl *gomock.Controller) map[string]service.Service
		url        string
		expectRes  string
	}{
		{
			name: "works",
			setupServs: func(ctrl *gomock.Controller) map[string]service.Service {
				serv1 := mockservice.NewMockService(ctrl)
				serv1.EXPECT().Name().Return("test1")
				serv1.EXPECT().Stats(gomock.Any()).Return(repo.Summary{})

				serv2 := mockservice.NewMockService(ctrl)
				serv2.EXPECT().Name().Return("test2")
				serv2.EXPECT().Stats(gomock.Any()).Return(repo.Summary{})

				return map[string]service.Service{
					"test1": serv1,
					"test2": serv2,
				}
			},
			url:       "https://localhost/data",
			expectRes: `{"test1":{"BookCount":0,"WriterCount":0,"ErrorCount":0,"UniqueBookCount":0,"MaxBookID":0,"LatestSuccessID":0,"DownloadCount":0,"StatusCount":null},"test2":{"BookCount":0,"WriterCount":0,"ErrorCount":0,"UniqueBookCount":0,"MaxBookID":0,"LatestSuccessID":0,"DownloadCount":0,"StatusCount":null}}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}

			res := httptest.NewRecorder()
			GeneralInfoAPIHandler(test.setupServs(ctrl)).ServeHTTP(res, req)

			assert.Equal(t, test.expectRes, strings.Trim(res.Body.String(), "\n"))
		})
	}
}

func Test_SiteInfoAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupServ func(ctrl *gomock.Controller) service.ReadDataService
		url       string
		expectRes string
	}{
		{
			name: "works",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().Stats(gomock.Any(), "test").Return(repo.Summary{})

				return serv
			},
			url:       "https://localhost/data",
			expectRes: `{"BookCount":0,"WriterCount":0,"ErrorCount":0,"UniqueBookCount":0,"MaxBookID":0,"LatestSuccessID":0,"DownloadCount":0,"StatusCount":null}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), ContextKeyReadDataServ, test.setupServ(ctrl))
			ctx = context.WithValue(ctx, ContextKeySiteName, "test")
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			SiteInfoAPIHandler(res, req)

			assert.Equal(t, test.expectRes, strings.Trim(res.Body.String(), "\n"))
		})
	}
}

func Test_BookSearchAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupServ     func(ctrl *gomock.Controller) service.ReadDataService
		url           string
		title, writer string
		limit, offset int
		expectRes     string
	}{
		{
			name: "works",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().SearchBooks(gomock.Any(), "title 1", "writer 1", 10, 0).Return([]model.Book{}, nil)

				return serv
			},
			url:       "https://localhost/data",
			title:     "title 1",
			writer:    "writer 1",
			limit:     10,
			offset:    0,
			expectRes: `{"books":[]}`,
		},
		{
			name: "error",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().SearchBooks(gomock.Any(), "title 1", "writer 1", 10, 0).Return(nil, errors.New("some error"))

				return serv
			},
			url:       "https://localhost/data",
			title:     "title 1",
			writer:    "writer 1",
			limit:     10,
			offset:    0,
			expectRes: `{"error":"some error"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), ContextKeyReadDataServ, test.setupServ(ctrl))
			ctx = context.WithValue(ctx, ContextKeyTitle, test.title)
			ctx = context.WithValue(ctx, ContextKeyWriter, test.writer)
			ctx = context.WithValue(ctx, ContextKeyLimit, test.limit)
			ctx = context.WithValue(ctx, ContextKeyOffset, test.offset)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookSearchAPIHandler(res, req)

			assert.Equal(t, test.expectRes, strings.Trim(res.Body.String(), "\n"))
		})
	}
}

func Test_BookRandomAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupServ     func(ctrl *gomock.Controller) service.ReadDataService
		url           string
		limit, offset int
		expectRes     string
	}{
		{
			name: "works",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().RandomBooks(gomock.Any(), 10).Return([]model.Book{}, nil)

				return serv
			},
			url:       "https://localhost/data",
			limit:     10,
			offset:    0,
			expectRes: `{"books":[]}`,
		},
		{
			name: "error",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().RandomBooks(gomock.Any(), 10).Return(nil, errors.New("some error"))

				return serv
			},
			url:       "https://localhost/data",
			limit:     10,
			offset:    0,
			expectRes: `{"error":"some error"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), ContextKeyReadDataServ, test.setupServ(ctrl))
			ctx = context.WithValue(ctx, ContextKeyLimit, test.limit)
			ctx = context.WithValue(ctx, ContextKeyOffset, test.offset)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookRandomAPIHandler(res, req)

			assert.Equal(t, test.expectRes, strings.Trim(res.Body.String(), "\n"))
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
				Title: "title 1", Writer: model.Writer{ID: 2, Name: "writer 1"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.StatusInProgress,
				IsDownloaded: true, Error: errors.New("error"),
			},
			expectRes: `{"site":"test","id":1,"hash_code":"2s","title":"title 1","writer":"writer 1","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"error":"error"}`,
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
			ctx := context.WithValue(req.Context(), ContextKeyBook, test.bk)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookInfoAPIHandler(res, req)

			assert.Equal(t, test.expectRes, strings.Trim(res.Body.String(), "\n"))
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
		setupServ func(ctrl *gomock.Controller) service.ReadDataService
		bk        *model.Book
		expectRes string
	}{
		{
			name: "works",
			url:  "https://localhost/data",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().
					BookContent(gomock.Any(), &model.Book{Site: "test", ID: 1, HashCode: 0, Status: model.StatusEnd, IsDownloaded: true}).
					Return("data", nil)

				return serv
			},
			bk:        &model.Book{Site: "test", ID: 1, HashCode: 0, Status: model.StatusEnd, IsDownloaded: true},
			expectRes: `data`,
		},
		{
			name: "bk is not download",
			url:  "https://localhost/data",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().
					BookContent(gomock.Any(), &model.Book{Site: "test", ID: 1, HashCode: 0}).
					Return("", errors.New("some error"))

				return serv
			},
			bk:        &model.Book{Site: "test", ID: 1, HashCode: 0},
			expectRes: `{"error":"some error"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}
			ctx := context.WithValue(req.Context(), ContextKeyReadDataServ, test.setupServ(ctrl))
			ctx = context.WithValue(ctx, ContextKeyBook, test.bk)
			req = req.WithContext(ctx)

			res := httptest.NewRecorder()
			BookDownloadAPIHandler(res, req)

			assert.Equal(t, test.expectRes, strings.Trim(res.Body.String(), "\n"))
		})
	}
}

func Test_DBStatAPIHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		setupServ func(ctrl *gomock.Controller) service.ReadDataService
		expectRes string
	}{
		{
			name: "works",
			url:  "https://localhost/data",
			setupServ: func(ctrl *gomock.Controller) service.ReadDataService {
				serv := mockservice.NewMockReadDataService(ctrl)
				serv.EXPECT().DBStats(gomock.Any()).Return(sql.DBStats{})

				return serv
			},
			expectRes: `{"stats":[{"MaxOpenConnections":0,"OpenConnections":0,"InUse":0,"Idle":0,"WaitCount":0,"WaitDuration":0,"MaxIdleClosed":0,"MaxIdleTimeClosed":0,"MaxLifetimeClosed":0}]}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Errorf("cannot init request: %v", err)
				return
			}

			res := httptest.NewRecorder()
			DBStatsAPIHandler(test.setupServ(ctrl)).ServeHTTP(res, req)

			assert.Equal(t, test.expectRes, strings.Trim(res.Body.String(), "\n"))
		})
	}
}
