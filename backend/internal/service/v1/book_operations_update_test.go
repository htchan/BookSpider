package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	clientmock "github.com/htchan/BookSpider/internal/mock/client/v2"
	repomock "github.com/htchan/BookSpider/internal/mock/repo"
	vendormock "github.com/htchan/BookSpider/internal/mock/vendorservice"
	"github.com/htchan/BookSpider/internal/model"
	serv "github.com/htchan/BookSpider/internal/service"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func Test_isNewBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     *model.Book
		bkInfo *vendor.BookInfo
		want   bool
	}{
		{
			name:   "book is not new, because id not exist before",
			bk:     &model.Book{Status: model.Error},
			bkInfo: &vendor.BookInfo{Title: "title", Writer: "writer", Type: "type"},
			want:   false,
		},
		{
			name: "book is not new, because key fields was not updated",
			bk: &model.Book{
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type", Status: model.InProgress,
			},
			bkInfo: &vendor.BookInfo{
				Title: "title", Writer: "writer", Type: "type",
				UpdateChapter: "chapter", UpdateDate: "date",
			},
			want: false,
		},
		{
			name:   "book is new, because id existed before and title was updated",
			bk:     &model.Book{Title: "title", Status: model.InProgress},
			bkInfo: &vendor.BookInfo{Title: "title 2"},
			want:   true,
		},
		{
			name:   "book is new, because id existed before and writer was updated",
			bk:     &model.Book{Writer: model.Writer{Name: "writer"}, Status: model.InProgress},
			bkInfo: &vendor.BookInfo{Writer: "writer 2"},
			want:   true,
		},
		{
			name:   "book is new, because id existed before and type was updated",
			bk:     &model.Book{Type: "type", Status: model.InProgress},
			bkInfo: &vendor.BookInfo{Type: "type 2"},
			want:   true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := isNewBook(test.bk, test.bkInfo)
			assert.Equal(t, test.want, got)
		})
	}

}

func Test_isBookUpdated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     *model.Book
		bkInfo *vendor.BookInfo
		want   bool
	}{
		{
			name:   "book is not updated, because fields was not updated",
			bk:     &model.Book{UpdateDate: "date", UpdateChapter: "chapter"},
			bkInfo: &vendor.BookInfo{UpdateDate: "date", UpdateChapter: "chapter"},
			want:   false,
		},
		{
			name:   "book is updated, because update date was updated",
			bk:     &model.Book{UpdateDate: "date", UpdateChapter: "chapter"},
			bkInfo: &vendor.BookInfo{UpdateDate: "date 1", UpdateChapter: "chapter"},
			want:   true,
		},
		{
			name:   "book is updated, because update chapter was updated",
			bk:     &model.Book{UpdateDate: "date", UpdateChapter: "chapter"},
			bkInfo: &vendor.BookInfo{UpdateDate: "date", UpdateChapter: "chapter 1"},
			want:   true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := isBookUpdated(test.bk, test.bkInfo)
			assert.Equal(t, test.want, got)
		})

	}
}

