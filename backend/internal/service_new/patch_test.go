package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	config "github.com/htchan/BookSpider/internal/config_new"
	mockclient "github.com/htchan/BookSpider/internal/mock/client/v2"
	mockparser "github.com/htchan/BookSpider/internal/mock/parser"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestServiceImp_checkBookStorate(t *testing.T) {
	t.Parallel()
	os.Mkdir("test-check-book-storage", os.ModePerm)
	os.WriteFile("./test-check-book-storage/1.txt", []byte("some content"), os.ModeAppend)

	t.Cleanup(func() {
		os.RemoveAll("test-check-book-storage")
	})

	tests := []struct {
		name   string
		serv   ServiceImp
		bk     *model.Book
		wantBk *model.Book
		want   bool
	}{
		{
			name:   "book download + file exist",
			serv:   ServiceImp{conf: config.SiteConfig{Storage: "test-check-book-storage"}},
			bk:     &model.Book{ID: 1, IsDownloaded: true},
			wantBk: &model.Book{ID: 1, IsDownloaded: true},
			want:   false,
		},
		{
			name:   "book not download + file not exist",
			serv:   ServiceImp{conf: config.SiteConfig{Storage: "test-check-book-storage"}},
			bk:     &model.Book{ID: 2, IsDownloaded: false},
			wantBk: &model.Book{ID: 2, IsDownloaded: false},
			want:   false,
		},
		{
			name:   "book download + file not exist",
			serv:   ServiceImp{conf: config.SiteConfig{Storage: "test-check-book-storage"}},
			bk:     &model.Book{ID: 2, IsDownloaded: true},
			wantBk: &model.Book{ID: 2, IsDownloaded: false},
			want:   true,
		},
		{
			name:   "book not download + file exist",
			serv:   ServiceImp{conf: config.SiteConfig{Storage: "test-check-book-storage"}},
			bk:     &model.Book{ID: 1, IsDownloaded: false},
			wantBk: &model.Book{ID: 1, IsDownloaded: true},
			want:   true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.serv.checkBookStorage(test.bk)
			assert.Equal(t, test.wantBk, test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImp_PatchDownloadStatus(t *testing.T) {
	t.Parallel()
	os.Mkdir("./test-patch-status", os.ModePerm)
	os.WriteFile("./test-patch-status/1.txt", []byte("data"), os.ModeAppend)
	os.WriteFile("./test-patch-status/3.txt", []byte("data"), os.ModeAppend)
	t.Cleanup(func() {
		os.RemoveAll("./test-patch-status")
	})

	tests := []struct {
		name         string
		setupServ    func(ctrl *gomock.Controller) ServiceImp
		wantError    bool
		wantErrorStr string
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)

				bookChan := make(chan model.Book, 4)
				bookChan <- model.Book{ID: 1, IsDownloaded: false}
				bookChan <- model.Book{ID: 2, IsDownloaded: true}
				bookChan <- model.Book{ID: 3, IsDownloaded: true}
				bookChan <- model.Book{ID: 4, IsDownloaded: false}
				close(bookChan)

				rpo.EXPECT().FindAllBooks().Return(bookChan, nil)
				rpo.EXPECT().UpdateBook(&model.Book{ID: 1, IsDownloaded: true})
				rpo.EXPECT().UpdateBook(&model.Book{ID: 2, IsDownloaded: false})
				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					conf:   config.SiteConfig{Storage: "./test-patch-status"},
				}
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "fail to find all books",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)

				rpo.EXPECT().FindAllBooks().Return(nil, errors.New("some error"))
				return ServiceImp{
					rpo:  rpo,
					conf: config.SiteConfig{Storage: "./test-patch-status"},
				}
			},
			wantError:    true,
			wantErrorStr: "patch download status fail: some error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)

			err := serv.PatchDownloadStatus()
			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}
		})
	}
}

func TestServiceImp_PatchMissingRecords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupServ    func(ctrl *gomock.Controller) ServiceImp
		wantError    bool
		wantErrorStr string
	}{
		{
			name: "no records missing",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().Stats().Return(repo.Summary{
					MaxBookID: 5,
				})
				for i := 1; i <= 5; i++ {
					rpo.EXPECT().FindBookById(i).Return(&model.Book{ID: i}, nil)
				}

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					conf:   config.SiteConfig{Storage: "./test-patch-status"},
				}
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "patch missing record in middle",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/5").Return("content", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("content").Return(nil, parse.ErrParseBookFieldsNotFound)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().Stats().Return(repo.Summary{
					MaxBookID: 5,
				})
				rpo.EXPECT().FindBookById(1).Return(&model.Book{ID: 1}, nil)
				rpo.EXPECT().FindBookById(2).Return(&model.Book{ID: 2}, nil)
				rpo.EXPECT().FindBookById(3).Return(&model.Book{ID: 3}, nil)
				rpo.EXPECT().FindBookById(4).Return(nil, sql.ErrConnDone)
				rpo.EXPECT().FindBookById(5).Return(nil, fmt.Errorf("fail to query book by site id: %w", repo.ErrBookNotExist))

				err := fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound)
				rpo.EXPECT().CreateBook(&model.Book{Site: "test-patch-missing-records", ID: 5, HashCode: model.GenerateHash()})
				rpo.EXPECT().SaveError(&model.Book{Site: "test-patch-missing-records", ID: 5, HashCode: model.GenerateHash(), Error: err}, err)

				return ServiceImp{
					name:   "test-patch-missing-records",
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					parser: p,
					conf: config.SiteConfig{
						URL:     config.URLConfig{Base: "http://test.com/%v"},
						Storage: "./test-patch-status",
					},
				}
			},
			wantError:    false,
			wantErrorStr: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)

			err := serv.PatchMissingRecords()
			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}
		})
	}
}
