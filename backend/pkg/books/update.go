package books

import (
	"errors"
	"github.com/htchan/BookSpider/internal/utils"
)

func (book *Book) extractInfo() (string, string, string, string, string, error) {
	// get online resource, try maximum 10 times if it keeps failed
	html, trial := utils.GetWeb(book.metaInfo.baseUrl, 10, book.decoder, book.CONST_SLEEP)
	if _, err := utils.Search(html, book.metaInfo.titleRegex); err != nil || !book.validHTML(html, book.metaInfo.baseUrl, trial) {
		book.Log(map[string]interface{}{
			"url": book.metaInfo.baseUrl, "error": "extract base html fail",
			"html": html, "trial": trial, "stage": "update",
		})
		return "", "", "", "", "", errors.New("invalid base page html")
	}
	// extract info from source
	var err, tempErr error
	var result [5]string
	regex := [5]string {
		book.metaInfo.titleRegex, book.metaInfo.writerRegex, book.metaInfo.typeRegex,
		book.metaInfo.lastUpdateRegex, book.metaInfo.lastChapterRegex,
	}
	for i := 0; i < 5; i++ {
		result[i], tempErr = utils.Search(html, regex[i])
		if tempErr != nil {
			err = tempErr
		}
	}
	if err != nil {
		book.Log(map[string]interface{}{
			"title": result[0], "writer": result[1], "type": result[2],
			"lastUpdate": result[3], "lastChapter": result[4],
			"message": "extract html fail", "error": err.Error(), "stage": "update",
		})
	}
	return result[0], result[1], result[2], result[3], result[4], err
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
