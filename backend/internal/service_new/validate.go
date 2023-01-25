package service

import (
	"fmt"
	"strings"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
)

func isEnd(chapter string) bool {
	//TODO: fetch all chapter
	//hint: use book.generateEmptyChapters
	//TODO: check last n chapter to see if they contains any end keywords
	//hint: use len(chapters) and the n should come from book config
	chapter = strings.ReplaceAll(chapter, " ", "")
	fmt.Println(chapter)
	for _, keyword := range repo.ChapterEndKeywords {
		if strings.Contains(chapter, keyword) {
			return true
		}
	}
	return false
}

func (serv *ServiceImp) ValidateBookEnd(bk *model.Book) error {
	isUpdated := false
	if isEnd(bk.UpdateChapter) {
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
