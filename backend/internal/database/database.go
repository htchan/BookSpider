package database

import (
	"time"
	"fmt"
	"strconv"
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

type SummaryRecord struct {
	BookCount, ErrorCount, WriterCount, UniqueBookCount int
	MaxBookId int
	LatestSuccessId int `json:"latestSuccessBookId"`
	StatusCount map[StatusCode]int
}

type Record interface{}

type Rows interface {
	Scan() (Record, error)
	ScanCurrent() (Record, error)
	Next() bool
	Close() error
}

type BookRecord struct {
	Site string
	Id, HashCode int
	Title string
	WriterId int
	Type string
	UpdateDate string
	UpdateChapter string
	Status StatusCode
}

func (record *BookRecord) Parameters() (parameters []interface{}) {
	parameters = make([]interface{}, 9)
	parameters[0] = record.Site
	parameters[1] = record.Id
	parameters[2] = record.HashCode
	parameters[3] = record.Title
	parameters[4] = record.WriterId
	parameters[5] = record.Type
	parameters[6] = record.UpdateDate
	parameters[7] = record.UpdateChapter
	parameters[8] = record.Status
	return
}
func (record *BookRecord) String() string {
	return fmt.Sprintf(
		"%v-%v-%v",
		record.Site,
		strconv.Itoa(record.Id),
		strconv.FormatInt(int64(record.HashCode), 36))
}

type WriterRecord struct {
	Id int
	Name string
}

func (record *WriterRecord) Parameters() (parameters []interface{}) {
	parameters = make([]interface{}, 2)
	parameters[0] = record.Id
	parameters[1] = record.Name
	return
}

type ErrorRecord struct {
	Site string
	Id int
	Error error
}

func (record *ErrorRecord) Parameters() (parameters []interface{}) {
	parameters = make([]interface{}, 3)
	parameters[0] = record.Site
	parameters[1] = record.Id
	parameters[2] = record.Error.Error()
	return
}

func GenerateHash() int {
	return int(time.Now().UnixMilli())
}

func StatustoString(status StatusCode) string {
	statusList := []string{ "error", "in_progress", "end", "download" }
	return statusList[status]
}