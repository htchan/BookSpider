package repo

import (
	"database/sql"
	"errors"

	"github.com/htchan/BookSpider/internal/model"
)

var ChapterEndKeywords = []string{
	"后记", "後記", "新书", "新書", "结局", "結局", "感言",
	"尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本" /*"结束", "結束", */, "完結",
	"完结", "终结", "終結", "番外", "结尾", "結尾", "全书完", "全書完", "全本完",
}

var BookNotExist = errors.New("no records found")

type Repostory interface {
	Migrate() error
	// book related
	CreateBook(*model.Book) error
	UpdateBook(*model.Book) error
	// the system will not delete exiting books

	FindBookById(id int) (*model.Book, error) // return book with the largest hash code
	FindBookByIdHash(id, hash int) (*model.Book, error)
	FindBooksByStatus(status model.StatusCode) (<-chan model.Book, error)
	FindAllBooks() (<-chan model.Book, error)
	FindBooksForUpdate() (<-chan model.Book, error)
	FindBooksForDownload() (<-chan model.Book, error)
	FindBooksByTitleWriter(title, writer string, limit, offset int) ([]model.Book, error)
	FindBooksByRandom(limit int) ([]model.Book, error)
	UpdateBooksStatus() error

	// writer related
	SaveWriter(*model.Writer) error // create and update id in writer
	// the system will not delete / update existing writers

	// error related
	SaveError(*model.Book, error) error // create / update / delete errors depends on error content

	// database
	Backup(path string) error
	DBStats() sql.DBStats // return empty if repo is not based on db
	Stats() Summary

	Close() error
}
