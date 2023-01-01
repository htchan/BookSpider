package parser

import "github.com/htchan/BookSpider/internal/model"

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

func (fields *ParsedBookFields) Equal(other *ParsedBookFields) bool {
	if fields == other {
		return true
	} else if fields == nil || other == nil {
		return false
	}

	return fields.title == other.title && fields.writer == other.writer &&
		fields.bookType == other.bookType && fields.updateDate == other.updateDate &&
		fields.updateChapter == other.updateChapter
}

func (fields *ParsedChapterList) Equal(other *ParsedChapterList) bool {
	if fields == other {
		return true
	} else if fields == nil || other == nil {
		return false
	}

	if len(fields.chapters) != len(other.chapters) {
		return false
	}

	for i := range fields.chapters {
		chap, otherChap := fields.chapters[i], other.chapters[i]
		if chap.title != otherChap.title || chap.url != otherChap.url {
			return false
		}
	}

	return true
}

func (fields *ParsedChapterFields) Equal(other *ParsedChapterFields) bool {
	if fields == other {
		return true
	} else if fields == nil || other == nil {
		return false
	}

	return fields.title == other.title && fields.content == other.content
}
