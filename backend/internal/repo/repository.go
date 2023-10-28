package repo

import (
	"database/sql"

	"github.com/htchan/BookSpider/internal/model"
)

//go:generate mockgen -destination=../mock/repo/repository.go -package=mockrepo . Repository
type Repository interface {
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

	FindBookGroupByID(id int) (model.BookGroup, error)
	FindBookGroupByIDHash(id, hashCode int) (model.BookGroup, error)

	FindAllBookIDs() ([]int, error)

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
