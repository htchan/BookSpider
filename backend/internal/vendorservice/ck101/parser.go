package ck101

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
)

func (p *VendorService) ParseDoc(body string) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(strings.NewReader(body))
}

func (p *VendorService) ParseBook(body string) (*vendor.BookInfo, error) {
	doc, docErr := p.ParseDoc(body)
	if docErr != nil {
		return nil, fmt.Errorf("parse body fail: %w", docErr)
	}

	var parseErr error

	// parse title
	title := doc.Find(bookTitleGoquerySelector).AttrOr("content", "")
	if title == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookTitleNotFound)
	}

	// parse writer
	writer := doc.Find(bookWriterGoquerySelector).AttrOr("content", "")
	if writer == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookWriterNotFound)
	}

	// parse type
	bookType := doc.Find(bookTypeGoquerySelector).AttrOr("content", "")
	if bookType == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookTypeNotFound)
	}

	// parse date
	date := vendor.GetGoqueryContentWithoutChildren(doc.Find(bookDateGoquerySelector))
	date = strings.ReplaceAll(date, "更新时间：", "")
	if date == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookDateNotFound)
	}

	// parse chapter
	chapter := vendor.GetGoqueryContentWithoutChildren(doc.Find(bookChapterGoquerySelector))
	if chapter == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookChapterNotFound)
	}

	if parseErr != nil {
		parseErr = errors.Join(parseErr, vendor.ErrFieldsNotFound)
	}

	return &vendor.BookInfo{
		Title:         title,
		Writer:        writer,
		Type:          bookType,
		UpdateDate:    date,
		UpdateChapter: chapter,
		IsEnd:         vendor.CheckChapterEnd(chapter) || vendor.CheckDateEnd(date),
	}, parseErr
}

func (p *VendorService) ParseChapterList(_, body string) (vendor.ChapterList, error) {
	doc, docErr := p.ParseDoc(body)
	if docErr != nil {
		return nil, fmt.Errorf("parse body fail: %w", docErr)
	}

	var chapterList vendor.ChapterList
	var parseErr error
	doc.Find(chapterListItemGoquerySelector).Each(func(i int, s *goquery.Selection) {
		url := s.AttrOr("href", "")
		if url == "" {
			parseErr = errors.Join(
				parseErr,
				fmt.Errorf("parse chapter url fail: %d, %w", i, vendor.ErrChapterListUrlNotFound),
			)
		}

		title := vendor.GetGoqueryContentWithoutChildren(s)
		if title == "" {
			parseErr = errors.Join(
				parseErr,
				fmt.Errorf("parse chapter url fail: %d, %w", i, vendor.ErrChapterListTitleNotFound),
			)
		}

		chapterList = append(chapterList, vendor.ChapterListInfo{
			URL:   p.ChapterListURL(url),
			Title: title,
		})
	})

	if len(chapterList) == 0 {
		return nil, vendor.ErrChapterListEmpty
	}

	if parseErr != nil {
		parseErr = errors.Join(parseErr, vendor.ErrFieldsNotFound)
	}

	return chapterList, parseErr
}

func (p *VendorService) ParseChapter(body string) (*vendor.ChapterInfo, error) {
	doc, docErr := p.ParseDoc(body)
	if docErr != nil {
		return nil, fmt.Errorf("parse body fail: %w", docErr)
	}

	var parseErr error

	// parse title
	title := vendor.GetGoqueryContentWithoutChildren(doc.Find(chapterTitleGoquerySelector))
	if title == "" {
		parseErr = errors.Join(parseErr, vendor.ErrChapterTitleNotFound)
	}

	// parse content
	content := vendor.GetGoqueryContentWithoutChildren(doc.Find(chapterContentGoquerySelector))
	if content == "" {
		parseErr = errors.Join(parseErr, vendor.ErrChapterContentNotFound)
	}

	if parseErr != nil {
		parseErr = errors.Join(parseErr, vendor.ErrFieldsNotFound)
	}

	return &vendor.ChapterInfo{
		Title: title,
		Body:  content,
	}, parseErr
}

func (p *VendorService) IsAvailable(body string) bool {
	return strings.Contains(body, "黃金屋")
}

func (p *VendorService) FindMissingIds(ids []int) []int {
	var missingIDs []int

	sort.Ints(ids)

	idPointer, i := 0, 1
	for idPointer < len(ids) && ids[len(ids)-1] > i {
		if i == ids[idPointer] {
			i++
			idPointer++
		} else if i > ids[idPointer] {
			idPointer++
		} else if i < ids[idPointer] {
			missingIDs = append(missingIDs, i)
			i++
		}
	}

	return missingIDs
}
