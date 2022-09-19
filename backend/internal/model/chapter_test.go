package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		index  int
		url    string
		title  string
		expect Chapter
	}{
		{
			name:   "works",
			index:  10,
			url:    "url",
			title:  "title",
			expect: Chapter{Index: 10, URL: "url", Title: "title"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := NewChapter(test.index, test.url, test.title)

			if !cmp.Equal(result, test.expect) {
				t.Error(cmp.Diff(result, test.expect))
			}
		})
	}
}

func TestChapter_ContentString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		chapter Chapter
		expect  string
	}{
		{
			name:    "correct format",
			chapter: Chapter{Index: 1, URL: "url", Title: "title", Content: "content"},
			expect:  "title\n" + CONTENT_SEP + "\ncontent\n" + CONTENT_SEP + "\n",
		},
		{
			name:    "empty content",
			chapter: Chapter{Index: 1, URL: "url", Title: "title", Content: ""},
			expect:  "title\n" + CONTENT_SEP + "\n\n" + CONTENT_SEP + "\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.chapter.ContentString()
			if result != test.expect {
				t.Errorf(cmp.Diff(result, test.expect))
			}
		})
	}
}
