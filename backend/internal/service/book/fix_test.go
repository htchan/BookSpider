package book

import (
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

func Test_checkStorage(t *testing.T) {
	t.Parallel()
	os.Create("./1.txt")

	t.Cleanup(func() {
		os.Remove("./1.txt")
	})

	tests := []struct {
		name       string
		bk         model.Book
		stConf     config.SiteConfig
		expect     bool
		expectBook model.Book
	}{
		{
			name:       "return false for existing file and downloaded book",
			bk:         model.Book{ID: 1, HashCode: 0, Status: model.End, IsDownloaded: true},
			stConf:     config.SiteConfig{Storage: "."},
			expect:     false,
			expectBook: model.Book{ID: 1, HashCode: 0, Status: model.End, IsDownloaded: true},
		},
		{
			name:       "return true for existing file and non downloaded book",
			bk:         model.Book{ID: 1, HashCode: 0, Status: model.End, IsDownloaded: false},
			stConf:     config.SiteConfig{Storage: "."},
			expect:     true,
			expectBook: model.Book{ID: 1, HashCode: 0, Status: model.End, IsDownloaded: true},
		},
		{
			name:       "return true for non existing file and downloaded book",
			bk:         model.Book{ID: 2, HashCode: 0, Status: model.End, IsDownloaded: true},
			stConf:     config.SiteConfig{Storage: "."},
			expect:     true,
			expectBook: model.Book{ID: 2, HashCode: 0, Status: model.End, IsDownloaded: false},
		},
		{
			name:       "return false for non existing file and non downloaded book",
			bk:         model.Book{ID: 2, HashCode: 0, Status: model.End, IsDownloaded: false},
			stConf:     config.SiteConfig{Storage: "."},
			expect:     false,
			expectBook: model.Book{ID: 2, HashCode: 0, Status: model.End, IsDownloaded: false},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := checkStorage(&test.bk, &test.stConf)
			if result != test.expect {
				t.Errorf("got: %v; want: %v", result, test.expect)
			}
			if !cmp.Equal(test.bk, test.expectBook) {
				t.Errorf(cmp.Diff(test.bk, test.expectBook))
			}
		})
	}
}

func Test_Fix(t *testing.T) {
	t.Parallel()
}
