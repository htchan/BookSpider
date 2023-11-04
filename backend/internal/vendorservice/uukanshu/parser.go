package uukanshu

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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
	title := doc.Find(bookTitleGoquerySelector).AttrOr("title", "")
	title = strings.ReplaceAll(title, "最新章节", "")
	if title == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookTitleNotFound)
	}

	// parse writer
	writer := vendor.GetGoqueryContent(doc.Find(bookWriterGoquerySelector))
	if writer == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookWriterNotFound)
	}

	// parse type
	bookType := vendor.GetGoqueryContent(doc.Find(bookTypeGoquerySelector))
	if bookType == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookTypeNotFound)
	}

	// parse dateStr
	dateStr := vendor.GetGoqueryContent(doc.Find(bookDateGoquerySelector))
	dateStr = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(dateStr, "更新时间：", ""), " ", ""), "\t", ""), "\n", "")
	if dateStr == "" {
		parseErr = errors.Join(parseErr, vendor.ErrBookDateNotFound)
	}

	date := time.Now().UTC().Truncate(24 * time.Hour)
	if strings.Contains(dateStr, "年") {
		yr, _ := strconv.Atoi(dateStr[:strings.Index(dateStr, "年")])
		date = date.AddDate(-yr, 0, 0)
		date = time.Date(date.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	} else if strings.Contains(dateStr, "月") {
		mn, _ := strconv.Atoi(dateStr[:strings.Index(dateStr, "月")])
		date = date.AddDate(0, -mn, 0)
		date = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else if strings.Contains(dateStr, "日") {
		day, _ := strconv.Atoi(dateStr[:strings.Index(dateStr, "日")])
		date = date.AddDate(0, 0, -day)
	}

	// parse chapter
	chapter := vendor.GetGoqueryContent(doc.Find(bookChapterGoquerySelector))
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
		UpdateDate:    date.Format(time.DateOnly),
		UpdateChapter: chapter,
	}, parseErr
}

func (p *VendorService) ParseChapterList(body string) (vendor.ChapterList, error) {
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

		title := vendor.GetGoqueryContent(s)
		if title == "" {
			parseErr = errors.Join(
				parseErr,
				fmt.Errorf("parse chapter url fail: %d, %w", i, vendor.ErrChapterListTitleNotFound),
			)
		}

		chapterList = append(chapterList, vendor.ChapterListInfo{
			URL:   url,
			Title: title,
		})
	})

	if len(chapterList) == 0 {
		return nil, vendor.ErrChapterListEmpty
	}

	if parseErr != nil {
		parseErr = errors.Join(parseErr, vendor.ErrFieldsNotFound)
	}

	resultChapterList := make(vendor.ChapterList, len(chapterList))
	for i := range chapterList {
		resultChapterList[i] = chapterList[len(chapterList)-1-i]
	}

	return resultChapterList, parseErr
}

func (p *VendorService) ParseChapter(body string) (*vendor.ChapterInfo, error) {
	doc, docErr := p.ParseDoc(body)
	if docErr != nil {
		return nil, fmt.Errorf("parse body fail: %w", docErr)
	}

	var parseErr error

	// parse title
	title := vendor.GetGoqueryContent(doc.Find(chapterTitleGoquerySelector))
	if title == "" {
		parseErr = errors.Join(parseErr, vendor.ErrChapterTitleNotFound)
	}

	// parse content
	content := vendor.GetGoqueryContent(doc.Find(chapterContentGoquerySelector))
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
