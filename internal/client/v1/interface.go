package client

import (
	"context"
	"time"
)

type BookInfo struct {
	Title         string
	Author        string
	Type          string
	UpdateChapter string
	UpdateDate    time.Time
}

type ChapterEntry struct {
	Title string
	URL   string
}

type ChapterEntryList []ChapterEntry

type ChapterContent struct {
	Title string
	Body  string
}

type Client interface {
	GetBookInfo(ctx context.Context, bookID string) (*BookInfo, error)
	GetBookChapterList(ctx context.Context, bookID string) (ChapterEntryList, error)
	GetChapterContent(ctx context.Context, chapter ChapterEntry) (*ChapterContent, error)
	Available(ctx context.Context) bool
}
