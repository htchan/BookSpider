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

func Test_baseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     model.Book
		conf   config.BookConfig
		expect string
	}{
		{
			name:   "works",
			bk:     model.Book{ID: 123},
			conf:   config.BookConfig{URLConfig: config.URLConfig{Base: "abcdef/%v"}},
			expect: "abcdef/123",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := baseURL(test.bk, test.conf)
			if result != test.expect {
				t.Errorf(cmp.Diff(result, test.expect))
			}
		})
	}
}

func Test_fetchInfo(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()
	t.Cleanup(func() {
		server.Close()
	})

	c := client.NewClient(config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)

	type ExpectStruct struct {
		title, writer, typeStr, date, chapterStr string
		err                                      bool
	}

	tests := []struct {
		name         string
		client       *client.CircuitBreakerClient
		bookKey, url string
		bkConf       config.BookConfig
		expect       ExpectStruct
	}{
		{
			name:    "nothing change",
			client:  &c,
			bookKey: "test_book",
			url:     server.URL + "/no-update-book",
			expect: ExpectStruct{
				title: "title", writer: "writer", typeStr: "type",
				date: "1234", chapterStr: "chapter", err: false,
			},
		},
		{
			name:    "error",
			client:  &c,
			bookKey: "test_book",
			url:     server.URL + "/error",
			expect: ExpectStruct{
				title: "", writer: "", typeStr: "",
				date: "", chapterStr: "", err: true,
			},
		},
		{
			name:    "chapter changed to end",
			client:  &c,
			bookKey: "test_book",
			url:     server.URL + "/update-book/chapter-end",
			expect: ExpectStruct{
				title: "title", writer: "writer", typeStr: "type",
				date: "1234", chapterStr: "chapter後記", err: false,
			},
		},
		{
			name:    "chapter changed to not end",
			client:  &c,
			bookKey: "test_book",
			url:     server.URL + "/update-book/chapter-not-end",
			expect: ExpectStruct{
				title: "title", writer: "writer", typeStr: "type",
				date: "1234", chapterStr: "chapter-new", err: false,
			},
		},
		{
			name:    "title changed",
			client:  &c,
			bookKey: "test_book",
			url:     server.URL + "/update-book/title",
			expect: ExpectStruct{
				title: "title-new", writer: "writer", typeStr: "type",
				date: "1234", chapterStr: "chapter", err: false,
			},
		},
		{
			name:    "replace unwant content in response",
			client:  &c,
			bookKey: "test_book",
			url:     server.URL + "/update-book/chapter-not-end",
			bkConf:  config.BookConfig{UnwantContent: []string{"ne"}},
			expect: ExpectStruct{
				title: "title", writer: "writer", typeStr: "type",
				date: "1234", chapterStr: "chapter-w", err: false,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			title, writer, typeStr, date, chapterStr, err := fetchInfo(test.url, test.client, test.bookKey, test.bkConf)
			if (err != nil) != test.expect.err {
				t.Errorf("got error: %v; want: %v", err, test.expect.err)
			}
			if title != test.expect.title {
				t.Errorf("got title: %v; want: %v", title, test.expect.title)
			}
			if writer != test.expect.writer {
				t.Errorf("got writer: %v; want: %v", writer, test.expect.writer)
			}
			if typeStr != test.expect.typeStr {
				t.Errorf("got typeStr: %v; want: %v", typeStr, test.expect.typeStr)
			}
			if date != test.expect.date {
				t.Errorf("got date: %v; want: %v", date, test.expect.date)
			}
			if chapterStr != test.expect.chapterStr {
				t.Errorf("got chapterStr: %v; want: %v", chapterStr, test.expect.chapterStr)
			}
		})
	}
}

func Test_isNewBook(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		bk            model.Book
		title, writer string
		expect        bool
	}{
		{
			name:  "return true for different title",
			bk:    model.Book{Title: "title", Writer: model.Writer{Name: "writer"}, Status: model.InProgress},
			title: "title new", writer: "writer",
			expect: true,
		},
		{
			name:  "return true for different writer",
			bk:    model.Book{Title: "title", Writer: model.Writer{Name: "writer"}, Status: model.InProgress},
			title: "title", writer: "writer new",
			expect: true,
		},
		{
			name:  "return true for different title and writer",
			bk:    model.Book{Title: "title", Writer: model.Writer{Name: "writer"}, Status: model.InProgress},
			title: "title new", writer: "writer new",
			expect: true,
		},
		{
			name:  "return false for same title and writer",
			bk:    model.Book{Title: "title", Writer: model.Writer{Name: "writer"}, Status: model.InProgress},
			title: "title", writer: "writer",
			expect: false,
		},
		{
			name:  "return false for error book",
			bk:    model.Book{Title: "title", Writer: model.Writer{Name: "writer"}, Status: model.Error},
			title: "title new", writer: "writer new",
			expect: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := isNewBook(test.bk, test.title, test.writer)
			if result != test.expect {
				t.Errorf("got: %v; want: %v", result, test.expect)
			}
		})
	}
}

