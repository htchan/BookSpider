package repo

import (
	"context"
	"database/sql"

	"github.com/htchan/BookSpider/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func GetTracer() trace.Tracer {
	return otel.Tracer("htchan/BookSpider/repository")
}

//go:generate go tool mockgen -destination=../mock/repo/repository.go -package=mockrepo . Repository
type Repository interface {
	// book related
	CreateBook(context.Context, *model.Book) error
	UpdateBook(context.Context, *model.Book) error
	// the system will not delete exiting books

	FindBookById(ctx context.Context, site string, id int) (*model.Book, error) // return book with the largest hash code
	FindBookByIdHash(ctx context.Context, site string, id, hash int) (*model.Book, error)
	FindBooksByStatus(ctx context.Context, status model.StatusCode) (<-chan model.Book, error)
	FindAllBooks(ctx context.Context, site string) (<-chan model.Book, error)
	FindBooksForUpdate(ctx context.Context, site string) (<-chan model.Book, error)
	FindBooksForDownload(ctx context.Context, site string) (<-chan model.Book, error)
	FindBooksByTitleWriter(ctx context.Context, title, writer string, limit, offset int) ([]model.Book, error)
	FindBooksByRandom(ctx context.Context, limit int) ([]model.Book, error)
	UpdateBooksStatus(context.Context) error

	FindBookGroupByID(ctx context.Context, site string, id int) (model.BookGroup, error)
	FindBookGroupByIDHash(ctx context.Context, site string, id, hashCode int) (model.BookGroup, error)

	FindAllBookIDs(ctx context.Context, site string) ([]int, error)

	// writer related
	SaveWriter(context.Context, *model.Writer) error // create and update id in writer
	// the system will not delete / update existing writers

	// error related
	SaveError(context.Context, *model.Book, error) error // create / update / delete errors depends on error content

	// database
	Backup(ctx context.Context, site, path string) error
	DBStats(context.Context) sql.DBStats // return empty if repo is not based on db
	Stats(ctx context.Context, site string) Summary

	Close() error
}
