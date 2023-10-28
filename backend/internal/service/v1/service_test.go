package service

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	circuitbreaker "github.com/htchan/BookSpider/internal/client/v2/circuit_breaker"
	"github.com/htchan/BookSpider/internal/client/v2/retry"
	"github.com/htchan/BookSpider/internal/client/v2/simple"
	config "github.com/htchan/BookSpider/internal/config_new"
	mockclient "github.com/htchan/BookSpider/internal/mock/client/v2"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	mockvendor "github.com/htchan/BookSpider/internal/mock/vendorservice"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	serv "github.com/htchan/BookSpider/internal/service"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		siteName      string
		rpo           repo.Repository
		vendorService vendor.VendorService
		sema          *semaphore.Weighted
		conf          config.SiteConfig
		want          *ServiceImpl
	}{
		{
			name: "happy flow",
			want: &ServiceImpl{
				cli: retry.NewClient(
					&retry.RetryClientConfig{},
					circuitbreaker.NewClient(
						&circuitbreaker.CircuitBreakerClientConfig{},
						simple.NewClient(&simple.SimpleClientConfig{}),
					),
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewService(test.siteName, test.rpo, test.vendorService, test.sema, test.conf)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImpl_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		serv *ServiceImpl
		want string
	}{
		{
			name: "return site name",
			serv: &ServiceImpl{name: "test"},
			want: "test",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.serv.Name()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImpl_bookFileLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		serv *ServiceImpl
		bk   *model.Book
		want string
	}{
		{
			name: "return book location for book with id only",
			serv: &ServiceImpl{conf: config.SiteConfig{Storage: "/data"}},
			bk:   &model.Book{ID: 123, HashCode: 0},
			want: "/data/123.txt",
		},
		{
			name: "return book location for book with id and hash code",
			serv: &ServiceImpl{conf: config.SiteConfig{Storage: "/data"}},
			bk:   &model.Book{ID: 123, HashCode: 456},
			want: "/data/123-vco.txt",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.serv.bookFileLocation(test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImpl_checkBookStorage(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll("./check-book-storage"))
	})

	if !assert.NoError(t, os.Mkdir("./check-book-storage", os.ModePerm)) ||
		!assert.NoError(t, os.WriteFile("./check-book-storage/123.txt", []byte("test"), 0644)) {
		return
	}

	tests := []struct {
		name   string
		serv   *ServiceImpl
		bk     *model.Book
		want   bool
		wantBk *model.Book
	}{
		{
			name:   "book status is downloaded and file exists",
			serv:   &ServiceImpl{conf: config.SiteConfig{Storage: "./check-book-storage/"}},
			bk:     &model.Book{ID: 123, HashCode: 0, IsDownloaded: true},
			want:   false,
			wantBk: &model.Book{ID: 123, HashCode: 0, IsDownloaded: true},
		},
		{
			name:   "book status is downloaded and file does not exist",
			serv:   &ServiceImpl{conf: config.SiteConfig{Storage: "./check-book-storage/"}},
			bk:     &model.Book{ID: 456, HashCode: 0, IsDownloaded: true},
			want:   true,
			wantBk: &model.Book{ID: 456, HashCode: 0, IsDownloaded: false},
		},
		{
			name:   "book status is not downloaded and file does not exist",
			serv:   &ServiceImpl{conf: config.SiteConfig{Storage: "./check-book-storage/"}},
			bk:     &model.Book{ID: 456, HashCode: 0, IsDownloaded: false},
			want:   false,
			wantBk: &model.Book{ID: 456, HashCode: 0, IsDownloaded: false},
		},
		{
			name:   "book status is not downloaded and file exists",
			serv:   &ServiceImpl{conf: config.SiteConfig{Storage: "./check-book-storage/"}},
			bk:     &model.Book{ID: 123, HashCode: 0, IsDownloaded: false},
			want:   true,
			wantBk: &model.Book{ID: 123, HashCode: 0, IsDownloaded: true},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.serv.checkBookStorage(test.bk)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantBk, test.bk)
		})

	}
}

func TestServiceImpl_PatchDownloadStatus(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll("./patch-download-status"))
	})

	if !assert.NoError(t, os.Mkdir("./patch-download-status", os.ModePerm)) ||
		!assert.NoError(t, os.WriteFile("./patch-download-status/123.txt", []byte("test"), 0644)) {
		return
	}

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		wantError  error
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				bkCh := make(chan model.Book, 10)

				go func() {
					bkCh <- model.Book{ID: 123, HashCode: 0, IsDownloaded: true}
					bkCh <- model.Book{ID: 456, HashCode: 0, IsDownloaded: true}
					close(bkCh)
				}()

				rpo.EXPECT().FindAllBooks().Return(bkCh, nil)
				rpo.EXPECT().UpdateBook(&model.Book{ID: 456, HashCode: 0, IsDownloaded: false}).Return(nil)

				return &ServiceImpl{
					conf: config.SiteConfig{Storage: "./patch-download-status"},
					rpo:  rpo,
					sema: semaphore.NewWeighted(1),
				}
			},
			wantError: nil,
		},
		{
			name: "FindAllBooks returns error",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				rpo.EXPECT().FindAllBooks().Return(nil, service.ErrUnavailable)

				return &ServiceImpl{
					conf: config.SiteConfig{Storage: "./patch-download-status"},
					rpo:  rpo,
					sema: semaphore.NewWeighted(1),
				}
			},
			wantError: service.ErrUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.getService(ctrl)
			err := serv.PatchDownloadStatus(context.Background())
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_PatchMissingRecords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		wantError  error
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				vendorService := mockvendor.NewMockVendorService(ctrl)
				cli := mockclient.NewMockBookClient(ctrl)

				hashcode := model.GenerateHash()

				rpo.EXPECT().FindAllBookIDs().Return([]int{1, 2, 4}, nil)
				vendorService.EXPECT().FindMissingIds([]int{1, 2, 4}).Return([]int{3})
				rpo.EXPECT().CreateBook(&model.Book{Site: "serv", ID: 3, HashCode: hashcode}).Return(nil)
				vendorService.EXPECT().BookURL("3").Return("http://testing.com/1234")
				cli.EXPECT().Get(gomock.Any(), "http://testing.com/1234").Return("result", nil)
				vendorService.EXPECT().ParseBook("result").Return(&vendor.BookInfo{
					Title: "title", Writer: "writer", Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
				}, nil)
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)
				rpo.EXPECT().UpdateBook(&model.Book{
					Site: "serv", ID: 3, HashCode: hashcode,
					Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter", Status: model.StatusInProgress,
				}).Return(nil)
				rpo.EXPECT().SaveError(&model.Book{
					Site: "serv", ID: 3, HashCode: hashcode,
					Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter", Status: model.StatusInProgress,
				}, nil).Return(nil)

				return &ServiceImpl{
					name:          "serv",
					rpo:           rpo,
					cli:           cli,
					vendorService: vendorService,
					sema:          semaphore.NewWeighted(1),
				}
			},
			wantError: nil,
		},
		{
			name: "FindAllBooks returns error",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				rpo.EXPECT().FindAllBookIDs().Return(nil, service.ErrUnavailable)

				return &ServiceImpl{
					rpo:  rpo,
					sema: semaphore.NewWeighted(1),
				}
			},
			wantError: service.ErrUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.getService(ctrl)
			err := serv.PatchMissingRecords(context.Background())
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_CheckAvailability(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		wantError  error
	}{
		{
			name: "site available",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				vendorService := mockvendor.NewMockVendorService(ctrl)
				cli := mockclient.NewMockBookClient(ctrl)

				vendorService.EXPECT().AvailabilityURL().Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("result", nil)
				vendorService.EXPECT().IsAvailable("result").Return(true)

				return &ServiceImpl{
					name:          "serv",
					cli:           cli,
					vendorService: vendorService,
					sema:          semaphore.NewWeighted(1),
				}
			},
			wantError: nil,
		},
		{
			name: "site unavailable",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				vendorService := mockvendor.NewMockVendorService(ctrl)
				cli := mockclient.NewMockBookClient(ctrl)

				vendorService.EXPECT().AvailabilityURL().Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("result", nil)
				vendorService.EXPECT().IsAvailable("result").Return(false)

				return &ServiceImpl{
					name:          "serv",
					cli:           cli,
					vendorService: vendorService,
					sema:          semaphore.NewWeighted(1),
				}
			},
			wantError: serv.ErrUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.getService(ctrl)
			err := serv.CheckAvailability(context.Background())
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
