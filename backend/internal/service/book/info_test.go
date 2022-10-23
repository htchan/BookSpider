package book

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
)

func Test_BookFileLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     model.Book
		stConf config.SiteConfig
		expect string
	}{
		{
			name:   "works with hashcode = 0",
			bk:     model.Book{Site: "test", ID: 1, HashCode: 0},
			stConf: config.SiteConfig{Storage: "/test"},
			expect: "/test/1.txt",
		},
		{
			name:   "works with hashcode = 1-0",
			bk:     model.Book{Site: "test", ID: 1, HashCode: 100},
			stConf: config.SiteConfig{Storage: "/test"},
			expect: "/test/1-v2s.txt",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := BookFileLocation(&test.bk, &test.stConf)

			if result != test.expect {
				t.Errorf("got: %v; want: %v", result, test.expect)
			}
		})
	}
}

func Test_Info(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     model.Book
		expect string
	}{
		{
			name: "works",
			bk: model.Book{
				Site: "site", ID: 1, HashCode: 30,
				Title: "title", Writer: model.Writer{ID: 1, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.End, IsDownloaded: false, Error: errors.New("data"),
			},
			expect: `{"site":"site","id":1,"hash_code":"u","title":"title","writer":"writer","type":"type","update_date":"date","update_chapter":"chapter","status":"END","is_downloaded":false,"error":"data"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := Info(&test.bk)
			if result != test.expect {
				t.Error(result, test.expect)
			}
		})
	}
}

func Test_Content(t *testing.T) {
	t.Parallel()

	os.Mkdir("./storage", os.ModePerm)
	os.WriteFile("./storage/1.txt", []byte("hello"), os.ModePerm)

	t.Cleanup(func() {
		os.RemoveAll("./storage")
	})

	tests := []struct {
		name      string
		stConf    *config.SiteConfig
		bk        *model.Book
		expect    string
		expectErr bool
	}{
		{
			name:      "return existing file of download book",
			stConf:    &config.SiteConfig{Storage: "./storage"},
			bk:        &model.Book{ID: 1, IsDownloaded: true},
			expect:    "hello",
			expectErr: false,
		},
		{
			name:      "return error for not existing book",
			stConf:    &config.SiteConfig{Storage: "./storage"},
			bk:        &model.Book{ID: 2, IsDownloaded: false},
			expect:    "",
			expectErr: true,
		},
		{
			name:      "return error for not downloaded book",
			stConf:    &config.SiteConfig{Storage: "./storage"},
			bk:        &model.Book{ID: 1},
			expect:    "",
			expectErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := Content(test.bk, test.stConf)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v, want err: %v", err, test.expectErr)
			}

			if result != test.expect {
				t.Errorf("content diff: %v", cmp.Diff(result, test.expect))
			}
		})
	}
}
