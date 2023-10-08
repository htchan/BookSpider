package service

import (
	"context"
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
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestServiceImp_downloadURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		serv   ServiceImp
		bk     *model.Book
		expect string
	}{
		{
			name: "populate download url for book",
			serv: ServiceImp{conf: config.SiteConfig{
				URL: config.URLConfig{Download: "some test url %v"},
			}},
			bk:     &model.Book{ID: 123},
			expect: "some test url 123",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			url := test.serv.downloadURL(test.bk)
			assert.Equal(t, url, test.expect)
		})
	}
}

func TestServiceImp_chapterURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		serv   ServiceImp
		bk     *model.Book
		ch     *model.Chapter
		expect string
	}{
		{
			name: "chapter url start with http",
			serv: ServiceImp{conf: config.SiteConfig{URL: config.URLConfig{
				Download:      "some test url %v",
				ChapterPrefix: "some test chapter prefix",
			}}},
			bk:     &model.Book{ID: 123},
			ch:     &model.Chapter{URL: "http://test.com"},
			expect: "http://test.com",
		},
		{
			name: "chapter url start with https",
			serv: ServiceImp{conf: config.SiteConfig{URL: config.URLConfig{
				Download:      "some test url %v",
				ChapterPrefix: "some test chapter prefix",
			}}},
			bk:     &model.Book{ID: 123},
			ch:     &model.Chapter{URL: "https://test.com"},
			expect: "https://test.com",
		},
		{
			name: "chapter url start with /",
			serv: ServiceImp{conf: config.SiteConfig{URL: config.URLConfig{
				Download:      "some test url %v",
				ChapterPrefix: "some test chapter prefix",
			}}},
			bk:     &model.Book{ID: 123},
			ch:     &model.Chapter{URL: "/data"},
			expect: "some test chapter prefix/data",
		},
		{
			name: "chapter url start with some random string",
			serv: ServiceImp{conf: config.SiteConfig{URL: config.URLConfig{
				Download:      "some test url %v",
				ChapterPrefix: "some test chapter prefix",
			}}},
			bk:     &model.Book{ID: 123},
			ch:     &model.Chapter{URL: "abc"},
			expect: "some test url 123/abc",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			url := test.serv.chapterURL(test.bk, test.ch)
			assert.Equal(t, url, test.expect)
		})
	}
}

func TestServiceImp_downloadChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupServ      func(ctrl *gomock.Controller) ServiceImp
		bk             *model.Book
		ch             *model.Chapter
		expectBook     *model.Book
		expectChapter  *model.Chapter
		expectedError  bool
		expectErrorStr string
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/chapters/1").Return("some result", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapter("some result").Return(parse.NewParsedChapterFields(
					"parsed title",
					"parsed content 1   \n\n\n\n   parsed content 2",
				), nil)

				return ServiceImp{
					client: c,
					parser: p,
				}
			},
			bk:         &model.Book{ID: 1},
			ch:         &model.Chapter{URL: "http://test.com/chapters/1"},
			expectBook: &model.Book{},
			expectChapter: &model.Chapter{
				URL:     "http://test.com/chapters/1",
				Title:   "parsed title",
				Content: "parsed content 1\n parsed content 2",
			},
			expectedError:  false,
			expectErrorStr: "",
		},
		{
			name: "fail to send request",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/chapters/1").Return("", errors.New("new error"))

				p := mockparser.NewMockParser(ctrl)

				return ServiceImp{
					client: c,
					parser: p,
				}
			},
			bk:         &model.Book{ID: 1},
			ch:         &model.Chapter{URL: "http://test.com/chapters/1"},
			expectBook: &model.Book{},
			expectChapter: &model.Chapter{
				URL: "http://test.com/chapters/1",
			},
			expectedError:  true,
			expectErrorStr: "fetch chapter html fail: new error",
		},
		{
			name: "fail to parse chapter",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/chapters/1").Return("some result", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapter("some result").Return(nil, errors.New("some error"))

				return ServiceImp{
					client: c,
					parser: p,
				}
			},
			bk:         &model.Book{ID: 1},
			ch:         &model.Chapter{URL: "http://test.com/chapters/1"},
			expectBook: &model.Book{},
			expectChapter: &model.Chapter{
				URL: "http://test.com/chapters/1",
			},
			expectedError:  true,
			expectErrorStr: "parse chapter html fail: some error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			err := serv.downloadChapter(test.bk, test.ch)

			if test.expectedError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.expectErrorStr)
			}
		})
	}
}

