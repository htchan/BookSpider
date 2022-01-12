package database

const (
	BOOK_RECORD_FIELDS = "site, id, hash_code, title, writer_id, type, " +
		"update_date, update_chapter, status "
	WRITER_RECORD_FIELDS = "id, name"
	ERROR_RECORD_FIELDS = "site, id, data"
)

type DB interface {
	QueryBookBySiteIdHash(string, int, int) Rows
	QueryBooksByTitle(string) Rows
	QueryBooksByPartialTitle([]string) Rows
	QueryBooksByWriterId(int) Rows
	QueryBooksByStatus(StatusCode) Rows
	QueryWriterById(int) Rows
	QueryWriterByName(string) Rows
	QueryWritersByPartialName([]string) Rows
	QueryErrorBySiteId(string, int) Rows

	CreateBookRecord(*BookRecord) error
	CreateWriterRecord(*WriterRecord) error
	CreateErrorRecord(*ErrorRecord) error

	UpdateBookRecord(*BookRecord) error

	DeleteBookRecord([]BookRecord) error
	DeleteWriterRecord([]WriterRecord) error
	DeleteErrorRecord([]ErrorRecord) error

	Summary(string) SummaryRecord

	Close() error
}

type StatusCode int

const (
	Error = iota
	InProgress
	End
	Download
)

type SummaryRecord struct {
	BookCount, ErrorCount, WriterCount, UniqueBookCount int
	MaxBookId, LatestSuccessId int
	StatusCount map[StatusCode]int
}

type Record interface{}

type Rows interface {
	Scan() (Record, error)
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

type WriterRecord struct {
	Id int
	Name string
}

type ErrorRecord struct {
	Site string
	Id int
	Error error
}