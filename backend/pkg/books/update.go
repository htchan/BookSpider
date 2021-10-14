package books

import (
	"errors"
	"github.com/htchan/BookSpider/internal/utils"
)

func (book *Book) extractInfo() (string, string, string, string, string, error) {
	// get online resource, try maximum 10 times if it keeps failed
	html, trial := utils.GetWeb(book.metaInfo.baseUrl, 10, book.decoder)
	if _, err := utils.Search(html, book.metaInfo.titleRegex); err != nil || !book.validHTML(html, book.metaInfo.baseUrl, trial) {
		book.Log(map[string]interface{}{
			"url": book.metaInfo.baseUrl, "error": "extract base html fail", "stage": "update",
		})
		return "", "", "", "", "", errors.New("invalid base page html")
	}
	// extract info from source
	var err error
	title, errTitle := utils.Search(html, book.metaInfo.titleRegex)
	if errTitle != nil {
		err = errTitle
	}
	writer, errWriter := utils.Search(html, book.metaInfo.writerRegex)
	if errWriter != nil {
		err = errWriter
	}
	typeName, errTypeName := utils.Search(html, book.metaInfo.typeRegex)
	if errTypeName != nil {
		err = errTypeName
	}
	lastUpdate, errLastUpdate := utils.Search(html, book.metaInfo.lastUpdateRegex)
	if errLastUpdate != nil {
		err = errLastUpdate
	}
	lastChapter, errLastChapter := utils.Search(html, book.metaInfo.lastChapterRegex)
	if errLastChapter != nil {
		err = errLastChapter
	}
	if err != nil {
		book.Log(map[string]interface{}{
			"title": title, "writer": writer, "type": typeName,
			"lastUpdate": lastUpdate, "lastChapter": lastChapter,
			"message": "extract html fail", "stage": "update",
		})
	}
	return title, writer, typeName, lastUpdate, lastChapter, err
}

func (book *Book) Update() bool {
	title, writer, typeName, lastUpdate, lastChapter, err := book.extractInfo()
	if err != nil {
		return false
	}
	// check difference
	update := false
	if lastUpdate != book.LastUpdate || lastChapter != book.LastChapter {
		update = true
	}
	if title != book.Title || writer != book.Writer || typeName != book.Type {
		update = true
		if book.DownloadFlag {
			book.Log(map[string]interface{}{
				"old": map[string]interface{}{
					"title": book.Title, "writer": book.Writer, "type": book.Type,
				},
				"new": map[string]interface{}{
					"title": title, "writer": writer, "type": typeName,
				},
				"message": "already download", "stage": "update",
			})
		}
		book.Version++
		book.EndFlag, book.DownloadFlag, book.ReadFlag = false, false, false
	}
	if update {
		// sync with online info
		book.Title, book.Writer, book.Type = title, writer, typeName
		book.LastUpdate, book.LastChapter = lastUpdate, lastChapter
	}
	return update
}
