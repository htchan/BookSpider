package book

import (
	"github.com/htchan/BookSpider/internal/book/model"
	"strings"
	"time"
)

var (
	endKeywords = []string{
		"后记", "後記", "新书", "新書", "结局", "結局", "感言",
		"尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本" /*"结束", "結束", */, "完結",
		"完结", "终结", "終結", "番外", "结尾", "結尾", "全书完", "全書完", "全本完",
	}
)

func (book *Book) isEnd() bool {
	for _, keyword := range endKeywords {
		if strings.Contains(book.UpdateChapter, keyword) {
			return true
		}
	}
	return false
}

func (book *Book) Process() bool {
	isUpdated, errUpdate := book.Update()
	if errUpdate != nil {
		return false
	}
	updateDate, err := time.Parse(book.BookConfig.UpdateDateLayout, book.UpdateDate)
	if isUpdated && book.isEnd() ||
		book.Status != model.Download && err == nil && updateDate.Add(365*24*time.Hour).Before(time.Now()) {
		isUpdated = true
		book.Status = model.End
	}
	if book.Status == model.End {
		//TODO: log the error
		/*errDownload := */
		book.Download()
	}
	return isUpdated
}
