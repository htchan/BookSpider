package database

import (
	"time"
)

const (
	BOOK_RECORD_FIELDS = "site, id, hash_code, title, writer_id, type, " +
		"update_date, update_chapter, status"
	WRITER_RECORD_FIELDS = "id, name"
	ERROR_RECORD_FIELDS = "site, id, data"
)

type DB interface {
	QueryBookBySiteIdHash(string, int, int) Rows
	QueryBooksByPartialTitleAndWriter(titles []string, writers []int) Rows
	QueryBooksByWriterId(int) Rows
	QueryBooksByStatus(StatusCode) Rows
	// QueryBooksPotentialDuplicated() Rows
	QueryWriterById(int) Rows
	QueryWriterByName(string) Rows
	QueryWritersByPartialName([]string) Rows
	QueryErrorBySiteId(string, int) Rows

	QueryBooksOrderByUpdateDate() Rows
	QueryBooksWithRandomOrder(n int, status StatusCode) Rows

	CreateBookRecord(*BookRecord, *WriterRecord) error
	CreateWriterRecord(*WriterRecord) error
	CreateErrorRecord(*ErrorRecord) error

	UpdateBookRecord(*BookRecord, *WriterRecord) error
	UpdateErrorRecord(*ErrorRecord) error

	UpdateBookRecordsStatusByChapter() error

	DeleteBookRecords([]BookRecord) error
	DeleteWriterRecords([]WriterRecord) error
	DeleteErrorRecords([]ErrorRecord) error

	Summary(siteName string) SummaryRecord
	Backup(directory, filename string) error

	Commit() error
	Close() error
}

type StatusCode int

const (
	Error = iota
	InProgress
	End
	Download
)
var StatusCodeMap = map[string]StatusCode{ "ERROR": Error, "INPROGRESS": InProgress, "END": End, "DOWNLOAD": Download }


type Rows interface {
	Scan() (Record, error)
	ScanCurrent() (Record, error)
	Next() bool
	Close() error
}

func GenerateHash() int {
	return int(time.Now().UnixMilli())
}

func StatustoString(status StatusCode) string {
	statusList := []string{ "error", "in_progress", "end", "download" }
	return statusList[status]
}