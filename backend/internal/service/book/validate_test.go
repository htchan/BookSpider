package book

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/model"
)

func Test_isEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chapter string
		expect  bool
	}{
		{
			name:    "chapter not end",
			chapter: "still in progress",
			expect:  false,
		},
		{
			name:    "chapter already end",
			chapter: "chapter - 全文結",
			expect:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := isEnd(test.chapter)
			if result != test.expect {
				t.Errorf("got: %v; want: %v", result, test.expect)
			}
		})
	}
}

func Test_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		bk         model.Book
		expect     bool
		expectErr  bool
		expectBook model.Book
	}{
		{
			name:       "validate does not change in progress book with is downloaded = false",
			bk:         model.Book{UpdateChapter: "chapter", Status: model.InProgress, IsDownloaded: false},
			expect:     false,
			expectErr:  false,
			expectBook: model.Book{UpdateChapter: "chapter", Status: model.InProgress, IsDownloaded: false},
		},
		{
			name:       "validate does not change in progress book with is downloaded = true",
			bk:         model.Book{UpdateChapter: "chapter", Status: model.InProgress, IsDownloaded: true},
			expect:     false,
			expectErr:  false,
			expectBook: model.Book{UpdateChapter: "chapter", Status: model.InProgress, IsDownloaded: true},
		},
		{
			name:       "validate does not end book with is_downloaded = true",
			bk:         model.Book{UpdateChapter: "chapter - 全文結", Status: model.End, IsDownloaded: true},
			expect:     false,
			expectErr:  false,
			expectBook: model.Book{UpdateChapter: "chapter - 全文結", Status: model.End, IsDownloaded: true},
		},
		{
			name:       "validate does not end book with is_downloaded = false",
			bk:         model.Book{UpdateChapter: "chapter - 全文結", Status: model.End, IsDownloaded: false},
			expect:     false,
			expectErr:  false,
			expectBook: model.Book{UpdateChapter: "chapter - 全文結", Status: model.End, IsDownloaded: false},
		},
		{
			name:       "validate set book to end and is_downloaded to false if chapter is end",
			bk:         model.Book{UpdateChapter: "chapter - 全文結", Status: model.InProgress, IsDownloaded: true},
			expect:     true,
			expectErr:  false,
			expectBook: model.Book{UpdateChapter: "chapter - 全文結", Status: model.End, IsDownloaded: false},
		},
		{
			name:       "validate set book status to in progress if chapter is not end",
			bk:         model.Book{UpdateChapter: "chapter", Status: model.End, IsDownloaded: true},
			expect:     true,
			expectErr:  false,
			expectBook: model.Book{UpdateChapter: "chapter", Status: model.InProgress, IsDownloaded: true},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result, err := Validate(&test.bk)
			if (err != nil) != test.expectErr {
				t.Errorf("got: %v; want: %v", err, test.expectErr)
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
