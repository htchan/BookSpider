package parse

import (
	"errors"

	"github.com/htchan/BookSpider/internal/model"
)

type ParsedBookFields struct {
	title         string
	writer        string
	bookType      string
	updateDate    string
	updateChapter string
}

type ParsedChapterList struct {
	chapters []struct {
		url   string
		title string
	}
}

type ParsedChapterFields struct {
	title   string
	content string
}

var (
	ErrParseBookFieldsNotFound        = errors.New("parse book fail: fields not found")
	ErrParseChapterFieldsNotFound     = errors.New("parse chapter fail: fields not found")
	ErrParseChapterListNoChapterFound = errors.New("no chapters found")
)

func IsNewBook(fields *ParsedBookFields, bk *model.Book) bool {
	return (bk.Title != "" && bk.Title != fields.title) ||
		(bk.Writer.Name != "" && bk.Writer.Name != fields.writer) ||
		(bk.Type != "" && bk.Type != fields.bookType)
}

func IsUpdatedBook(fields *ParsedBookFields, bk *model.Book) bool {
	return bk.Title != fields.title ||
		bk.Writer.Name != fields.writer ||
		bk.Type != fields.bookType ||
		bk.UpdateChapter != fields.updateChapter ||
		bk.UpdateDate != fields.updateDate
}

func NewParsedBookFields(title, writer, bookType, updateDate, updateChapter string) *ParsedBookFields {
	fields := ParsedBookFields{
		title:         title,
		writer:        writer,
		bookType:      bookType,
		updateDate:    updateDate,
		updateChapter: updateChapter,
	}

	return &fields
}

func (fields *ParsedBookFields) Validate() error {
	if fields.title == "" || fields.writer == "" || fields.bookType == "" ||
		fields.updateDate == "" || fields.updateChapter == "" {
		return ErrParseBookFieldsNotFound
	}

	return nil
}

func (fields *ParsedChapterList) Append(url, title string) {
	fields.chapters = append(fields.chapters, struct {
		url   string
		title string
	}{url, title})
}

func (fields *ParsedChapterList) Validate() error {
	if len(fields.chapters) == 0 {
		return ErrParseChapterListNoChapterFound
	}

	return nil
}

func NewParsedChapterFields(title, content string) *ParsedChapterFields {
	fields := ParsedChapterFields{
		title:   title,
		content: content,
	}
	return &fields
}

func (fields *ParsedChapterFields) Validate() error {
	if fields.title == "" || fields.content == "" {
		return ErrParseChapterFieldsNotFound
	}

	return nil
}

func (fields *ParsedBookFields) Populate(bk *model.Book) {
	bk.Title = fields.title
	bk.Writer.Name = fields.writer
	bk.Type = fields.bookType
	bk.UpdateDate = fields.updateDate
	bk.UpdateChapter = fields.updateChapter
}

func (fields *ParsedChapterList) Populate(chapters *model.Chapters) {
	for i, item := range fields.chapters {
		*chapters = append(*chapters, model.Chapter{
			Index: i,
			URL:   item.url,
			Title: item.title,
		})
	}
}

func (fields *ParsedChapterFields) Populate(ch *model.Chapter) {
	ch.Title = fields.title
	ch.Content = fields.content
}
