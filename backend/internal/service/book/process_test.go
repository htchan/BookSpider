package book

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/model"
)

func Test_Process(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()

	c := client.NewClient(config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)

	os.Mkdir("process_test", os.ModePerm)
	os.Create("process_test/12.txt")
	os.Create("process_test/10-v1.txt")
	os.Create("process_test/13.txt")
	t.Cleanup(func() {
		server.Close()
		os.RemoveAll("process_test")
	})

	tests := []struct {
		name               string
		bk                 model.Book
		bkConf             config.BookConfig
		stConf             config.SiteConfig
		c                  *client.CircuitBreakerClient
		expect             bool
		expectErr          bool
		expectBook         model.Book
		expectContentExist bool
		expectLocation     string
		expectContent      string
	}{
		{
			name: "error book has no update",
			bk:   model.Book{ID: 1, Status: model.Error},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/empty/%v",
			}},
			stConf:             config.SiteConfig{BookKey: "test_book"},
			c:                  &c,
			expect:             false,
			expectErr:          true,
			expectBook:         model.Book{ID: 1, Status: model.Error, Error: errors.New("zero length")},
			expectContentExist: false,
		},
		{
			name: "error book has updates",
			bk:   model.Book{ID: 2, Status: model.Error},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/no-update-book/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 2, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "in progress book has updated with new title",
			bk: model.Book{
				ID: 3, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/title/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 3, HashCode: model.GenerateHash(),
				Title: "title-new", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "in progress book has updated with non end chapter",
			bk: model.Book{
				ID: 4, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/chapter-not-end/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 4, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter-new",
				Status: model.InProgress, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "in progress book has updated with end chapter and fail to download",
			bk: model.Book{
				ID: 5, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base:     server.URL + "/update-book/chapter-end/%v",
				Download: server.URL + "/chapter-header/empty/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: true,
			expectBook: model.Book{
				ID: 5, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "in progress book has updated with end chapter and success to download",
			bk: model.Book{
				ID: 6, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base:          server.URL + "/update-book/chapter-end/%v",
				Download:      server.URL + "/chapter-header/valid/%v",
				ChapterPrefix: server.URL + "/chapter/valid",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book", Storage: "process_test"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 6, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: true, Error: nil,
			},
			expectContentExist: true,
			expectLocation:     "process_test/6.txt",
			expectContent: "title\nwriter\n" + model.CONTENT_SEP + "\n" +
				"\n1\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n2\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n3\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n4\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP + "\n",
		},
		{
			name: "end book that has chapter update",
			bk: model.Book{
				ID: 7, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.End, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/chapter-not-end/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 7, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter-new",
				Status: model.InProgress, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "end book that has title update",
			bk: model.Book{
				ID: 8, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.End, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/title/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 8, HashCode: model.GenerateHash(),
				Title: "title-new", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "end book that fail to download",
			bk: model.Book{
				ID: 9, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: false, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base:     server.URL + "/update-book/chapter-end/%v",
				Download: server.URL + "/chapter-header/empty/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    false,
			expectErr: true,
			expectBook: model.Book{
				ID: 9, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: false, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "end book that success to download",
			bk: model.Book{
				ID: 10, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base:          server.URL + "/update-book/chapter-end/%v",
				Download:      server.URL + "/chapter-header/valid/%v",
				ChapterPrefix: server.URL + "/chapter/valid",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book", Storage: "process_test"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 10, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: true, Error: nil,
			},
			expectContentExist: true,
			expectLocation:     "process_test/10.txt",
			expectContent: "title\nwriter\n" + model.CONTENT_SEP + "\n" +
				"\n1\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n2\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n3\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n4\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP + "\n",
		},
		{
			name: "end book has no update but an existing file",
			bk: model.Book{
				ID: 10, HashCode: 1,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: false, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/chapter-end/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book", Storage: "process_test"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 10, HashCode: 1,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: true, Error: nil,
			},
			expectContentExist: true,
			expectLocation:     "process_test/10-v1.txt",
			expectContent:      "",
		},
		{
			name: "downloaded book without file has no update and download fail",
			bk: model.Book{
				ID: 11, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: true, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base:     server.URL + "/update-book/chapter-end/%v",
				Download: server.URL + "/chapter-header/empty/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book", Storage: "process_test"},
			c:         &c,
			expect:    true,
			expectErr: true,
			expectBook: model.Book{
				ID: 11, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: false, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "downloaded book with file has no update",
			bk: model.Book{
				ID: 12, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: true, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/chapter-end/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book", Storage: "process_test"},
			c:         &c,
			expect:    false,
			expectErr: false,
			expectBook: model.Book{
				ID: 12, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記",
				Status: model.End, IsDownloaded: true, Error: nil,
			},
			expectContentExist: false,
		},
		{
			name: "downloaded book that has update",
			bk: model.Book{
				ID: 13, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.End, IsDownloaded: true, Error: nil,
			},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Base: server.URL + "/update-book/chapter-not-end/%v",
			}},
			stConf:    config.SiteConfig{BookKey: "test_book", Storage: "process_test"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				ID: 13, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter-new",
				Status: model.InProgress, IsDownloaded: true, Error: nil,
			},
			expectContentExist: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := Process(&test.bk, &test.bkConf, &test.stConf, test.c)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if result != test.expect {
				t.Errorf("got result: %v; want result: %v", result, test.expect)
			}

			if !cmp.Equal(test.bk, test.expectBook) {
				t.Errorf("book different: %v", cmp.Diff(test.bk, test.expectBook))
			}

			if !test.expectContentExist {
				return
			}

			content, err := os.ReadFile(test.expectLocation)
			if err != nil {
				t.Errorf("fail to read file: %v; error: %v", test.expectLocation, err)
			}
			if string(content) != test.expectContent {
				t.Errorf("content diff: %v", cmp.Diff(string(content), test.expectContent))
			}
		})
	}
}
