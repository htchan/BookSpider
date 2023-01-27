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
	selector string
	attr     string
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
	if s.attr == "" {
		return selection.Text()
	}
	return selection.AttrOr(s.attr, "")
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
		titleSelector:            Selector{selector: conf.Title.Selector, attr: conf.Title.Attr},
		writerSelector:           Selector{selector: conf.Writer.Selector, attr: conf.Writer.Attr},
		bookTypeSelector:         Selector{selector: conf.BookType.Selector, attr: conf.BookType.Attr},
		lastUpdateSelector:       Selector{selector: conf.LastUpdate.Selector, attr: conf.LastUpdate.Attr},
		lastChapterSelector:      Selector{selector: conf.LastChapter.Selector, attr: conf.LastChapter.Attr},
		bookChapterURLSelector:   Selector{selector: conf.BookChapterURL.Selector, attr: conf.BookChapterURL.Attr},
		bookChapterTitleSelector: Selector{selector: conf.BookChapterTitle.Selector, attr: conf.BookChapterTitle.Attr},
		ChapterTitleSelector:     Selector{selector: conf.ChapterTitle.Selector, attr: conf.ChapterTitle.Attr},
		ChapterContentSelector:   Selector{selector: conf.ChapterContent.Selector, attr: conf.ChapterContent.Attr},
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
		title := parser.bookChapterTitleSelector.Parse(s)
		chapters.Append(url, title)
	})

	err = chapters.Validate()
	if err != nil {
		return nil, err
	}

	return &chapters, nil
}

func (parser *GoqueryParser) ParseChapter(html string) (*parse.ParsedChapterFields, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse chapter fail: %w", err)
	}

	content := ""
	title := parser.ChapterTitleSelector.Parse(doc.Find(parser.ChapterTitleSelector.selector))

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
