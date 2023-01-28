package goquery

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/parse"
)

type Selector struct {
	selector        string
	attr            string
	unwantedContent []string
}

type GoqueryParser struct {
	titleSelector            Selector
	writerSelector           Selector
	bookTypeSelector         Selector
	lastUpdateSelector       Selector
	lastChapterSelector      Selector
	bookChapterURLSelector   Selector
	bookChapterTitleSelector Selector
	ChapterTitleSelector     Selector
	ChapterContentSelector   Selector
}

var (
	ErrBookInfoSelectorEmpty       = errors.New("book info selector is empty")
	ErrChapterListSelectorEmpty    = errors.New("chapter list selector is empty")
	ErrBookChapterSelectorMismatch = errors.New("book chapter url selector different from title selector")
	ErrChapterSelectorEmpty        = errors.New("chapter selector is empty")
)

var _ parse.Parser = (*GoqueryParser)(nil)

func (s *Selector) Parse(selection *goquery.Selection) string {
	result := selection.AttrOr(s.attr, "")
	if s.attr == "" {
		result = selection.Children().Remove().End().Text()
	}

	for _, content := range s.unwantedContent {
		result = strings.ReplaceAll(result, content, "")
	}

	return result
}

func LoadParser(conf *config.GoquerySelectorsConfig) (*GoqueryParser, error) {
	if conf.Title.Selector == "" || conf.Writer.Selector == "" || conf.BookType.Selector == "" ||
		conf.LastUpdate.Selector == "" || conf.LastChapter.Selector == "" {
		return nil, ErrBookInfoSelectorEmpty
	}

	if conf.BookChapterURL.Selector == "" && conf.BookChapterTitle.Selector == "" {
		return nil, ErrChapterListSelectorEmpty
	}

	if conf.BookChapterURL.Selector != conf.BookChapterTitle.Selector {
		return nil, ErrBookChapterSelectorMismatch
	}

	if conf.ChapterContent.Selector == "" || conf.ChapterTitle.Selector == "" {
		return nil, ErrChapterSelectorEmpty
	}

	return &GoqueryParser{
		titleSelector:            Selector{selector: conf.Title.Selector, attr: conf.Title.Attr, unwantedContent: conf.Title.UnwantedContent},
		writerSelector:           Selector{selector: conf.Writer.Selector, attr: conf.Writer.Attr, unwantedContent: conf.Writer.UnwantedContent},
		bookTypeSelector:         Selector{selector: conf.BookType.Selector, attr: conf.BookType.Attr, unwantedContent: conf.BookType.UnwantedContent},
		lastUpdateSelector:       Selector{selector: conf.LastUpdate.Selector, attr: conf.LastUpdate.Attr, unwantedContent: conf.LastUpdate.UnwantedContent},
		lastChapterSelector:      Selector{selector: conf.LastChapter.Selector, attr: conf.LastChapter.Attr, unwantedContent: conf.LastChapter.UnwantedContent},
		bookChapterURLSelector:   Selector{selector: conf.BookChapterURL.Selector, attr: conf.BookChapterURL.Attr, unwantedContent: conf.BookChapterURL.UnwantedContent},
		bookChapterTitleSelector: Selector{selector: conf.BookChapterTitle.Selector, attr: conf.BookChapterTitle.Attr, unwantedContent: conf.BookChapterTitle.UnwantedContent},
		ChapterTitleSelector:     Selector{selector: conf.ChapterTitle.Selector, attr: conf.ChapterTitle.Attr, unwantedContent: conf.ChapterTitle.UnwantedContent},
		ChapterContentSelector:   Selector{selector: conf.ChapterContent.Selector, attr: conf.ChapterContent.Attr, unwantedContent: conf.ChapterContent.UnwantedContent},
	}, nil
}

func (parser *GoqueryParser) ParseBook(html string) (*parse.ParsedBookFields, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse book fail: %w", err)
	}

	fields := parse.NewParsedBookFields(
		parser.titleSelector.Parse(doc.Find(parser.titleSelector.selector)),
		parser.writerSelector.Parse(doc.Find(parser.writerSelector.selector)),
		parser.bookTypeSelector.Parse(doc.Find(parser.bookTypeSelector.selector)),
		parser.lastUpdateSelector.Parse(doc.Find(parser.lastUpdateSelector.selector)),
		parser.lastChapterSelector.Parse(doc.Find(parser.lastChapterSelector.selector)),
	)

	err = fields.Validate()
	if err != nil {
		return nil, err
	}

	return fields, nil
}

func (parser *GoqueryParser) ParseChapterList(html string) (*parse.ParsedChapterList, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse chapter list fail: %w", err)
	}

	var chapters parse.ParsedChapterList

	doc.Find(parser.bookChapterURLSelector.selector).Each(func(i int, s *goquery.Selection) {
		url := parser.bookChapterURLSelector.Parse(s)
		title := strings.TrimSpace(parser.bookChapterTitleSelector.Parse(s))
		chapters.Append(url, title)
	})

	err = chapters.Validate()
	if err != nil {
		return nil, err
	}

	return &chapters, nil
}

func (parser *GoqueryParser) ParseChapter(html string) (*parse.ParsedChapterFields, error) {
	replaceItems := []struct {
		old, new string
	}{
		{"<br />", "\n"},
		{"&nbsp;", ""},
		{"<b>", ""},
		{"</b>", ""},
		{"<p>", ""},
		{"</p>", ""},
		{"                ", ""},
		{"<p/>", "\n"},
	}
	for _, replaceItem := range replaceItems {
		html = strings.ReplaceAll(
			html, replaceItem.old, replaceItem.new)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse chapter fail: %w", err)
	}

	content := ""
	title := strings.TrimSpace(parser.ChapterTitleSelector.Parse(doc.Find(parser.ChapterTitleSelector.selector)))

	doc.Find(parser.ChapterContentSelector.selector).Each(func(i int, s *goquery.Selection) {
		content += strings.TrimSpace(parser.ChapterContentSelector.Parse(s)) + "\n"
	})

	fields := parse.NewParsedChapterFields(title, content)

	err = fields.Validate()
	if err != nil {
		return nil, err
	}

	return fields, nil
}
