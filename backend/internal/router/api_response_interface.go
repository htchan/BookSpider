package router

import (
	"database/sql"

	"github.com/htchan/BookSpider/internal/model"
)

type errResp struct {
	Error string `json:"error"`
}

type booksResp struct {
	Books []model.Book `json:"books"`
}

type dbStatsResp struct {
	Stats []sql.DBStats `json:"stats"`
}
