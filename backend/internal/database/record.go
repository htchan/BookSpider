package database

import (
	"fmt"
	"strconv"
)

type SummaryRecord struct {
	BookCount, ErrorCount, WriterCount, UniqueBookCount int
	MaxBookId int
	LatestSuccessId int `json:"latestSuccessBookId"`
	StatusCount map[StatusCode]int
}

type Record interface{}

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

func (record BookRecord)Equal(compare BookRecord) bool {
	return record.Site == compare.Site && record.Id == compare.Id &&
		record.HashCode == compare.HashCode && record.Title == compare.Title &&
		record.WriterId == compare.WriterId && record.Type == compare.Type &&
		record.UpdateDate == compare.UpdateDate &&
		record.UpdateChapter == compare.UpdateChapter &&
		record.Status == compare.Status
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

func (record WriterRecord)Equal(compare WriterRecord) bool {
	return record.Id == compare.Id && record.Name == compare.Name
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

func (record ErrorRecord)Equal(compare ErrorRecord) bool {
	return record.Site == compare.Site && record.Id == compare.Id &&
		record.Error.Error() == compare.Error.Error()
}