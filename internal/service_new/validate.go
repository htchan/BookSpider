package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/htchan/BookSpider/internal/model"
)

func isEnd(bk *model.Book) bool {
	//TODO: fetch all chapter
	//hint: use book.generateEmptyChapters
	//TODO: check last n chapter to see if they contains any end keywords
	//hint: use len(chapters) and the n should come from book config
	if bk.UpdateDate < strconv.Itoa(time.Now().Year()-1) {
		return true
	}

	chapter := strings.ReplaceAll(bk.UpdateChapter, " ", "")
	for _, keyword := range model.ChapterEndKeywords {
		if strings.Contains(chapter, keyword) {
			return true
		}
	}
	return false
}

func (serv *ServiceImp) ValidateBookEnd(bk *model.Book) error {
	isUpdated := false
	if isEnd(bk) {
		if bk.Status != model.StatusEnd {
			bk.IsDownloaded = false
			bk.Status = model.StatusEnd
			isUpdated = true
		}
	} else {
		if bk.Status != model.StatusInProgress {
			bk.Status = model.StatusInProgress
			isUpdated = true
		}
	}

	if isUpdated {
		err := serv.rpo.UpdateBook(bk)
		if err != nil {
			return fmt.Errorf("update book in DB fail: %w", err)
		}
	}

	return nil
}

func (serv *ServiceImp) ValidateEnd() error {
	// TODO: move this to Update service
	return serv.rpo.UpdateBooksStatus()
}
