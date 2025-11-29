package service

import (
	"context"
	"database/sql"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
)

type PatchStorageResult struct {
	FileExistWithNonDownloadStatus int64
	FileMissingWithDownloadStatus  int64
}

// pull latest data from vendor and update the local storage
//
//go:generate go tool mockgen -destination=../../mock/service/v1/book_service.go -package=mockservice . BookService
type BookService interface {
	// CreateBook(context.Context, *model.Book) error
	UpdateBook(context.Context, *model.Book) error
	DownloadBook(context.Context, *model.Book) error
	ProcessBook(context.Context, *model.Book) error // do all create / update + download (optional)
}

//go:generate go tool mockgen -destination=../../mock/service/v1/read_data_service.go -package=mockservice . ReadDataService
type ReadDataService interface {
	Book(ctx context.Context, site, id, hash string) (*model.Book, error)
	BookContent(context.Context, *model.Book) (string, error)
	BookChapters(context.Context, *model.Book) (model.Chapters, error)
	BookGroup(ctx context.Context, site, id, hash string) (*model.Book, *model.BookGroup, error)
	SearchBooks(ctx context.Context, title, writer string, limit, offset int) ([]model.Book, error)
	RandomBooks(ctx context.Context, limit int) ([]model.Book, error)

	Stats(context.Context, string) repo.Summary
	DBStats(context.Context) sql.DBStats
}
