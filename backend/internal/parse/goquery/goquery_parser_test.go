package goquery

import (
	"testing"

	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/stretchr/testify/assert"
)

func Test_LoadParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		conf           *config.GoquerySelectorConfig
		expectedParser *GoqueryParser
		expectError    bool
	}{
		{
			name: "load parser successfully",
			conf: &config.GoquerySelectorConfig{
				Title:          "title",
				Writer:         "writer",
				BookType:       "book-type",
				LastUpdate:     "update-date",
				LastChapter:    "update-chapter",
				BookChapter:    "book-chapter",
				ChapterTitle:   "chapter-title",
				ChapterContent: "chapter-content",
			},
			expectedParser: &GoqueryParser{
				titleSelector:          "title",
				writerSelector:         "writer",
				bookTypeSelector:       "book-type",
				lastUpdateSelector:     "update-date",
				lastChapterSelector:    "update-chapter",
				bookChapterSelector:    "book-chapter",
				ChapterTitleSelector:   "chapter-title",
				ChapterContentSelector: "chapter-content",
			},
			expectError: false,
		},
		{
			name: "missing book info selector",
			conf: &config.GoquerySelectorConfig{
				BookChapter:    "book-chapter",
				ChapterTitle:   "chapter-title",
				ChapterContent: "chapter-content",
			},
			expectedParser: nil,
			expectError:    true,
		},
		{
			name: "missing book chapter selector",
			conf: &config.GoquerySelectorConfig{
				Title:          "title",
				Writer:         "writer",
				BookType:       "book-type",
				LastUpdate:     "update-date",
				LastChapter:    "update-chapter",
				ChapterTitle:   "chapter-title",
				ChapterContent: "chapter-content",
			},
			expectedParser: nil,
			expectError:    true,
		},
		{
			name: "missing chapter info selector",
			conf: &config.GoquerySelectorConfig{
				Title:       "title",
				Writer:      "writer",
				BookType:    "book-type",
				LastUpdate:  "update-date",
				LastChapter: "update-chapter",
				BookChapter: "book-chapter",
			},
			expectedParser: nil,
			expectError:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			parser, err := LoadParser(test.conf)

			if (err != nil) != test.expectError {
				t.Errorf("error fidd")
				t.Errorf("expect error exist: %v", test.expectError)
				t.Errorf("got error: %v", err)
			}

			assert.Equal(t, test.expectedParser, parser)
		})
	}
}

func TestParser_ParseBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		parser       GoqueryParser
		html         string
		expectFields *parse.ParsedBookFields
		expectError  bool
	}{
		{
			name: "success parse html to book info",
			parser: GoqueryParser{
				titleSelector:       "title",
				writerSelector:      "writer",
				bookTypeSelector:    "type",
				lastUpdateSelector:  "date",
				lastChapterSelector: "chapter",
			},
			html: `<html><head>
			<title>title</title>
			<writer>writer</writer>
			<type>type</type>
			<date>date</date>
			<chapter>chapter</chapter>
			</head></html>`,
			expectFields: parse.NewParsedBookFields(
				"title",
				"writer",
				"type",
				"date",
				"chapter",
			),
			expectError: false,
		},
		{
			name: "some book info in missing",
			parser: GoqueryParser{
				titleSelector:       "title",
				writerSelector:      "writer",
				bookTypeSelector:    "type",
				lastUpdateSelector:  "date",
				lastChapterSelector: "chapter",
			},
			html: `<html><head>
			<title>title</title>
			<type>type</type>
			<date>date</date>
			<chapter>chapter</chapter>
			</head></html>`,
			expectFields: nil,
			expectError:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fields, err := test.parser.ParseBook(test.html)
			if (err != nil) != test.expectError {
				t.Errorf("error diff")
				t.Errorf("expect error exist: %v", test.expectError)
				t.Errorf("got error: %v", err)
			}

			assert.Equal(t, test.expectFields, fields)
		})
	}
}

func TestParser_ParserChapterList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		parser            GoqueryParser
		html              string
		expectChapterList *parse.ParsedChapterList
		expectError       bool
	}{
		{
			name: "success parse html to book info",
			parser: GoqueryParser{
				bookChapterSelector: "ul>li",
			},
			html: `<html><body>
			<ul>
				<li href="link 1">chap 1</li>
				<li href="link 2">chap 2</li>
			</ul>
			</body></html>`,
			expectChapterList: func() *parse.ParsedChapterList {
				var fields parse.ParsedChapterList
				fields.Append("link 1", "chap 1")
				fields.Append("link 2", "chap 2")

				return &fields
			}(),
			expectError: false,
		},
		{
			name: "some chapter missing link / title",
			parser: GoqueryParser{
				bookChapterSelector: "ul>li",
			},
			html: `<html><body>
			<ul>
				<li href="link 1">chap 1</li>
				<li href="">chap 2</li>
				<li href="link 3"></li>
			</ul>
			</body></html>`,
			expectChapterList: func() *parse.ParsedChapterList {
				var fields parse.ParsedChapterList
				fields.Append("link 1", "chap 1")
				fields.Append("", "chap 2")
				fields.Append("link 3", "")

				return &fields
			}(),
			expectError: false,
		},
		{
			name: "no chapters",
			parser: GoqueryParser{
				bookChapterSelector: "ul>li",
			},
			html: `<html><body>
			</body></html>`,
			expectChapterList: nil,
			expectError:       true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			chapterList, err := test.parser.ParseChapterList(test.html)
			if (err != nil) != test.expectError {
				t.Errorf("error diff")
				t.Errorf("expect error exist: %v", test.expectError)
				t.Errorf("got error: %v", err)
			}

			assert.Equal(t, test.expectChapterList, chapterList)
		})
	}
}

func TestParser_parserChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		parser       GoqueryParser
		html         string
		expectFields *parse.ParsedChapterFields
		expectError  bool
	}{
		{
			name: "success parse html to book info",
			parser: GoqueryParser{
				ChapterTitleSelector:   "title",
				ChapterContentSelector: "content",
			},
			html: `<html><body>
			<title>title</title>
			<content>
				some long long long 
				long long long 
				long long long content
			</content>
			</body></html>`,
			expectFields: parse.NewParsedChapterFields(
				"title",
				"some long long long \n\t\t\t\tlong long long \n\t\t\t\tlong long long content\n",
			),
			expectError: false,
		},
		{
			name: "content not found",
			parser: GoqueryParser{
				ChapterTitleSelector:   "title",
				ChapterContentSelector: "content",
			},
			html: `<html><body>
			<title>title</title>
			</body></html>`,
			expectFields: nil,
			expectError:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fields, err := test.parser.ParseChapter(test.html)
			if (err != nil) != test.expectError {
				t.Errorf("error diff")
				t.Errorf("expect error exist: %v", test.expectError)
				t.Errorf("got error: %v", err)
			}

			assert.Equal(t, test.expectFields, fields)
		})
	}

}
