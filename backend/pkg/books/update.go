package books

import (
	"errors"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/logging"
	"github.com/htchan/ApiParser"
)

func (book *Book) fetchInfo() (title, writer, typeString, updateDate, updateChapter string, err error) {
	defer utils.Recover(func() {
		title, writer, typeString, updateDate, updateChapter = "", "", "", "", ""
	})
	// get online resource, try maximum 10 times if it keeps failed
	html, trial := utils.GetWeb(
		book.config.BaseUrl, 10, book.config.Decoder, book.config.ConstSleep)
	if trial > 0 {
		// logging.LogBookEvent(book.String(), "source_info", "trial", trial)
	}
	err = book.validHTML(html)
	utils.CheckError(err)
	
	responseApi := ApiParser.Parse(book.config.SourceKey + ".info", html)
	title, okTitle := responseApi.Data["Title"]
	writer, okWriter := responseApi.Data["Writer"]
	typeString, okTypeString := responseApi.Data["Type"]
	updateDate, okUpdateDate := responseApi.Data["LastUpdate"]
	updateChapter, okUpdateChapter := responseApi.Data["LastChapter"]
	if !(okTitle && okWriter && okTypeString && okUpdateDate && okUpdateChapter) {
		err = errors.New("some data not found")
		panic(err)
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
