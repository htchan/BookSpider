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

func Test_downloadURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     model.Book
		bkConf config.BookConfig
		expect string
	}{
		{
			name:   "works",
			bk:     model.Book{ID: 1},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{Download: "http://test.com/download/%v"}},
			expect: "http://test.com/download/1",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := downloadURL(&test.bk, test.bkConf)
			if result != test.expect {
				t.Errorf(cmp.Diff(result, test.expect))
			}
		})
	}
}

func Test_fetchChaptersHeaderInfo(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()
	t.Cleanup(func() {
		server.Close()
	})

	c := client.NewClient(config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)

	tests := []struct {
		name           string
		bk             model.Book
		bkConf         config.BookConfig
		stConf         config.SiteConfig
		c              *client.CircuitBreakerClient
		expectChapters model.Chapters
		expectErr      bool
	}{
		{
			name:   "valid chapters",
			bk:     model.Book{ID: 1},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{Download: server.URL + "/chapter-header/valid/%v"}},
			stConf: config.SiteConfig{BookKey: "test_book"},
			c:      &c,
			expectChapters: model.Chapters{
				{Index: 0, URL: "/1", Title: "1"},
				{Index: 1, URL: "/2", Title: "2"},
				{Index: 2, URL: "/3", Title: "3"},
				{Index: 3, URL: "/4", Title: "4"},
			},
			expectErr: false,
		},
		{
			name:   "imbalance chapters",
			bk:     model.Book{ID: 1},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{Download: server.URL + "/chapter-header/imbalance-url-title/%v"}},
			stConf: config.SiteConfig{BookKey: "test_book"},
			c:      &c,
			expectChapters: model.Chapters{
				{Index: 0, URL: "/1", Title: "1"},
				{Index: 1, URL: "/2", Title: "2"},
				{Index: 2, URL: "/3", Title: "3"},
			},
			expectErr: false,
		},
		{
			name:           "empty chapters",
			bk:             model.Book{ID: 1},
			bkConf:         config.BookConfig{URLConfig: config.URLConfig{Download: server.URL + "/chapter-header/empty/%v"}},
			stConf:         config.SiteConfig{BookKey: "test_book"},
			c:              &c,
			expectChapters: nil,
			expectErr:      true,
		},
		{
			name:           "unrecognize chapters",
			bk:             model.Book{ID: 1},
			bkConf:         config.BookConfig{URLConfig: config.URLConfig{Download: server.URL + "/chapter-header/url-not-found/%v"}},
			stConf:         config.SiteConfig{BookKey: "test_book"},
			c:              &c,
			expectChapters: nil,
			expectErr:      true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := fetchChaptersHeaderInfo(&test.bk, test.bkConf, test.stConf, test.c)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if !cmp.Equal(result, test.expectChapters) {
				t.Error(cmp.Diff(result, test.expectChapters))
			}
		})
	}
}

func Test_downloadChapters(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()
	t.Cleanup(func() {
		server.Close()
	})

	c := client.NewClient(config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)

	tests := []struct {
		name           string
		bk             model.Book
		chapters       model.Chapters
		bkConf         config.BookConfig
		stConf         config.SiteConfig
		c              *client.CircuitBreakerClient
		expectChapters model.Chapters
		expectErr      bool
	}{
		{
			name: "works",
			bk:   model.Book{ID: 1},
			chapters: model.Chapters{
				{URL: server.URL + "/chapter/extra-content", Title: "title 1"},
				{URL: server.URL + "/chapter/extra-content", Title: "title 2"},
			},
			bkConf: config.BookConfig{},
			stConf: config.SiteConfig{BookKey: "test_book"},
			c:      &c,
			expectChapters: model.Chapters{
				{URL: server.URL + "/chapter/extra-content", Title: "title 1", Content: "success-\ncontent-regex"},
				{URL: server.URL + "/chapter/extra-content", Title: "title 2", Content: "success-\ncontent-regex"},
			},
			expectErr: false,
		},
		{
			name: "invalid chapter",
			bk:   model.Book{ID: 1},
			chapters: model.Chapters{
				{URL: server.URL + "/chapter/extra-content", Title: "title 1"},
				{URL: server.URL + "/chapter/unrecognize", Title: "title 2"},
			},
			bkConf: config.BookConfig{},
			stConf: config.SiteConfig{BookKey: "test_book"},
			c:      &c,
			expectChapters: model.Chapters{
				{URL: server.URL + "/chapter/extra-content", Title: "title 1", Content: "success-\ncontent-regex"},
				{URL: server.URL + "/chapter/unrecognize", Title: "title 2", Error: errors.New("chapter content not found")},
			},
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := downloadChapters(&test.bk, test.chapters, test.bkConf, test.stConf, test.c)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want err: %v", err, test.expectErr)
			}

			if !cmp.Equal(test.chapters, test.expectChapters) {
				t.Error(cmp.Diff(test.chapters, test.expectChapters))
			}
		})
	}
}

