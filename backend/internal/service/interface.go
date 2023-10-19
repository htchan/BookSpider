package service

import (
	"context"
	"database/sql"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
)

type BookOperation func(context.Context, *model.Book) error

type SiteOperation func(context.Context) error

type Service interface {
	Name() string
	// Backup() error
	PatchDownloadStatus(context.Context) error
	PatchMissingRecords(context.Context) error
	CheckAvailability(context.Context) error

	UpdateBook(context.Context, *model.Book) error
	Update(context.Context) error

	ExploreBook(context.Context, *model.Book) error
	Explore(context.Context) error

	DownloadBook(context.Context, *model.Book) error
	Download(context.Context) error

	ValidateBookEnd(context.Context, *model.Book) error
	ValidateEnd(context.Context) error

	ProcessBook(context.Context, *model.Book) error
	Process(context.Context) error

	BookInfo(context.Context, *model.Book) string
	BookContent(context.Context, *model.Book) (string, error)
	BookChapters(context.Context, *model.Book) (model.Chapters, error)
	Book(ctx context.Context, id, hash string) (*model.Book, error)
	BookGroup(ctx context.Context, id, hash string) (*model.Book, *model.BookGroup, error)
	QueryBooks(ctx context.Context, title, writer string, limit, offset int) ([]model.Book, error)
	RandomBooks(ctx context.Context, limit int) ([]model.Book, error)

	Stats(context.Context) repo.Summary
	DBStats(context.Context) sql.DBStats
}