func TestServiceImp_downloadChapterList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupServ      func(ctrl *gomock.Controller) ServiceImp
		bk             *model.Book
		expectBook     *model.Book
		expectChapters model.Chapters
		expectedError  bool
		expectErrorStr string
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/1/download").Return("chapter list", nil)
				c.EXPECT().Get(gomock.Any(), "http://test.com/bk/1/chapter/1").Return("chapter 1", nil)

				parsedChapterList := parse.ParsedChapterList{}
				parsedChapterList.Append("http://test.com/bk/1/chapter/1", "test chapter 1")

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapterList("chapter list").Return(&parsedChapterList, nil)
				p.EXPECT().ParseChapter("chapter 1").Return(
					parse.NewParsedChapterFields("title 1", "content 1"), nil,
				)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					parser: p,
					conf: config.SiteConfig{
						URL: config.URLConfig{Download: "http://test.com/book/%v/download"},
					},
				}
			},
			bk:         &model.Book{ID: 1},
			expectBook: &model.Book{ID: 1},
			expectChapters: model.Chapters{
				{URL: "http://test.com/bk/1/chapter/1", Title: "title 1", Content: "content 1"},
			},
			expectedError:  false,
			expectErrorStr: "",
		},
		{
			name: "fail to send request",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/1/download").Return("", errors.New("some error"))

				p := mockparser.NewMockParser(ctrl)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					parser: p,
					conf: config.SiteConfig{
						URL: config.URLConfig{Download: "http://test.com/book/%v/download"},
					},
				}
			},
			bk:             &model.Book{ID: 1},
			expectBook:     &model.Book{ID: 1},
			expectChapters: nil,
			expectedError:  true,
			expectErrorStr: "fetch chapter list html fail: some error",
		},
		{
			name: "fail to parse chapter list",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/1/download").Return("chapter list", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapterList("chapter list").Return(nil, errors.New("some error"))

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					parser: p,
					conf: config.SiteConfig{
						URL: config.URLConfig{Download: "http://test.com/book/%v/download"},
					},
				}
			},
			bk:             &model.Book{ID: 1},
			expectBook:     &model.Book{ID: 1},
			expectChapters: nil,
			expectedError:  true,
			expectErrorStr: "parse chapter list html fail: some error",
		},
		{
			name: "some download chapter failed",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/1/download").Return("chapter list", nil)
				c.EXPECT().Get(gomock.Any(), "http://test.com/bk/1/chapter/1").Return("chapter 1", nil)
				c.EXPECT().Get(gomock.Any(), "http://test.com/bk/1/chapter/2").Return("chapter 2", nil)

				parsedChapterList := parse.ParsedChapterList{}
				parsedChapterList.Append("http://test.com/bk/1/chapter/1", "test chapter 1")
				parsedChapterList.Append("http://test.com/bk/1/chapter/2", "test chapter 2")

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapterList("chapter list").Return(&parsedChapterList, nil)
				p.EXPECT().ParseChapter("chapter 1").Return(
					nil, errors.New("chapter 1 fail"),
				)
				p.EXPECT().ParseChapter("chapter 2").Return(
					parse.NewParsedChapterFields("title 2", "content 2"), nil,
				)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					parser: p,
					conf: config.SiteConfig{
						URL: config.URLConfig{Download: "http://test.com/book/%v/download"},
					},
				}
			},
			bk:         &model.Book{ID: 1},
			expectBook: &model.Book{ID: 1},
			expectChapters: model.Chapters{
				{
					Index: 0,
					URL:   "http://test.com/bk/1/chapter/1",
					Title: "test chapter 1",
					Error: fmt.Errorf("parse chapter html fail: %w", errors.New("chapter 1 fail")),
				},
				{
					Index:   1,
					URL:     "http://test.com/bk/1/chapter/2",
					Title:   "title 2",
					Content: "content 2",
				},
			},
			expectedError:  false,
			expectErrorStr: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			chapters, err := serv.downloadChapterList(test.bk)

			if test.expectedError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.expectErrorStr)
			}

			assert.Equal(t, test.expectChapters, chapters)
		})
	}
}

