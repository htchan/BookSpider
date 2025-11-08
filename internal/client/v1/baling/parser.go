package baling

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/htchan/BookSpider/internal/client/v1"
)

func parseBook(body string) (*client.BookInfo, error) {
	doc, docErr := client.ParseDoc(body)
	if docErr != nil {
		return nil, fmt.Errorf("parse body fail: %w", docErr)
	}

	var parseErr error

	// parse title
	title := doc.Find(bookTitleGoquerySelector).AttrOr("content", "")
	if title == "" {
		parseErr = errors.Join(parseErr, client.ErrBookTitleNotFound)
	}

	// parse writer
	writer := doc.Find(bookWriterGoquerySelector).AttrOr("content", "")
	if writer == "" {
		parseErr = errors.Join(parseErr, client.ErrBookWriterNotFound)
	}

	// parse type
	bookType := doc.Find(bookTypeGoquerySelector).AttrOr("content", "")
	if bookType == "" {
		parseErr = errors.Join(parseErr, client.ErrBookTypeNotFound)
	}

	// parse dateStr
	date := time.Now().UTC().Truncate(time.Minute)
	var parseDateErr error
	dateStr := client.GetGoqueryContentWithoutChildren(doc.Find(bookDateGoquerySelector))
	dateStr = strings.ReplaceAll(dateStr, "更新时间：", "")
	if dateStr == "" {
		parseErr = errors.Join(parseErr, client.ErrBookDateNotFound)
	} else {
		date, parseDateErr = time.Parse("2006-01-02", dateStr)
		if parseDateErr != nil {
			parseErr = errors.Join(parseErr, client.ErrBookDateParseFail)
		}
	}

	// parse chapter
	chapter := client.GetGoqueryContentWithoutChildren(doc.Find(bookChapterGoquerySelector))
	if chapter == "" {
		parseErr = errors.Join(parseErr, client.ErrBookChapterNotFound)
	}

	if parseErr != nil {
		parseErr = errors.Join(parseErr, client.ErrFieldsNotFound)
	}

	return &client.BookInfo{
		Title:         title,
		Author:        writer,
		Type:          bookType,
		UpdateDate:    date.UTC().Truncate(time.Second),
		UpdateChapter: chapter,
	}, parseErr
}

func parseChapterList(body string) (client.ChapterEntryList, error) {
	doc, docErr := client.ParseDoc(body)
	if docErr != nil {
		return nil, fmt.Errorf("parse body fail: %w", docErr)
	}

	var chapterList client.ChapterEntryList
	var parseErr error
	doc.Find(chapterListItemGoquerySelector).Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("href", "")
		if url == "" {
			parseErr = errors.Join(
				parseErr,
				fmt.Errorf("parse chapter url fail: %d, %w", i, client.ErrChapterListUrlNotFound),
			)
		}

		title := client.GetGoqueryContentWithoutChildren(s)
		if title == "" {
			parseErr = errors.Join(
				parseErr,
				fmt.Errorf("parse chapter url fail: %d, %w", i, client.ErrChapterListTitleNotFound),
			)
		}

		chapterList = append(chapterList, client.ChapterEntry{
			URL:   chapterURL(url),
			Title: title,
		})
	})

	if len(chapterList) == 0 {
		return nil, client.ErrChapterListEmpty
	}

	if parseErr != nil {
		parseErr = errors.Join(parseErr, client.ErrFieldsNotFound)
	}

	return chapterList, parseErr
}

func parseChapter(body string) (*client.ChapterContent, error) {
	doc, docErr := client.ParseDoc(body)
	if docErr != nil {
		return nil, fmt.Errorf("parse body fail: %w", docErr)
	}

	var parseErr error

	// parse title
	title := client.GetGoqueryContentWithoutChildren(doc.Find(chapterTitleGoquerySelector))
	if title == "" {
		parseErr = errors.Join(parseErr, client.ErrChapterTitleNotFound)
	}

	// parse content
	content := client.GetGoqueryContentWithoutChildren(doc.Find(chapterContentGoquerySelector))
	if content == "" {
		parseErr = errors.Join(parseErr, client.ErrChapterContentNotFound)
	}

	if parseErr != nil {
		parseErr = errors.Join(parseErr, client.ErrFieldsNotFound)
	}

	return &client.ChapterContent{
		Title: title,
		Body:  content,
	}, parseErr
}

func isAvailable(body string) bool {
	return strings.Contains(body, "80txt")
}
