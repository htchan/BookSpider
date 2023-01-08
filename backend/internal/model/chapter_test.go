package model

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
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

			assert.Equal(t, result, test.expect)
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

func TestChapter_optimizeContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chapter Chapter
		expect  Chapter
	}{
		{
			name:    "remove specific string",
			chapter: Chapter{Content: "&nbsp;<b></b></p>                "},
			expect:  Chapter{Content: ""},
		},
		{
			name:    "replace specific string to \\n",
			chapter: Chapter{Content: "<br /><p/>"},
			expect:  Chapter{Content: ""},
		},
		{
			name:    "remove space / tab in each line",
			chapter: Chapter{Content: " abc \n\tdef\t"},
			expect:  Chapter{Content: "abc\n\ndef"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			test.chapter.OptimizeContent()

			assert.Equal(t, test.chapter, test.expect)
		})
	}
}

func Test_removeEmptyLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		lines  []string
		expect []string
	}{
		{
			name:   "remove empty lines in input",
			lines:  []string{" \t", "\t abc\t ", "\t "},
			expect: []string{"abc"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := removeEmptyLines(test.lines)
			assert.Equal(t, result, test.expect)
		})
	}
}

func Test_StringToContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		content   string
		expect    Chapters
		expectErr bool
	}{
		{
			name: "works for new version",
			content: `
			title
			writer
			--------------------
			
			chapter title 1
			--------------------
			content1
			
			content2
			--------------------
			chapter title 2
			--------------------
			content3
			--------------------
			chapter title 3
			--------------------
			content4
			
			content5
			content6
			--------------------
			`,
			expect: Chapters{
				Chapter{Index: 0, Title: "chapter title 1", Content: "content1\n\ncontent2"},
				Chapter{Index: 1, Title: "chapter title 2", Content: "content3"},
				Chapter{Index: 2, Title: "chapter title 3", Content: "content4\n\ncontent5\n\ncontent6"},
			},
			expectErr: false,
		},
		{
			name: "works for old version",
			content: `
			title
			writer
			--------------------
			
			chapter title 1
			--------------------
			content1
			
			content2

			chapter title 2
			--------------------
			content3
			
			chapter title 3
			--------------------
			content4
			
			content5
			content6`,
			expect: Chapters{
				Chapter{Index: 0, Title: "chapter title 1", Content: "content1\n\ncontent2"},
				Chapter{Index: 1, Title: "chapter title 2", Content: "content3"},
				Chapter{Index: 2, Title: "chapter title 3", Content: "content4\n\ncontent5\n\ncontent6"},
			},
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := StringToChapters(test.content)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v, expect error: %v", err, test.expectErr)
			}
			assert.Equal(t, result, test.expect)
		})
	}
}