func TestServiceImp_saveContent(t *testing.T) {
	t.Parallel()

	os.Mkdir("test-save-content", os.ModePerm)
	t.Cleanup(func() {
		os.RemoveAll("test-save-content")
	})

	tests := []struct {
		name               string
		serv               ServiceImp
		location           string
		bk                 *model.Book
		chapters           model.Chapters
		expectFileLocation string
		expectContent      string
		expectedError      bool
		expectErrorStr     string
	}{
		{
			name:     "happy flow",
			serv:     ServiceImp{},
			location: "./test-save-content/1.txt",
			bk:       &model.Book{Title: "title", Writer: model.NewWriter("writer")},
			chapters: model.Chapters{
				{Title: "chapter title", Content: "chapter content"},
			},
			expectFileLocation: "./test-save-content/1.txt",
			expectContent: `title
writer
--------------------

chapter title
--------------------
chapter content
--------------------
`,
			expectedError:  false,
			expectErrorStr: "",
		},
		{
			name:     "location not exist",
			serv:     ServiceImp{},
			location: "./test-save-content/test/2.txt",
			bk:       &model.Book{Title: "title", Writer: model.NewWriter("writer")},
			chapters: model.Chapters{
				{Title: "chapter title", Content: "chapter content"},
			},
			expectFileLocation: "",
			expectContent:      "",
			expectedError:      true,
			expectErrorStr:     "Save book fail: open ./test-save-content/test/2.txt: no such file or directory",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.serv.saveContent(test.location, test.bk, test.chapters)

			if test.expectedError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.expectErrorStr)
			}

			if test.expectContent != "" {
				if assert.FileExists(t, test.expectFileLocation) {
					defer os.Remove(test.expectFileLocation)
				}
				b, err := os.ReadFile(test.expectFileLocation)
				if err != nil {
					t.Errorf("read file error: %v", err)
				}
				assert.Equal(t, test.expectContent, string(b))
			}
		})
	}
}

