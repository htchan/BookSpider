package goquery

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/parse"
)

type GoqueryParser struct {
	titleSelector          string
	writerSelector         string
	bookTypeSelector       string
	lastUpdateSelector     string
	lastChapterSelector    string
	bookChapterSelector    string
	ChapterTitleSelector   string
	ChapterContentSelector string
}

var (
	ErrBookInfoSelectorEmpty    = errors.New("book info selector is empty")
	ErrChapterListSelectorEmpty = errors.New("chapter list selector is empty")
	ErrChapterSelectorEmpty     = errors.New("chapter selector is empty")
)

var _ parse.Parser = (*GoqueryParser)(nil)

func LoadParser(conf *config.GoquerySelectorConfig) (*GoqueryParser, error) {
	if conf.Title == "" || conf.Writer == "" || conf.BookType == "" ||
		conf.LastUpdate == "" || conf.LastChapter == "" {
		return nil, ErrBookInfoSelectorEmpty
	}

	if conf.BookChapter == "" {
		return nil, ErrChapterListSelectorEmpty
	}

	if conf.ChapterContent == "" || conf.ChapterTitle == "" {
		return nil, ErrChapterSelectorEmpty
	}

	return &GoqueryParser{
		titleSelector:          conf.Title,
		writerSelector:         conf.Writer,
		bookTypeSelector:       conf.BookType,
		lastUpdateSelector:     conf.LastUpdate,
		lastChapterSelector:    conf.LastChapter,
		bookChapterSelector:    conf.BookChapter,
		ChapterTitleSelector:   conf.ChapterTitle,
		ChapterContentSelector: conf.ChapterContent,
	}, nil
}

func (parser *GoqueryParser) ParseBook(html string) (*parse.ParsedBookFields, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse book fail: %w", err)
	}

	fields := parse.NewParsedBookFields(
		doc.Find(parser.titleSelector).Text(),
		doc.Find(parser.writerSelector).Text(),
		doc.Find(parser.bookTypeSelector).Text(),
		doc.Find(parser.lastUpdateSelector).Text(),
		doc.Find(parser.lastChapterSelector).Text(),
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

	doc.Find(parser.bookChapterSelector).Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		title := s.Text()
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
	title := doc.Find(parser.ChapterTitleSelector).Text()

	doc.Find(parser.ChapterContentSelector).Each(func(i int, s *goquery.Selection) {
		content += strings.TrimSpace(s.Text()) + "\n"
	})

	fields := parse.NewParsedChapterFields(title, content)

	err = fields.Validate()
	if err != nil {
		return nil, err
	}

	return fields, nil
}
