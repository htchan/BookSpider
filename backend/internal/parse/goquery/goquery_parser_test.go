package goquery

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/stretchr/testify/assert"
)

func TestSeletor_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		selector Selector
		html     string
		want     string
	}{
		{
			name:     "happy flow with empty attr (return text)",
			selector: Selector{selector: "a", attr: ""},
			html:     `<html><a>test</a></html>`,
			want:     "test",
		},
		{
			name:     "happy flow with existing attr",
			selector: Selector{selector: "a", attr: "href"},
			html:     `<html><a href="abc">test</a></html>`,
			want:     "abc",
		},
		{
			name:     "not existing attr",
			selector: Selector{selector: "a", attr: "href"},
			html:     `<html><a class="abc">test</a></html>`,
			want:     "",
		},
		{
			name:     "selection not find",
			selector: Selector{selector: "li", attr: "href"},
			html:     `<html><a href="abc">test</a></html>`,
			want:     "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(test.html))
			if err != nil {
				t.Errorf("parse goquery.selection fail: %v", err)
			}

			selection := doc.Find(test.selector.selector)
			got := test.selector.Parse(selection)

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_LoadParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		conf           *config.GoquerySelectorsConfig
		expectedParser *GoqueryParser
		expectError    bool
	}{
		{
			name: "load parser successfully",
			conf: &config.GoquerySelectorsConfig{
				Title:            config.GoquerySelectorConfig{Selector: "title", Attr: ""},
				Writer:           config.GoquerySelectorConfig{Selector: "writer", Attr: ""},
				BookType:         config.GoquerySelectorConfig{Selector: "book-type", Attr: ""},
				LastUpdate:       config.GoquerySelectorConfig{Selector: "update-date", Attr: ""},
				LastChapter:      config.GoquerySelectorConfig{Selector: "update-chapter", Attr: ""},
				BookChapterURL:   config.GoquerySelectorConfig{Selector: "book-chapter", Attr: "href"},
				BookChapterTitle: config.GoquerySelectorConfig{Selector: "book-chapter", Attr: ""},
				ChapterTitle:     config.GoquerySelectorConfig{Selector: "chapter-title", Attr: ""},
				ChapterContent:   config.GoquerySelectorConfig{Selector: "chapter-content", Attr: ""},
			},
			expectedParser: &GoqueryParser{
				titleSelector:            Selector{"title", ""},
				writerSelector:           Selector{"writer", ""},
				bookTypeSelector:         Selector{"book-type", ""},
				lastUpdateSelector:       Selector{"update-date", ""},
				lastChapterSelector:      Selector{"update-chapter", ""},
				bookChapterURLSelector:   Selector{"book-chapter", "href"},
				bookChapterTitleSelector: Selector{"book-chapter", ""},
				ChapterTitleSelector:     Selector{"chapter-title", ""},
				ChapterContentSelector:   Selector{"chapter-content", ""},
			},
			expectError: false,
		},
		{
			name: "missing book info selector",
			conf: &config.GoquerySelectorsConfig{
				BookChapterURL:   config.GoquerySelectorConfig{Selector: "book-chapter", Attr: ""},
				BookChapterTitle: config.GoquerySelectorConfig{Selector: "book-chapter", Attr: ""},
				ChapterTitle:     config.GoquerySelectorConfig{Selector: "chapter-title", Attr: ""},
				ChapterContent:   config.GoquerySelectorConfig{Selector: "chapter-content", Attr: ""},
			},
			expectedParser: nil,
			expectError:    true,
		},
		{
			name: "missing book chapter selector",
			conf: &config.GoquerySelectorsConfig{
				Title:          config.GoquerySelectorConfig{Selector: "title", Attr: ""},
				Writer:         config.GoquerySelectorConfig{Selector: "writer", Attr: ""},
				BookType:       config.GoquerySelectorConfig{Selector: "book-type", Attr: ""},
				LastUpdate:     config.GoquerySelectorConfig{Selector: "update-date", Attr: ""},
				LastChapter:    config.GoquerySelectorConfig{Selector: "update-chapter", Attr: ""},
				ChapterTitle:   config.GoquerySelectorConfig{Selector: "chapter-title", Attr: ""},
				ChapterContent: config.GoquerySelectorConfig{Selector: "chapter-content", Attr: ""},
			},
			expectedParser: nil,
			expectError:    true,
		},
		{
			name: "missing chapter info selector",
			conf: &config.GoquerySelectorsConfig{
				Title:            config.GoquerySelectorConfig{Selector: "title", Attr: ""},
				Writer:           config.GoquerySelectorConfig{Selector: "writer", Attr: ""},
				BookType:         config.GoquerySelectorConfig{Selector: "book-type", Attr: ""},
				LastUpdate:       config.GoquerySelectorConfig{Selector: "update-date", Attr: ""},
				LastChapter:      config.GoquerySelectorConfig{Selector: "update-chapter", Attr: ""},
				BookChapterURL:   config.GoquerySelectorConfig{Selector: "book-chapter", Attr: "href"},
				BookChapterTitle: config.GoquerySelectorConfig{Selector: "book-chapter", Attr: ""},
			},
			expectedParser: nil,
			expectError:    true,
		},
		{
			name: "mismatch chapter url and title selector",
			conf: &config.GoquerySelectorsConfig{
				Title:            config.GoquerySelectorConfig{Selector: "title", Attr: ""},
				Writer:           config.GoquerySelectorConfig{Selector: "writer", Attr: ""},
				BookType:         config.GoquerySelectorConfig{Selector: "book-type", Attr: ""},
				LastUpdate:       config.GoquerySelectorConfig{Selector: "update-date", Attr: ""},
				LastChapter:      config.GoquerySelectorConfig{Selector: "update-chapter", Attr: ""},
				BookChapterURL:   config.GoquerySelectorConfig{Selector: "book-chapter-url", Attr: "href"},
				BookChapterTitle: config.GoquerySelectorConfig{Selector: "book-chapter-title", Attr: ""},
				ChapterTitle:     config.GoquerySelectorConfig{Selector: "chapter-title", Attr: ""},
				ChapterContent:   config.GoquerySelectorConfig{Selector: "chapter-content", Attr: ""},
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
				titleSelector:       Selector{"title", ""},
				writerSelector:      Selector{"writer", ""},
				bookTypeSelector:    Selector{"type", ""},
				lastUpdateSelector:  Selector{"date", ""},
				lastChapterSelector: Selector{"chapter", ""},
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
				titleSelector:       Selector{"title", ""},
				writerSelector:      Selector{"writer", ""},
				bookTypeSelector:    Selector{"type", ""},
				lastUpdateSelector:  Selector{"date", ""},
				lastChapterSelector: Selector{"chapter", ""},
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
				bookChapterURLSelector:   Selector{"ul>li", "href"},
				bookChapterTitleSelector: Selector{"ul>li", ""},
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
				bookChapterURLSelector:   Selector{"ul>li", "href"},
				bookChapterTitleSelector: Selector{"ul>li", ""},
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
				bookChapterURLSelector:   Selector{"ul>li", "href"},
				bookChapterTitleSelector: Selector{"ul>li", ""},
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
				ChapterTitleSelector:   Selector{"title", ""},
				ChapterContentSelector: Selector{"content", ""},
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
				ChapterTitleSelector:   Selector{"title", ""},
				ChapterContentSelector: Selector{"content", ""},
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