func TestServiceImp_DownloadBook(t *testing.T) {
	t.Parallel()
	os.Mkdir("test-download-book", os.ModePerm)
	t.Cleanup(func() {
		os.RemoveAll("test-download-book")
	})

	tests := []struct {
		name               string
		setupServ          func(ctrl *gomock.Controller) ServiceImp
		bk                 *model.Book
		expectBk           *model.Book
		expectFileLocation string
		expectContent      string
		expectedError      bool
		expectErrorStr     string
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/1/download").Return("chapter list", nil)
				c.EXPECT().Get(gomock.Any(), "http://test.com/bk/1/chapter/1").Return("chapter 1", nil)

				parsedChapterList := parse.ParsedChapterList{}
				parsedChapterList.Append("http://test.com/bk/1/chapter/1", "test chapter 1")

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapterList("chapter list").Return(&parsedChapterList, nil)
				p.EXPECT().ParseChapter("chapter 1").Return(
					parse.NewParsedChapterFields("title 1", "content 1"), nil,
				)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(gomock.Any()).Return(nil)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:     config.URLConfig{Download: "http://test.com/book/%v/download"},
						Storage: "./test-download-book/",
					},
				}
			},
			bk: &model.Book{
				ID:     1,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.End,
			},
			expectBk: &model.Book{
				ID:           1,
				Title:        "title",
				Writer:       model.Writer{Name: "writer"},
				Status:       model.End,
				IsDownloaded: true,
			},
			expectFileLocation: "./test-download-book/1.txt",
			expectContent:      "title\nwriter\n--------------------\n\ntitle 1\n--------------------\ncontent 1\n--------------------\n",
			expectedError:      false,
			expectErrorStr:     "",
		},
		{
			name: "book is not end",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				return ServiceImp{}
			},
			bk: &model.Book{
				ID:     2,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.InProgress,
			},
			expectBk: &model.Book{
				ID:     2,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.InProgress,
			},
			expectFileLocation: "./test-download-book/2.txt",
			expectContent:      "",
			expectedError:      true,
			expectErrorStr:     "book status not ready for download. status: INPROGRESS",
		},
		{
			name: "book already downloaded",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				return ServiceImp{}
			},
			bk: &model.Book{
				ID:           3,
				Title:        "title",
				Writer:       model.Writer{Name: "writer"},
				Status:       model.End,
				IsDownloaded: true,
			},
			expectBk: &model.Book{
				ID:           3,
				Title:        "title",
				Writer:       model.Writer{Name: "writer"},
				Status:       model.End,
				IsDownloaded: true,
			},
			expectFileLocation: "./test-download-book/3.txt",
			expectContent:      "",
			expectedError:      true,
			expectErrorStr:     "book was downloaded",
		},
		{
			name: "fail to fetch chapter list",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/4/download").Return("", errors.New("url invalid"))

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					conf: config.SiteConfig{
						URL:     config.URLConfig{Download: "http://test.com/book/%v/download"},
						Storage: "./test-download-book/",
					},
				}
			},
			bk: &model.Book{
				ID:     4,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.End,
			},
			expectBk: &model.Book{
				ID:     4,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.End,
			},
			expectFileLocation: "./test-download-book/4.txt",
			expectContent:      "",
			expectedError:      true,
			expectErrorStr:     "Download chapters fail: fetch chapter list html fail: url invalid",
		},
		{
			name: "fail to save content to file",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/5/download").Return("chapter list", nil)
				c.EXPECT().Get(gomock.Any(), "http://test.com/bk/5/chapter/1").Return("chapter 1", nil)

				parsedChapterList := parse.ParsedChapterList{}
				parsedChapterList.Append("http://test.com/bk/5/chapter/1", "test chapter 1")

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapterList("chapter list").Return(&parsedChapterList, nil)
				p.EXPECT().ParseChapter("chapter 1").Return(
					parse.NewParsedChapterFields("title 1", "content 1"), nil,
				)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					parser: p,
					conf: config.SiteConfig{
						URL:     config.URLConfig{Download: "http://test.com/book/%v/download"},
						Storage: "./test-download-book/not-exist",
					},
				}
			},
			bk: &model.Book{
				ID:     5,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.End,
			},
			expectBk: &model.Book{
				ID:     5,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.End,
			},
			expectFileLocation: "./test-download-book/5.txt",
			expectContent:      "",
			expectedError:      true,
			expectErrorStr:     "save content fail: Save book fail: open test-download-book/not-exist/5.txt: no such file or directory",
		},
		{
			name: "fail to update book in DB",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/6/download").Return("chapter list", nil)
				c.EXPECT().Get(gomock.Any(), "http://test.com/bk/6/chapter/1").Return("chapter 1", nil)

				parsedChapterList := parse.ParsedChapterList{}
				parsedChapterList.Append("http://test.com/bk/6/chapter/1", "test chapter 1")

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseChapterList("chapter list").Return(&parsedChapterList, nil)
				p.EXPECT().ParseChapter("chapter 1").Return(
					parse.NewParsedChapterFields("title 1", "content 1"), nil,
				)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(gomock.Any()).Return(errors.New("update bk error"))

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:     config.URLConfig{Download: "http://test.com/book/%v/download"},
						Storage: "./test-download-book/",
					},
				}
			},
			bk: &model.Book{
				ID:     6,
				Title:  "title",
				Writer: model.Writer{Name: "writer"},
				Status: model.End,
			},
			expectBk: &model.Book{
				ID:           6,
				Title:        "title",
				Writer:       model.Writer{Name: "writer"},
				Status:       model.End,
				IsDownloaded: true,
			},
			expectFileLocation: "./test-download-book/6.txt",
			expectContent:      "",
			expectedError:      true,
			expectErrorStr:     "update book fail: update bk error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			err := serv.DownloadBook(test.bk)

			if test.expectedError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.expectErrorStr)
			}

			if test.expectContent != "" {
				if assert.FileExists(t, test.expectFileLocation) {
					defer os.Remove(test.expectFileLocation)
				}
				b, err := os.ReadFile(test.expectFileLocation)
				if err != nil {
					t.Errorf("read file error: %v", err)
				}
				assert.Equal(t, test.expectContent, string(b))
			}
			assert.Equal(t, test.bk, test.expectBk)
		})
	}
}

