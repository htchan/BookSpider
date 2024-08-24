package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/htchan/BookSpider/internal/config/v2"
	clientmock "github.com/htchan/BookSpider/internal/mock/client/v2"
	repomock "github.com/htchan/BookSpider/internal/mock/repo"
	vendormock "github.com/htchan/BookSpider/internal/mock/vendorservice"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	serv "github.com/htchan/BookSpider/internal/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestServiceImpl_ExploreBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		getServ   func(ctrl *gomock.Controller) *ServiceImpl
		bk        *model.Book
		wantBk    *model.Book
		wantError error
	}{
		{
			name: "book status is not error",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				return &ServiceImpl{}
			},
			bk:        &model.Book{Status: model.StatusEnd},
			wantBk:    &model.Book{Status: model.StatusEnd},
			wantError: serv.ErrBookStatusNotError,
		},
		{
			name: "input is a completely new book (error is nil)",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)

				rpo.EXPECT().CreateBook(&model.Book{ID: 1, Status: model.StatusError})
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("", serv.ErrUnavailable)
				expectErr := fmt.Errorf("get book page failed: %w", serv.ErrUnavailable)
				rpo.EXPECT().SaveError(&model.Book{ID: 1, Status: model.StatusError, Error: expectErr}, expectErr).Return(nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk:        &model.Book{ID: 1, Status: model.StatusError, Error: nil},
			wantBk:    &model.Book{ID: 1, Status: model.StatusError, Error: fmt.Errorf("get book page failed: %w", serv.ErrUnavailable)},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "input is not a completely book (error is not nil)",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)

				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("", serv.ErrUnavailable)
				expectErr := fmt.Errorf("get book page failed: %w", serv.ErrUnavailable)
				rpo.EXPECT().SaveError(&model.Book{ID: 1, Status: model.StatusError, Error: expectErr}, expectErr).Return(nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk:        &model.Book{ID: 1, Status: model.StatusError, Error: serv.ErrUnavailable},
			wantBk:    &model.Book{ID: 1, Status: model.StatusError, Error: fmt.Errorf("get book page failed: %w", serv.ErrUnavailable)},
			wantError: serv.ErrUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := test.getServ(ctrl).ExploreBook(context.Background(), test.bk, nil)
			assert.Equal(t, test.wantBk, test.bk)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_Explore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		getServ   func(ctrl *gomock.Controller) *ServiceImpl
		wantError error
	}{
		{
			name: "early quit if explore existing book reaching limit",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)

				rpo.EXPECT().Stats().Return(repo.Summary{LatestSuccessID: 0, MaxBookID: 5})
				rpo.EXPECT().FindBookById(1).Return(&model.Book{ID: 1, Status: model.StatusError, Error: serv.ErrUnavailable}, nil)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("", serv.ErrUnavailable)
				expectErr := fmt.Errorf("get book page failed: %w", serv.ErrUnavailable)
				rpo.EXPECT().SaveError(&model.Book{ID: 1, Status: model.StatusError, Error: expectErr}, expectErr).Return(nil)

				return &ServiceImpl{
					rpo: rpo, vendorService: vendorService, cli: cli, sema: semaphore.NewWeighted(1),
					conf: config.SiteConfig{MaxExploreError: 1},
				}
			},
			wantError: nil,
		},
		{
			name: "quit if explore new book reaching limit",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)

				rpo.EXPECT().Stats().Return(repo.Summary{LatestSuccessID: 5, MaxBookID: 5})

				rpo.EXPECT().CreateBook(&model.Book{Site: "testing", ID: 6, HashCode: model.GenerateHash(), Status: model.StatusError}).Return(nil)
				vendorService.EXPECT().BookURL("6").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("", serv.ErrUnavailable)
				expectErr := fmt.Errorf("get book page failed: %w", serv.ErrUnavailable)
				rpo.EXPECT().SaveError(&model.Book{Site: "testing", ID: 6, HashCode: model.GenerateHash(), Status: model.StatusError, Error: expectErr}, expectErr).Return(nil)

				return &ServiceImpl{
					rpo: rpo, vendorService: vendorService, cli: cli, sema: semaphore.NewWeighted(1),
					name: "testing", conf: config.SiteConfig{MaxExploreError: 1},
				}
			},
			wantError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := test.getServ(ctrl).Explore(context.Background(), nil)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
