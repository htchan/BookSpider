package service

import (
	"context"
	"database/sql"
	"sync/atomic"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
)

type BookOperation func(context.Context, *model.Book) error

type SiteOperation func(context.Context) error

type UpdateStats struct {
	Total             atomic.Int64
	Fail              atomic.Int64
	Unchanged         atomic.Int64
	NewChapter        atomic.Int64
	NewEntity         atomic.Int64
	ErrorUpdated      atomic.Int64
	InProgressUpdated atomic.Int64
	EndUpdated        atomic.Int64
	DownloadedUpdated atomic.Int64
}

type DownloadStats struct {
	Total               atomic.Int64
	Success             atomic.Int64
	NoChapter           atomic.Int64
	TooManyFailChapters atomic.Int64
	RequestFail         atomic.Int64
}

type PatchStorageStats struct {
	FileExist   atomic.Int64
	FileMissing atomic.Int64
}

//go:generate go tool mockgen -destination=../mock/service/v1/service.go -package=mockservice . Service
type Service interface {
	Name() string
	// Backup() error
	PatchDownloadStatus(context.Context, *PatchStorageStats) error
	PatchMissingRecords(context.Context, *UpdateStats) error
	CheckAvailability(context.Context) error

	UpdateBook(context.Context, *model.Book, *UpdateStats) error
	Update(context.Context, *UpdateStats) error

	ExploreBook(context.Context, *model.Book, *UpdateStats) error
	Explore(context.Context, *UpdateStats) error

	DownloadBook(context.Context, *model.Book, *DownloadStats) error
	Download(context.Context, *DownloadStats) error

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

//go:generate go tool mockgen -destination=../mock/service/v1/read_data_service.go -package=mockservice . ReadDataService
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