func Test_headerInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     model.Book
		expect string
	}{
		{
			name:   "works",
			bk:     model.Book{Title: "title", Writer: model.Writer{Name: "name"}},
			expect: "title\nname\n" + model.CONTENT_SEP + "\n\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := headerInfo(&test.bk)
			if result != test.expect {
				t.Error(cmp.Diff(result, test.expect))
			}
		})
	}
}

func Test_saveContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		location           string
		bk                 model.Book
		chapters           model.Chapters
		expectContentExist bool
		expectContent      string
		expectErr          bool
	}{
		{
			name:     "works",
			location: "./content-1.txt",
			bk:       model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "name"}},
			chapters: model.Chapters{
				{Title: "title 1", Content: "content 1"},
				{Title: "title 2", Content: "content 2"},
			},
			expectContentExist: true,
			expectContent: "title\nname\n" + model.CONTENT_SEP + "\n\n" +
				"title 1\n" + model.CONTENT_SEP + "\ncontent 1\n" + model.CONTENT_SEP + "\n" +
				"title 2\n" + model.CONTENT_SEP + "\ncontent 2\n" + model.CONTENT_SEP + "\n",
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			t.Cleanup(func() {
				os.Remove(test.location)
			})

			err := saveContent(test.location, &test.bk, test.chapters)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if err != nil {
				return
			}

			content, err := os.ReadFile(test.location)
			if (err != nil) != !test.expectContentExist {
				t.Errorf("fail to open file: %v", test.location)
			}

			if test.expectContentExist && string(content) != test.expectContent {
				t.Error(cmp.Diff(string(content), test.expectContent))
			}
		})
	}
}

func Test_Download(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()

	c := client.NewClient(config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)
	os.Mkdir("download_test", os.ModePerm)
	t.Cleanup(func() {
		server.Close()
		os.RemoveAll("download_test")
	})

	tests := []struct {
		name               string
		bk                 model.Book
		bkConf             config.BookConfig
		stConf             config.SiteConfig
		c                  *client.CircuitBreakerClient
		expect             bool
		expectErr          bool
		expectLocation     string
		expectContentExist bool
		expectContent      string
	}{
		{
			name: "return true for not downloaded book",
			bk:   model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "name"}, Status: model.End},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Download: server.URL + "/chapter-header/valid/%v", ChapterPrefix: server.URL + "/chapter/valid",
			}},
			stConf:             config.SiteConfig{BookKey: "test_book", Storage: "./download_test"},
			c:                  &c,
			expect:             true,
			expectErr:          false,
			expectLocation:     "./download_test/1.txt",
			expectContentExist: true,
			expectContent: "title\nname\n" + model.CONTENT_SEP + "\n" +
				"\n1\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n2\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n3\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n4\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP + "\n",
		},
		{
			name:               "return false for downloaded book",
			bk:                 model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "name"}, Status: model.End, IsDownloaded: true},
			bkConf:             config.BookConfig{},
			stConf:             config.SiteConfig{},
			c:                  &c,
			expect:             false,
			expectErr:          false,
			expectContentExist: false,
		},
		{
			name:               "return false for In Progress book",
			bk:                 model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "name"}, Status: model.InProgress},
			bkConf:             config.BookConfig{},
			stConf:             config.SiteConfig{},
			c:                  &c,
			expect:             false,
			expectErr:          false,
			expectContentExist: false,
		},
		{
			name:               "return false for Error book",
			bk:                 model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "name"}, Status: model.Error},
			bkConf:             config.BookConfig{},
			stConf:             config.SiteConfig{},
			c:                  &c,
			expect:             false,
			expectErr:          false,
			expectContentExist: false,
		},
		{
			name: "return false if it fail to load chapter header info",
			bk:   model.Book{ID: 1, Title: "title", Writer: model.Writer{Name: "name"}, Status: model.End},
			bkConf: config.BookConfig{URLConfig: config.URLConfig{
				Download: server.URL + "/chapter-header/empty/%v", ChapterPrefix: server.URL + "/chapter/valid",
			}},
			stConf:             config.SiteConfig{BookKey: "test_book", Storage: "./download_test"},
			c:                  &c,
			expect:             false,
			expectErr:          true,
			expectContentExist: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := Download(&test.bk, test.bkConf, test.stConf, test.c)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if result != test.expect {
				t.Errorf("got: %v, want: %v", result, test.expect)
			}

			content, err := os.ReadFile(test.expectLocation)
			if (err != nil) != !test.expectContentExist {
				t.Errorf("fail to open file: %v", test.expectLocation)
			}

			if test.expectContentExist && string(content) != test.expectContent {
				t.Error(cmp.Diff(string(content), test.expectContent))
			}
		})
	}
}
