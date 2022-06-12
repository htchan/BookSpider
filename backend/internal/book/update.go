package book

import (
	"fmt"
	"errors"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/ApiParser"
)

func (book Book) baseURL() string {
	return fmt.Sprintf(book.BookConfig.URL.Base, book.BookModel.ID)
}

func (book Book) fetchInfo() (title, writer, typeStr, date, chapterStr string, err error) {
	html, err := book.Get(book.baseURL())
	if err != nil {
		return
	}
	responseApi := ApiParser.Parse(book.BookConfig.SourceKey + ".info", html)
	okMap := make(map[string]bool)
	title, okMap["title"] = responseApi.Data["Title"]
	writer, okMap["writer"] = responseApi.Data["Writer"]
	typeStr, okMap["type"] = responseApi.Data["Type"]
	date, okMap["date"] = responseApi.Data["LastUpdate"]
	chapterStr, okMap["chapter"] = responseApi.Data["LastChapter"]
	for key := range okMap {
		if !okMap[key] {
			err = errors.New(fmt.Sprintf("%v not found", key))
			return
		}
	}
	return
}

func (book Book) isNewBook(title, writer, typeStr string) bool {
	return title != book.BookModel.Title || writer != book.WriterModel.Name ||
	typeStr != book.BookModel.Type
}

func (book Book) isUpdated(title, writer, typeStr, date, chapterStr string) bool {
	return title != book.BookModel.Title || writer != book.WriterModel.Name ||
	typeStr != book.BookModel.Type || date != book.BookModel.UpdateDate ||
	chapterStr != book.BookModel.UpdateChapter
}

func (book *Book) Update() (bool, error) {
	title, writer, typeStr, date, chapterStr, err := book.fetchInfo()
	if err != nil {
		if book.BookModel.Status == model.Error {
			book.ErrorModel.Err = err
		}
		return false, err
	}
	if book.isNewBook(title, writer, typeStr) && book.Status != model.Error {
		book.BookModel.HashCode = model.GenerateHash()
	}
	isUpdated := book.isUpdated(title, writer, typeStr, date, chapterStr)
	if isUpdated {
		book.BookModel.Status = model.InProgress
		book.ErrorModel.Err = nil
		// populate the updated fields
		book.BookModel.Title = title
		book.WriterModel.Name = writer
		book.BookModel.Type = typeStr
		book.BookModel.UpdateDate = date
		book.BookModel.UpdateChapter = chapterStr
	}
	return isUpdated, err
}