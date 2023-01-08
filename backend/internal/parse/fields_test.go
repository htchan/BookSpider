package parse

import (
	"testing"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestParsedBookFields_Populate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fields     ParsedBookFields
		bk         model.Book
		expectBook model.Book
	}{
		{
			name: "populate book fields to model.Book",
			fields: ParsedBookFields{
				title:         "title",
				writer:        "writer",
				bookType:      "book-type",
				updateDate:    "update-date",
				updateChapter: "update-chapter",
			},
			bk: model.Book{},
			expectBook: model.Book{
				Title:         "title",
				Writer:        model.Writer{Name: "writer"},
				Type:          "book-type",
				UpdateDate:    "update-date",
				UpdateChapter: "update-chapter",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			test.fields.Populate(&test.bk)

			assert.Equal(t, test.expectBook, test.bk)
		})
	}
}

func TestParsedChapterList_Populate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		fields            ParsedChapterList
		chapterList       model.Chapters
		expectChapterList model.Chapters
	}{
		{
			name: "populate book fields to model.Chapters",
			fields: ParsedChapterList{
				chapters: []struct {
					url   string
					title string
				}{
					{url: "url 1", title: "title 1"},
					{url: "url 2", title: "title 2"},
				},
			},
			chapterList: model.Chapters{},
			expectChapterList: model.Chapters{
				{Index: 0, URL: "url 1", Title: "title 1"},
				{Index: 1, URL: "url 2", Title: "title 2"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			test.fields.Populate(&test.chapterList)

			assert.Equal(t, test.expectChapterList, test.chapterList)
		})
	}
}

func TestParsedChapterFields_Populate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fields     ParsedChapterFields
		bk         model.Chapter
		expectBook model.Chapter
	}{
		{
			name: "populate book fields to model.Chapter",
			fields: ParsedChapterFields{
				title:   "title",
				content: "some content",
			},
			bk: model.Chapter{},
			expectBook: model.Chapter{
				Title:   "title",
				Content: "some content",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			test.fields.Populate(&test.bk)

			assert.Equal(t, test.expectBook, test.bk)
		})
	}
}
