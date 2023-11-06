package service

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/htchan/BookSpider/internal/config/v2"
	clientmock "github.com/htchan/BookSpider/internal/mock/client/v2"
	repomock "github.com/htchan/BookSpider/internal/mock/repo"
	vendormock "github.com/htchan/BookSpider/internal/mock/vendorservice"
	"github.com/htchan/BookSpider/internal/model"
	serv "github.com/htchan/BookSpider/internal/service"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestServiceImpl_downloadChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		getService  func(ctrl *gomock.Controller) *ServiceImpl
		chapter     *model.Chapter
		wantChapter *model.Chapter
		wantError   error
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("chapter response", nil)
				vendorService.EXPECT().ParseChapter("chapter response").Return(&vendor.ChapterInfo{
					Title: "title", Body: "content content content",
				}, nil)

				return &ServiceImpl{cli: cli, vendorService: vendorService}
			},
			chapter: &model.Chapter{
				Index: 1, URL: "https://test.com",
			},
			wantChapter: &model.Chapter{
				Index: 1, URL: "https://test.com",
				Title: "title", Content: "content content content",
			},
			wantError: nil,
		},
		{
			name: "fail to send request",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				cli := clientmock.NewMockBookClient(ctrl)

				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("", serv.ErrUnavailable)

				return &ServiceImpl{cli: cli}
			},
			chapter: &model.Chapter{
				Index: 1, URL: "https://test.com",
			},
			wantChapter: &model.Chapter{
				Index: 1, URL: "https://test.com", Error: serv.ErrUnavailable,
			},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "fail to parse chapter",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				cli.EXPECT().Get(gomock.Any(), "https://test.com").Return("chapter response", nil)
				vendorService.EXPECT().ParseChapter("chapter response").Return(nil, serv.ErrUnavailable)

				return &ServiceImpl{cli: cli, vendorService: vendorService}
			},
			chapter: &model.Chapter{
				Index: 1, URL: "https://test.com",
			},
			wantChapter: &model.Chapter{
				Index: 1, URL: "https://test.com", Error: serv.ErrUnavailable,
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

			err := test.getService(ctrl).downloadChapter(context.Background(), test.chapter)
			assert.ErrorIs(t, err, test.wantError)
			assert.Equal(t, test.wantChapter, test.chapter)
		})
	}
}

func TestServiceImpl_DownloadBook(t *testing.T) {
	t.Parallel()

	os.Mkdir("./download-book", 0755)

	t.Cleanup(func() { os.RemoveAll("./download-book") })

	tests := []struct {
		name                 string
		getService           func(ctrl *gomock.Controller) *ServiceImpl
		book                 *model.Book
		wantBook             *model.Book
		wantError            error
		wantBookFileLocation string
		wantBookContent      string
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				vendorService.EXPECT().ChapterListURL("1").Return("https://test.com/chapter-list")
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter-list").Return("chapter list response", nil)
				vendorService.EXPECT().ParseChapterList("1", "chapter list response").Return(vendor.ChapterList{
					{URL: "https://test.com/chapter/1", Title: "title 1"},
					{URL: "https://test.com/chapter/2", Title: "title 2"},
				}, nil)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/1").Return("chapter 1 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 1 response").Return(&vendor.ChapterInfo{
					Title: "chapter title 1", Body: "content 1 content 1 content 1",
				}, nil)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/2").Return("chapter 2 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 2 response").Return(&vendor.ChapterInfo{
					Title: "chapter title 2", Body: "content 2 content 2 content 2",
				}, nil)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
					Status: model.StatusEnd, IsDownloaded: true,
				}).Return(nil)

				return &ServiceImpl{
					conf: config.SiteConfig{Storage: "./download-book"}, sema: semaphore.NewWeighted(1),
					rpo: rpo, cli: cli, vendorService: vendorService,
				}
			},
			book: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantBook: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: true,
			},
			wantError:            nil,
			wantBookFileLocation: "./download-book/1.txt",
			wantBookContent: `title 1
writer 1
--------------------

chapter title 1
--------------------
content 1 content 1 content 1
--------------------
chapter title 2
--------------------
content 2 content 2 content 2
--------------------
`,
		},
		{
			name: "book status is not end",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				return &ServiceImpl{}
			},
			book:      &model.Book{Status: model.StatusInProgress},
			wantBook:  &model.Book{Status: model.StatusInProgress},
			wantError: serv.ErrBookStatusNotEnd,
		},
		{
			name: "book already downloaded",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				return &ServiceImpl{}
			},
			book:      &model.Book{Status: model.StatusEnd, IsDownloaded: true},
			wantBook:  &model.Book{Status: model.StatusEnd, IsDownloaded: true},
			wantError: serv.ErrBookAlreadyDownloaded,
		},
		{
			name: "fail to send request",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				vendorService.EXPECT().ChapterListURL("1").Return("https://test.com/chapter-list")
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter-list").Return("", serv.ErrUnavailable)

				return &ServiceImpl{
					cli: cli, vendorService: vendorService, sema: semaphore.NewWeighted(1),
				}
			},
			book: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantBook: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "fail to parse chapter list",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				vendorService.EXPECT().ChapterListURL("1").Return("https://test.com/chapter-list")
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter-list").Return("chapter list response", nil)
				vendorService.EXPECT().ParseChapterList("1", "chapter list response").Return(nil, serv.ErrUnavailable)

				return &ServiceImpl{
					conf: config.SiteConfig{Storage: "./download-book"}, sema: semaphore.NewWeighted(1),
					cli: cli, vendorService: vendorService,
				}
			},
			book: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantBook: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantError: serv.ErrUnavailable,
		},
		{
			name: "download chapter fails reach threshold",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				vendorService.EXPECT().ChapterListURL("1").Return("https://test.com/chapter-list")
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter-list").Return("chapter list response", nil)
				vendorService.EXPECT().ParseChapterList("1", "chapter list response").Return(vendor.ChapterList{
					{URL: "https://test.com/chapter/1", Title: "title 1"},
					{URL: "https://test.com/chapter/2", Title: "title 2"},
				}, nil)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/1").Return("chapter 1 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 1 response").Return(nil, serv.ErrUnavailable)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/2").Return("chapter 2 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 2 response").Return(nil, serv.ErrUnavailable)

				return &ServiceImpl{
					conf: config.SiteConfig{Storage: "./download-book"}, sema: semaphore.NewWeighted(1),
					rpo: rpo, cli: cli, vendorService: vendorService,
				}
			},
			book: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantBook: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantError: serv.ErrTooManyFailedChapters,
		},
		{
			name: "fail to update book after download",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				vendorService.EXPECT().ChapterListURL("1").Return("https://test.com/chapter-list")
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter-list").Return("chapter list response", nil)
				vendorService.EXPECT().ParseChapterList("1", "chapter list response").Return(vendor.ChapterList{
					{URL: "https://test.com/chapter/1", Title: "title 1"},
					{URL: "https://test.com/chapter/2", Title: "title 2"},
				}, nil)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/1").Return("chapter 1 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 1 response").Return(&vendor.ChapterInfo{
					Title: "chapter title 1", Body: "content 1 content 1 content 1",
				}, nil)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/2").Return("chapter 2 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 2 response").Return(&vendor.ChapterInfo{
					Title: "chapter title 2", Body: "content 2 content 2 content 2",
				}, nil)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
					Status: model.StatusEnd, IsDownloaded: true,
				}).Return(serv.ErrUnavailable)

				return &ServiceImpl{
					conf: config.SiteConfig{Storage: "./download-book"}, sema: semaphore.NewWeighted(1),
					rpo: rpo, cli: cli, vendorService: vendorService,
				}
			},
			book: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: false,
			},
			wantBook: &model.Book{
				ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"},
				Status: model.StatusEnd, IsDownloaded: true,
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

			err := test.getService(ctrl).DownloadBook(context.Background(), test.book)
			assert.ErrorIs(t, err, test.wantError)
			assert.Equal(t, test.wantBook, test.book)

			if test.wantBookFileLocation != "" {
				content, err := os.ReadFile(test.wantBookFileLocation)
				assert.NoError(t, err)
				assert.Equal(t, test.wantBookContent, string(content))
			}
		})
	}
}