func Test_isUpdated(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                                     string
		bk                                       model.Book
		title, writer, typeStr, date, chapterStr string
		expect                                   bool
	}{
		{
			name: "return true for different title",
			bk: model.Book{
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			title: "title new", writer: "writer", typeStr: "type", date: "date", chapterStr: "chapter",
			expect: true,
		},
		{
			name: "return true for different writer",
			bk: model.Book{
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			title: "title", writer: "writer new", typeStr: "type", date: "date", chapterStr: "chapter",
			expect: true,
		},
		{
			name: "return true for different type",
			bk: model.Book{
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			title: "title", writer: "writer", typeStr: "type new", date: "date", chapterStr: "chapter",
			expect: true,
		},
		{
			name: "return true for different update date",
			bk: model.Book{
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			title: "title", writer: "writer", typeStr: "type", date: "date new", chapterStr: "chapter",
			expect: true,
		},
		{
			name: "return true for different chapter",
			bk: model.Book{
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			title: "title", writer: "writer", typeStr: "type", date: "date", chapterStr: "chapter new",
			expect: true,
		},
		{
			name: "return false for all parameter same",
			bk: model.Book{
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			title: "title", writer: "writer", typeStr: "type", date: "date", chapterStr: "chapter",
			expect: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := isUpdated(
				test.bk, test.title, test.writer, test.typeStr, test.date, test.chapterStr,
			)
			if result != test.expect {
				t.Errorf("got: %v; want: %v", result, test.expect)
			}
		})
	}
}

func Test_Update(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()
	t.Cleanup(func() {
		server.Close()
	})

	c := client.NewClient(config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)

	tests := []struct {
		name       string
		bk         model.Book
		bkConf     config.BookConfig
		stConf     config.SiteConfig
		c          *client.CircuitBreakerClient
		expect     bool
		expectErr  bool
		expectBook model.Book
	}{
		{
			name: "not updated",
			bk: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
			},
			bkConf:    config.BookConfig{URLConfig: config.URLConfig{Base: server.URL + "/no-update-book/%v"}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    false,
			expectErr: false,
			expectBook: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
			},
		},
		{
			name: "error",
			bk: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
			},
			bkConf:    config.BookConfig{URLConfig: config.URLConfig{Base: server.URL + "/error/%v"}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    false,
			expectErr: true,
			expectBook: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter", Error: errors.New("title not found"),
			},
		},
		{
			name: "updated chapter",
			bk: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
			},
			bkConf:    config.BookConfig{URLConfig: config.URLConfig{Base: server.URL + "/update-book/chapter-end/%v"}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter後記", Status: model.InProgress,
			},
		},
		{
			name: "updated title",
			bk: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter", Status: model.InProgress,
			},
			bkConf:    config.BookConfig{URLConfig: config.URLConfig{Base: server.URL + "/update-book/title/%v"}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				Site: "test", ID: 1, HashCode: model.GenerateHash(),
				Title: "title-new", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter", Status: model.InProgress,
			},
		},
		{
			name: "updated error status to in progress",
			bk: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Status: model.Error, Error: errors.New("error"),
			},
			bkConf:    config.BookConfig{URLConfig: config.URLConfig{Base: server.URL + "/no-update-book/%v"}},
			stConf:    config.SiteConfig{BookKey: "test_book"},
			c:         &c,
			expect:    true,
			expectErr: false,
			expectBook: model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "1234", UpdateChapter: "chapter",
				Status: model.InProgress, Error: nil,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := Update(&test.bk, test.bkConf, test.stConf, test.c)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v, want: %v", err, test.expectErr)
			}
			if result != test.expect {
				t.Errorf("got: %v; want: %v", result, test.expect)
			}
			if !cmp.Equal(test.bk, test.expectBook) {
				t.Errorf("book diff: %v", cmp.Diff(test.bk, test.expectBook))
			}
		})
	}
}
