package book

import (
	"fmt"
	"strings"

	"github.com/htchan/BookSpider/internal/model"
)

var (
	endKeywords = []string{
		// traditional chinese
		"番外", "結局", "新書", "完結", "尾聲", "感言", "後記", "完本",
		"全書完", "全文完", "全文終", "全文結", "劇終", "（完）", "終章",
		"外傳", "結尾",
		// simplified chinese
		"番外", "结局", "新书", "完结", "尾声", "感言", "后记", "完本",
		"全书完", "全文完", "全文终", "全文结", "剧终", "（完）", "终章",
		"外传", "结尾",
	}
)

func isEnd(chapter string) bool {
	//TODO: fetch all chapter
	//hint: use book.generateEmptyChapters
	//TODO: check last n chapter to see if they contains any end keywords
	//hint: use len(chapters) and the n should come from book config
	chapter = strings.ReplaceAll(chapter, " ", "")
	fmt.Println(chapter)
	for _, keyword := range endKeywords {
		if strings.Contains(chapter, keyword) {
			return true
		}
	}
	return false
}

func Validate(bk *model.Book) (bool, error) {
	isUpdated := false
	if isEnd(bk.UpdateChapter) {
		if bk.Status != model.End {
			bk.IsDownloaded = false
			isUpdated = true
		}
		bk.Status = model.End
	} else {
		if bk.Status != model.InProgress {
			isUpdated = true
		}
		bk.Status = model.InProgress
	}
	return isUpdated, nil
}
