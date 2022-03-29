package books

import (
	// "errors"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/logging"
)

func (book *Book) fetchInfo() (title, writer, typeString, updateDate, updateChapter string, err error) {
	operation := [5]struct{
		regex string
		result *string
	} {
		{book.config.TitleRegex, &title},
		{book.config.WriterRegex, &writer},
		{book.config.TypeRegex, &typeString},
		{book.config.LastUpdateRegex, &updateDate},
		{book.config.LastChapterRegex, &updateChapter},
	}
	defer utils.Recover(func() {
		for i := 0; i < 5; i++ {
			*operation[i].result = ""
		}
	})
	// get online resource, try maximum 10 times if it keeps failed
	html, _ := utils.GetWeb(
		book.config.BaseUrl, 10, book.config.Decoder, book.config.CONST_SLEEP)
	err = book.validHTML(html)
	utils.CheckError(err)
	
	// extract info from source
	var result string
	for i := 0; i < 5; i++ {
		result, err = utils.Search(html, operation[i].regex)
		utils.CheckError(err)
		*operation[i].result = result
	}
	return
}

func (book *Book) Update() bool {
	title, writer, typeString, updateDate, updateChapter, err := book.fetchInfo()
	if err != nil {
		if book.GetStatus() == database.Error {
			book.SetError(err)
		}
		logging.LogBookEvent(book.String(), "update", "fetch-fail", err)
		return false
	}
	// check difference
	update := false
	if title != book.GetTitle() || writer != book.GetWriter() || typeString != book.GetType() {
		update = true
		if book.GetStatus() != database.Error {
			book.bookRecord.HashCode = database.GenerateHash()
			logging.LogBookEvent(book.String(), "update", "new-record-created", nil)
		}
		book.SetStatus(database.InProgress)
		book.SetError(nil)
	} else if updateDate != book.GetUpdateDate() || updateChapter != book.GetUpdateChapter() {
		update = true
		book.SetStatus(database.InProgress)
		book.SetError(nil)
	}
	if update {
		// sync with online info
		book.SetTitle(title)
		book.SetWriter(writer)
		book.SetType(typeString)
		book.SetUpdateDate(updateDate)
		book.SetUpdateChapter(updateChapter)
		logging.LogBookEvent(book.String(), "update", "success", nil)
	}
	return update
}