func TestServiceImpl_UpdateBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		getServ   func(ctrl *gomock.Controller) *ServiceImpl
		bk        *model.Book
		wantBk    *model.Book
		wantError error
	}{
		{
			name: "no update for book with status error",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(nil, serv.ErrUnavailable)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk:        &model.Book{ID: 1, Status: model.Error},
			wantBk:    &model.Book{ID: 1, Status: model.Error},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "no update for existing book",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(&vendor.BookInfo{
					Title: "title", Writer: "writer", Type: "type", UpdateChapter: "chapter", UpdateDate: "date",
				}, nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk: &model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
			},
			wantBk: &model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
			},
			wantError: nil,
		},
		{
			name: "update book with status error",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(&vendor.BookInfo{
					Title: "title", Writer: "writer", Type: "type", UpdateChapter: "chapter", UpdateDate: "date",
				}, nil)
				bk := &model.Book{
					ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
				}
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)
				rpo.EXPECT().UpdateBook(bk).Return(nil)
				rpo.EXPECT().SaveError(bk, nil).Return(nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk: &model.Book{ID: 1, Status: model.Error},
			wantBk: &model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
			},
			wantError: nil,
		},
		{
			name: "update existing book",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(&vendor.BookInfo{
					Title: "title", Writer: "writer", Type: "type", UpdateChapter: "chapter 2", UpdateDate: "date 2",
				}, nil)
				bk := &model.Book{
					ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date 2", UpdateChapter: "chapter 2", Status: model.InProgress,
				}
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)
				rpo.EXPECT().UpdateBook(bk).Return(nil)
				rpo.EXPECT().SaveError(bk, nil).Return(nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk: &model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
			},
			wantBk: &model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date 2", UpdateChapter: "chapter 2", Status: model.InProgress,
			},
			wantError: nil,
		},
		{
			name: "create new books",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(&vendor.BookInfo{
					Title: "title 2", Writer: "writer", Type: "type", UpdateChapter: "chapter", UpdateDate: "date",
				}, nil)
				bk := &model.Book{
					ID: 1, HashCode: model.GenerateHash(), Title: "title 2", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
				}
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)
				rpo.EXPECT().CreateBook(bk).Return(nil)
				rpo.EXPECT().SaveError(bk, nil).Return(nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk: &model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
			},
			wantBk: &model.Book{ID: 1, HashCode: model.GenerateHash(), Title: "title 2", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
			},
			wantError: nil,
		},
		{
			name: "getting error when sending request",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("", serv.ErrUnavailable)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk:        &model.Book{ID: 1, Status: model.Error},
			wantBk:    &model.Book{ID: 1, Status: model.Error},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "getting error when parsing book",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(nil, serv.ErrUnavailable)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk:        &model.Book{ID: 1, Status: model.Error},
			wantBk:    &model.Book{ID: 1, Status: model.Error},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "getting error when updating book",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(&vendor.BookInfo{
					Title: "title", Writer: "writer", Type: "type", UpdateChapter: "chapter", UpdateDate: "date",
				}, nil)
				bk := &model.Book{
					ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
				}
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)
				rpo.EXPECT().UpdateBook(bk).Return(serv.ErrUnavailable)
				rpo.EXPECT().SaveError(bk, nil).Return(nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk: &model.Book{ID: 1, Status: model.Error},
			wantBk: &model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
			},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "getting error when creating book",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(&vendor.BookInfo{
					Title: "title", Writer: "writer", Type: "type", UpdateChapter: "chapter", UpdateDate: "date",
				}, nil)
				bk := &model.Book{
					ID: 1, HashCode: model.GenerateHash(), Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
				}
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)
				rpo.EXPECT().CreateBook(bk).Return(serv.ErrUnavailable)
				rpo.EXPECT().SaveError(bk, nil).Return(nil)

				return &ServiceImpl{rpo: rpo, vendorService: vendorService, cli: cli}
			},
			bk: &model.Book{ID: 1, Status: model.InProgress},
			wantBk: &model.Book{ID: 1, HashCode: model.GenerateHash(), Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
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

			err := test.getServ(ctrl).UpdateBook(context.Background(), test.bk)
			assert.Equal(t, test.wantBk, test.bk)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		getServ   func(ctrl *gomock.Controller) *ServiceImpl
		wantError error
	}{
		{
			name: "update existing book",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo, cli := repomock.NewMockRepository(ctrl), clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)
				bk := model.Book{
					ID: 1, Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "date", UpdateChapter: "chapter", Status: model.InProgress,
				}
				ch := make(chan model.Book)

				go func() {
					bk := bk
					ch <- bk
					close(ch)
				}()

				rpo.EXPECT().FindBooksForUpdate().Return(ch, nil)
				vendorService.EXPECT().BookURL("1").Return("https://test.com")
				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("response", nil)
				vendorService.EXPECT().ParseBook("response").Return(&vendor.BookInfo{
					Title: "title", Writer: "writer", Type: "type", UpdateChapter: "chapter 2", UpdateDate: "date 2",
				}, nil)

				bkUpdated := bk
				bkUpdated.UpdateDate, bkUpdated.UpdateChapter = "date 2", "chapter 2"

				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)
				rpo.EXPECT().UpdateBook(&bkUpdated).Return(nil)
				rpo.EXPECT().SaveError(&bkUpdated, nil).Return(nil)

				return &ServiceImpl{sema: semaphore.NewWeighted(1), rpo: rpo, vendorService: vendorService, cli: cli}
			},
			wantError: nil,
		},
		{
			name: "return error if find book for update failed",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)

				rpo.EXPECT().FindBooksForUpdate().Return(nil, serv.ErrUnavailable)

				return &ServiceImpl{rpo: rpo}
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

			err := test.getServ(ctrl).Update(context.Background())
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
