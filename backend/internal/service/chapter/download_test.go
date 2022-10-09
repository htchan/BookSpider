package chapter

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

func Test_chapterURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		bkID    int
		chapter *model.Chapter
		bkConf  config.BookConfig
		expect  string
	}{
		{
			name:    "works for chapter url start with http",
			bkID:    0,
			chapter: &model.Chapter{URL: "http://test.com/abc/def"},
			bkConf:  config.BookConfig{},
			expect:  "http://test.com/abc/def",
		},
		{
			name:    "works for chapter url start with /",
			bkID:    0,
			chapter: &model.Chapter{URL: "/abc/def"},
			bkConf:  config.BookConfig{URLConfig: config.URLConfig{ChapterPrefix: "http://test.com"}},
			expect:  "http://test.com/abc/def",
		},
		{
			name:    "works for chapter url not start with http or /, download url end with /",
			bkID:    15,
			chapter: &model.Chapter{URL: "abc/def"},
			bkConf:  config.BookConfig{URLConfig: config.URLConfig{Download: "http://test.com/%v/"}},
			expect:  "http://test.com/15/abc/def",
		},
		{
			name:    "works for chapter url not start with http or /, download url not end wih /",
			bkID:    15,
			chapter: &model.Chapter{URL: "abc/def"},
			bkConf:  config.BookConfig{URLConfig: config.URLConfig{Download: "http://test.com/%v"}},
			expect:  "http://test.com/15/abc/def",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := chapterURL(test.bkID, test.chapter, &test.bkConf)

			if result != test.expect {
				t.Errorf("got: %v; want: %v", result, test.expect)
			}
		})
	}
}

func Test_optimizeContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chapter model.Chapter
		expect  model.Chapter
	}{
		{
			name:    "remove specific string",
			chapter: model.Chapter{Content: "&nbsp;<b></b></p>                "},
			expect:  model.Chapter{Content: ""},
		},
		{
			name:    "replace specific string to \\n",
			chapter: model.Chapter{Content: "<br /><p/>"},
			expect:  model.Chapter{Content: "\n\n"},
		},
		{
			name:    "remove space / tab in each line",
			chapter: model.Chapter{Content: " abc \n\tdef\t"},
			expect:  model.Chapter{Content: "abc\ndef"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			optimizeContent(&test.chapter)

			if !cmp.Equal(test.chapter, test.expect) {
				t.Error(cmp.Diff(test.chapter, test.expect))
			}
		})
	}
}

func Test_Download(t *testing.T) {
	t.Parallel()
	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()
	t.Cleanup(func() {
		server.Close()
	})

	c := client.NewClient(config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)

	tests := []struct {
		name          string
		bkID          int
		chapter       model.Chapter
		bkConf        config.BookConfig
		stConf        config.SiteConfig
		c             *client.CircuitBreakerClient
		expectChapter model.Chapter
		expectErr     bool
	}{
		{
			name:          "download valid chapter",
			bkID:          1,
			chapter:       model.Chapter{URL: server.URL + "/chapter/valid"},
			bkConf:        config.BookConfig{},
			stConf:        config.SiteConfig{BookKey: "test_book"},
			c:             &c,
			expectChapter: model.Chapter{URL: server.URL + "/chapter/valid", Content: "success-content-regex"},
			expectErr:     false,
		},
		{
			name:          "download chapter with extra content",
			bkID:          1,
			chapter:       model.Chapter{URL: server.URL + "/chapter/extra-content"},
			bkConf:        config.BookConfig{},
			stConf:        config.SiteConfig{BookKey: "test_book"},
			c:             &c,
			expectChapter: model.Chapter{URL: server.URL + "/chapter/extra-content", Content: "success-\ncontent-regex"},
			expectErr:     false,
		},
		{
			name:          "download unrecognize chapter",
			bkID:          1,
			chapter:       model.Chapter{URL: server.URL + "/chapter/unrecognize"},
			bkConf:        config.BookConfig{},
			stConf:        config.SiteConfig{BookKey: "test_book"},
			c:             &c,
			expectChapter: model.Chapter{URL: server.URL + "/chapter/unrecognize", Error: errors.New("chapter content not found")},
			expectErr:     true,
		},
		{
			name:          "download empty chapter",
			bkID:          1,
			chapter:       model.Chapter{URL: server.URL + "/chapter/empty"},
			bkConf:        config.BookConfig{},
			stConf:        config.SiteConfig{BookKey: "test_book"},
			c:             &c,
			expectChapter: model.Chapter{URL: server.URL + "/chapter/empty", Error: errors.New("chapter content not found")},
			expectErr:     true,
		},
		{
			name:          "replace unwant content in response",
			bkID:          1,
			chapter:       model.Chapter{URL: server.URL + "/chapter/valid"},
			bkConf:        config.BookConfig{UnwantContent: []string{"regex"}},
			stConf:        config.SiteConfig{BookKey: "test_book"},
			c:             &c,
			expectChapter: model.Chapter{URL: server.URL + "/chapter/valid", Content: "success-content-"},
			expectErr:     false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := Download(test.bkID, &test.chapter, &test.bkConf, &test.stConf, test.c)

			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want err: %v", err, test.expectErr)
			}

			if !cmp.Equal(test.chapter, test.expectChapter) {
				t.Error(cmp.Diff(test.chapter, test.expectChapter))
			}
		})
	}
}