func TestServiceImp_Download(t *testing.T) {
	t.Parallel()

	os.Mkdir("./test-download", os.ModePerm)
	t.Cleanup(func() {
		os.Remove("./test-download")
	})

	type expectFile struct {
		expectFileLocation string
		expectContent      string
	}

	tests := []struct {
		name           string
		setupServ      func(ctrl *gomock.Controller) ServiceImp
		expectFiles    []expectFile
		expectedError  bool
		expectErrorStr string
	}{
		{
			name: "some book download success and some book failed",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {

				c := mockclient.NewMockBookClient(ctrl)
				p := mockparser.NewMockParser(ctrl)
				rpo := mockrepo.NewMockRepository(ctrl)
				n := 2

				downloadBookChan := make(chan model.Book, n+1)
				defer close(downloadBookChan)
				rpo.EXPECT().FindBooksForDownload().Return(downloadBookChan, nil)

				for i := 0; i < n; i++ {
					downloadBookChan <- model.Book{
						ID:     i,
						Title:  fmt.Sprintf("title %v", i),
						Writer: model.Writer{Name: fmt.Sprintf("writer %v", i)},
						Status: model.End,
					}

					c.EXPECT().Get(gomock.Any(), fmt.Sprintf("http://test.com/book/%v/download", i)).Return(fmt.Sprintf("chapter list %v", i), nil)

					if i%2 == 1 {
						parsedChapterList := parse.ParsedChapterList{}
						parsedChapterList.Append(fmt.Sprintf("http://test.com/bk/%v/chapter/1", i), "test chapter 1")

						p.EXPECT().ParseChapterList(fmt.Sprintf("chapter list %v", i)).Return(&parsedChapterList, nil)

						c.EXPECT().Get(gomock.Any(), fmt.Sprintf("http://test.com/bk/%v/chapter/1", i)).Return("chapter 1", nil)

						p.EXPECT().ParseChapter("chapter 1").Return(
							parse.NewParsedChapterFields("title 1", "content 1"), nil,
						)

						rpo.EXPECT().UpdateBook(
							&model.Book{
								ID:           i,
								Title:        fmt.Sprintf("title %v", i),
								Writer:       model.Writer{Name: fmt.Sprintf("writer %v", i)},
								Status:       model.End,
								IsDownloaded: true,
							},
						).Return(nil)
					} else {
						p.EXPECT().ParseChapterList(fmt.Sprintf("chapter list %v", i)).Return(nil, parse.ErrParseChapterListNoChapterFound)
					}
				}

				return ServiceImp{
					client: c,
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(2),
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:                    config.URLConfig{Download: "http://test.com/book/%v/download"},
						Storage:                "./test-download/",
						MaxDownloadConcurrency: 1,
					},
				}
			},
			expectFiles: []expectFile{
				{
					expectFileLocation: "./test-download/1.txt",
					expectContent:      "title 1\nwriter 1\n--------------------\n\ntitle 1\n--------------------\ncontent 1\n--------------------\n",
				},
			},
			expectedError:  false,
			expectErrorStr: "",
		},
		{
			name: "fail to find book for download",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksForDownload().Return(nil, errors.New("some error"))

				return ServiceImp{
					rpo: rpo,
				}
			},
			expectFiles:    []expectFile{},
			expectedError:  true,
			expectErrorStr: "fail to fetch books: some error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			err := serv.Download()

			if test.expectedError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.expectErrorStr)
			}

			for _, expectFile := range test.expectFiles {
				if assert.FileExists(t, expectFile.expectFileLocation) {
					defer os.Remove(expectFile.expectFileLocation)
				}
				b, err := os.ReadFile(expectFile.expectFileLocation)
				if err != nil {
					t.Errorf("read file error: %v", err)
				}
				assert.Equal(t, expectFile.expectContent, string(b))
			}
		})
	}
}
