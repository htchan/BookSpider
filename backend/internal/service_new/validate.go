package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
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
	for _, keyword := range repo.ChapterEndKeywords {
		if strings.Contains(chapter, keyword) {
			return true
		}
	}
	return false
}

func (serv *ServiceImp) ValidateBookEnd(bk *model.Book) error {
	isUpdated := false
	if isEnd(bk) {
		if bk.Status != model.End {
			bk.IsDownloaded = false
			bk.Status = model.End
			isUpdated = true
		}
	} else {
		if bk.Status != model.InProgress {
			bk.Status = model.InProgress
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
