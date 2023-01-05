package parse

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/model"
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

var _ Parser = (*GoqueryParser)(nil)

var (
	ErrBookInfoSelectorEmpty          = errors.New("book info selector is empty")
	ErrChapterListSelectorEmpty       = errors.New("chapter list selector is empty")
	ErrChapterSelectorEmpty           = errors.New("chapter selector is empty")
	ErrParseBookFieldsNotFound        = errors.New("parse book fail: fields not found")
	ErrParseChapterFieldsNotFound     = errors.New("parse chapter fail: fields not found")
	ErrParseChapterListNoChapterFound = errors.New("no chapters found")
)

func LoadGoqueryParser(conf *config.GoquerySelectorConfig) (*GoqueryParser, error) {
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

func IsNewBook(fields *ParsedBookFields, bk *model.Book) bool {
	// TODO
	return false
}

func IsUpdatedBook(fields *ParsedBookFields, bk *model.Book) bool {
	// TODO
	return false
}

func (parser *GoqueryParser) ParseBook(html string) (*ParsedBookFields, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse book fail: %w", err)
	}

	var fields ParsedBookFields

	fields.title = doc.Find(parser.titleSelector).Text()
	fields.writer = doc.Find(parser.writerSelector).Text()
	fields.bookType = doc.Find(parser.bookTypeSelector).Text()
	fields.updateDate = doc.Find(parser.lastUpdateSelector).Text()
	fields.updateChapter = doc.Find(parser.lastChapterSelector).Text()

	if fields.title == "" || fields.writer == "" || fields.bookType == "" ||
		fields.updateDate == "" || fields.updateChapter == "" {
		return nil, ErrParseBookFieldsNotFound
	}

	return &fields, nil
}

func (parser *GoqueryParser) ParseChapterList(html string) (*ParsedChapterList, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse chapter list fail: %w", err)
	}

	var chapters ParsedChapterList

	doc.Find(parser.bookChapterSelector).Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		title := s.Text()
		chapters.chapters = append(chapters.chapters, struct {
			url   string
			title string
		}{url: url, title: title})
	})

	if len(chapters.chapters) == 0 {
		return nil, ErrParseChapterListNoChapterFound
	}

	return &chapters, nil
}

func (parser *GoqueryParser) ParseChapter(html string) (*ParsedChapterFields, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parse chapter fail: %w", err)
	}

	var fields ParsedChapterFields

	fields.title = doc.Find(parser.ChapterTitleSelector).Text()

	doc.Find(parser.ChapterContentSelector).Each(func(i int, s *goquery.Selection) {
		fields.content += strings.TrimSpace(s.Text()) + "\n"
	})

	if fields.title == "" || fields.content == "" {
		return nil, ErrParseChapterFieldsNotFound
	}
	return &fields, nil
}

func (parser *GoqueryParser) Equal(target *GoqueryParser) bool {
	if parser == nil && target == nil {
		return true
	} else if parser == nil || target == nil {
		return false
	}

	return parser.titleSelector == target.titleSelector &&
		parser.writerSelector == target.writerSelector &&
		parser.bookTypeSelector == target.bookTypeSelector &&
		parser.lastUpdateSelector == target.lastUpdateSelector &&
		parser.lastChapterSelector == target.lastChapterSelector &&
		parser.bookChapterSelector == target.bookChapterSelector &&
		parser.ChapterTitleSelector == target.ChapterTitleSelector &&
		parser.ChapterContentSelector == target.ChapterContentSelector
}
