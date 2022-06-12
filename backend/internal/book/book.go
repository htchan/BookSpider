package book

import (
	"errors"
	"fmt"
	"strconv"
	"database/sql"
	"encoding/json"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/client"
)

type Book struct {
	model.BookModel
	model.WriterModel
	model.ErrorModel
	*config.BookConfig
	*client.CircuitBreakerClient
}

func NewBook(site string, id int, con *config.BookConfig, client *client.CircuitBreakerClient) Book {
	return Book{
		BookModel: model.BookModel{Site: site, ID: id, HashCode: model.GenerateHash()},
		ErrorModel: model.ErrorModel{Site: site, ID: id, Err: errors.New("new book")},
		BookConfig: con,
		CircuitBreakerClient: client,
	}
}

func LoadBook(db *sql.DB, site string, id, hash int, config *config.BookConfig, client *client.CircuitBreakerClient) (Book, error) {
	book := Book{BookConfig: config, CircuitBreakerClient: client}
	bookModel, err := model.QueryBookModel(db, site, id, hash)
	if err != nil {
		return book, err
	}

	return LoadBookFromModel(db, bookModel, config, client)
}

func LoadBookFromModel(db *sql.DB, bookModel model.BookModel, config *config.BookConfig, client *client.CircuitBreakerClient) (Book, error) {
	var (
		book = Book{
			BookModel: bookModel,
			BookConfig: config,
			CircuitBreakerClient: client,
		}
		err error
	)

	site, id, writerID := bookModel.Site, bookModel.ID, bookModel.WriterID

	book.ErrorModel, err = model.QueryErrorModel(db, site, id)
	if err != nil && err.Error() != "error model not exist" {
		return book, err
	}

	book.WriterModel, err = model.QueryWriterModel(db, writerID)
	return book, err
}

func (book *Book) Save(db *sql.DB) error {
	var err error
	err = model.SaveErrorModel(db, &book.ErrorModel)
	if err != nil {
		return err
	}
	err = model.SaveWriterModel(db, &book.WriterModel)
	if err != nil {
		return err
	}
	book.BookModel.WriterID = book.WriterModel.ID
	err = model.SaveBookModel(db, &book.BookModel)
	return err
}

func (book Book) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Site string `json:"site"`
		Id int `json:"id"`
		Hash string `json:"hash"`
		Title string `json:"title"`
		Writer string `json:"writer"`
		Type string `json:"type"`
		UpdateDate string `json:"update_date"`
		UpdateChapter string `json:"update_chapter"`
		Status string `json:"status"`
	} {
		Site: book.BookModel.Site,
		Id: book.BookModel.ID,
		Hash: strconv.FormatInt(int64(book.BookModel.HashCode), 36),
		Title: book.BookModel.Title,
		Writer: book.WriterModel.Name,
		Type: book.BookModel.Type,
		UpdateDate: book.BookModel.UpdateDate,
		UpdateChapter: book.BookModel.UpdateChapter,
		Status: model.StatusToString(book.BookModel.Status),
	})
}

func (book Book) String() string {
	return fmt.Sprintf(
		"bk.%v.%v.%v",
		book.BookModel.Site,
		strconv.Itoa(book.BookModel.ID),
		strconv.FormatInt(int64(book.BookModel.HashCode), 36),
	)
}