func TestServiceImpl_Download(t *testing.T) {
	t.Parallel()

	os.MkdirAll("./download", 0755)

	t.Cleanup(func() { os.RemoveAll("./download") })

	tests := []struct {
		name       string
		getService func(ctrl *gomock.Controller) *ServiceImpl
		wantError  error
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				cli := clientmock.NewMockBookClient(ctrl)
				vendorService := vendormock.NewMockVendorService(ctrl)

				ch := make(chan model.Book)
				go func() {
					ch <- model.Book{ID: 1, Title: "title 1", Writer: model.Writer{Name: "writer 1"}, Status: model.StatusEnd}
					ch <- model.Book{ID: 2, Title: "title 2", Writer: model.Writer{Name: "writer 2"}, Status: model.StatusEnd}
					close(ch)
				}()

				rpo.EXPECT().FindBooksForDownload().Return(ch, nil)
				vendorService.EXPECT().ChapterListURL("1").Return("https://test.com/chapter-list-1")
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter-list-1").Return("", serv.ErrUnavailable)
				vendorService.EXPECT().ChapterListURL("2").Return("https://test.com/chapter-list-2")
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter-list-2").Return("chapter list response", nil)
				vendorService.EXPECT().ParseChapterList("2", "chapter list response").Return(vendor.ChapterList{
					{URL: "https://test.com/chapter/1", Title: "title 1"},
					{URL: "https://test.com/chapter/2", Title: "title 2"},
				}, nil)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/1").Return("chapter 1 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 1 response").Return(&vendor.ChapterInfo{
					Title: "chapter title 1", Body: "content 1 content 1 content 1",
				}, nil)
				cli.EXPECT().Get(gomock.Any(), "https://test.com/chapter/2").Return("chapter 2 response", nil)
				vendorService.EXPECT().ParseChapter("chapter 2 response").Return(&vendor.ChapterInfo{
					Title: "chapter title 2", Body: "content 2 content 2 content 2",
				}, nil)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 2, Title: "title 2", Writer: model.Writer{Name: "writer 2"},
					Status: model.StatusEnd, IsDownloaded: true,
				}).Return(serv.ErrUnavailable)

				return &ServiceImpl{
					conf: config.SiteConfig{Storage: "./download-book", MaxDownloadConcurrency: 1},
					sema: semaphore.NewWeighted(2),
					rpo:  rpo, cli: cli, vendorService: vendorService,
				}
			},
			wantError: nil,
		},
		{
			name: "fail to find books for download",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksForDownload().Return(nil, serv.ErrUnavailable)

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

			err := test.getService(ctrl).Download(context.Background())
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
