package book

import (
	"errors"
	"fmt"
	"log"

	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
)

func baseURL(bk model.Book, config config.BookConfig) string {
	return fmt.Sprintf(config.URLConfig.Base, bk.ID)
}

func fetchInfo(c *client.CircuitBreakerClient, bookKey, url string) (title, writer, typeStr, date, chapStr string, err error) {
	html, err := c.Get(url)
	if err != nil {
		return
	}
	responseApi := ApiParser.Parse(bookKey+".info", html)
	okMap := make(map[string]bool)
	title, okMap["title"] = responseApi.Data["Title"]
	writer, okMap["writer"] = responseApi.Data["Writer"]
	typeStr, okMap["type"] = responseApi.Data["Type"]
	date, okMap["date"] = responseApi.Data["LastUpdate"]
	chapStr, okMap["chapter"] = responseApi.Data["LastChapter"]
	for _, key := range []string{"title", "writer", "type", "date", "chapter"} {
		if !okMap[key] {
			err = errors.New(fmt.Sprintf("%v not found", key))
			return
		}
	}
	return
}

func isNewBook(bk model.Book, title, writer string) bool {
	return bk.Status != model.Error && (title != bk.Title || writer != bk.Writer.Name)
}

func isUpdated(bk model.Book, title, writer, typeStr, date, chapStr string) bool {
	return title != bk.Title || writer != bk.Writer.Name ||
		typeStr != bk.Type || date != bk.UpdateDate ||
		chapStr != bk.UpdateChapter
}

func Update(bk *model.Book, bkConf config.BookConfig, stConf config.SiteConfig, c *client.CircuitBreakerClient) (bool, error) {
	title, writer, typeStr, date, chapStr, err := fetchInfo(c, stConf.BookKey, baseURL(*bk, bkConf))
	// TODO: log the response
	if err != nil {
		if bk.Status == model.Error {
			bk.Error = err
		}
		return false, err
	}

	if isNewBook(*bk, title, writer) {
		bk.HashCode = model.GenerateHash()
		log.Printf("[%v] new book found: title: %v", bk, bk.Title)
	}

	isUpdated := isUpdated(*bk, title, writer, typeStr, date, chapStr)
	if isUpdated {
		// TODO: log uipdated
		bk.Status = model.InProgress
		bk.Error = nil
		// populate updated fields
		bk.Title = title
		bk.Writer.Name = writer
		bk.Type = typeStr
		bk.UpdateDate = date
		bk.UpdateChapter = chapStr
		log.Printf("[%v] updated book found: title: %v", bk, bk.Title)
	} else {
		// TODO: log not updated, should I?
	}
	return isUpdated, nil
}
