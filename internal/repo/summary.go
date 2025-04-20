package repo

import "github.com/htchan/BookSpider/internal/model"

type Summary struct {
	BookCount       int
	WriterCount     int
	ErrorCount      int
	UniqueBookCount int
	MaxBookID       int
	LatestSuccessID int
	DownloadCount   int
	StatusCount     map[model.StatusCode]int
}
