package books

import (
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
	"strconv"
	"io/ioutil"
	"errors"
	"os"
	"fmt"
)

type Book struct {
    bookRecord *database.BookRecord
    writerRecord *database.WriterRecord
    errorRecord *database.ErrorRecord
    config configs.BookConfig
}

func NewBook(site string, id int, hash int, config *configs.BookConfig) (book *Book) {
	if hash == -1 {
		hash = database.GenerateHash()
	}
	return &Book{
		config: config.Populate(id),
		bookRecord: &database.BookRecord{
			Site: site,
			Id: id,
			HashCode: hash,
		},
		writerRecord: &database.WriterRecord{
			Id: 0,
			Name: "",
		},
	}
}

func LoadBook(db database.DB, site string, id int, hash int, config *configs.BookConfig) (book *Book) {
	defer utils.Recover(func() { book = nil })
	book = new(Book)
	book.config = config.Populate(id)

	query := db.QueryBookBySiteIdHash(site, id, hash)
	record, err := query.Scan()
	utils.CheckError(err)
	book.bookRecord = record.(*database.BookRecord)
	query.Close()

	query = db.QueryWriterById(book.bookRecord.WriterId)
	record, err = query.Scan()
	utils.CheckError(err)
	book.writerRecord = record.(*database.WriterRecord)
	query.Close()
	
	query = db.QueryErrorBySiteId(site, id)
	if query.Next() {
		record, err = query.ScanCurrent()
		utils.CheckError(err)
		book.errorRecord = record.(*database.ErrorRecord)
	}
	query.Close()
	return
}

func LoadBookByRecord(db database.DB, bookRecord *database.BookRecord, config *configs.BookConfig) (book *Book) {
	defer utils.Recover(func() { book = nil })
	book = new(Book)
	book.config = config.Populate(bookRecord.Id)

	book.bookRecord = bookRecord

	query := db.QueryWriterById(book.bookRecord.WriterId)
	record, err := query.Scan()
	utils.CheckError(err)
	book.writerRecord = record.(*database.WriterRecord)
	query.Close()
	
	query = db.QueryErrorBySiteId(bookRecord.Site, bookRecord.Id)
	if query.Next() {
		record, err = query.ScanCurrent()
		utils.CheckError(err)
		book.errorRecord = record.(*database.ErrorRecord)
	}
	query.Close()
	return
}

func (book *Book)saveWriterRecord(db database.DB) {
	if book.writerRecord.Id < 0 {
		query := db.QueryWriterByName(book.writerRecord.Name)
		defer query.Close()
		record, err := query.Scan()
		if err != nil {
			utils.CheckError(db.CreateWriterRecord(book.writerRecord))
			book.bookRecord.WriterId = book.writerRecord.Id
		} else {
			book.writerRecord.Id = record.(*database.WriterRecord).Id
			book.bookRecord.WriterId = record.(*database.WriterRecord).Id
		}
	}
}

func (book *Book)saveBookRecord(db database.DB) {
	query := db.QueryBookBySiteIdHash(
		book.bookRecord.Site,
		book.bookRecord.Id,
		book.bookRecord.HashCode)
	exist := query.Next()
	query.Close()
	if exist {
		utils.CheckError(db.UpdateBookRecord(book.bookRecord, book.writerRecord))
	} else {
		utils.CheckError(db.CreateBookRecord(book.bookRecord, book.writerRecord))
	}
}

func (book *Book)saveErrorRecord(db database.DB) {
	query := db.QueryErrorBySiteId(
		book.bookRecord.Site, 
		book.bookRecord.Id)
	errorRecord, err := query.Scan()
	query.Close()
	
	if book.bookRecord.Status != database.Error && err == nil {
		utils.CheckError(
			db.DeleteErrorRecords(
				[]database.ErrorRecord {*errorRecord.(*database.ErrorRecord) } ))
	} else if book.bookRecord.Status == database.Error && err == nil {
		utils.CheckError(db.UpdateErrorRecord(book.errorRecord))
	} else if book.bookRecord.Status == database.Error && err != nil {
		utils.CheckError(db.CreateErrorRecord(book.errorRecord))
	}
}

func (book *Book)Save(db database.DB) (result bool) {
	result = true
	// defer utils.Recover(func() { result = false })
	book.saveWriterRecord(db)
	book.saveBookRecord(db)
	book.saveErrorRecord(db)
	return
}

func (book *Book)GetInfo() (string, int, int) {
	return book.bookRecord.Site, book.bookRecord.Id, book.bookRecord.HashCode
}

func (book *Book)GetTitle() string { return book.bookRecord.Title }
func (book *Book)SetTitle(title string) { book.bookRecord.Title = title }

func (book *Book)GetWriter() string { return book.writerRecord.Name }
func (book *Book)SetWriter(name string) {
	book.bookRecord.WriterId = -1
	book.writerRecord.Id = -1
	book.writerRecord.Name = name
}

func (book *Book)GetType() string { return book.bookRecord.Type }
func (book *Book)SetType(typeString string) { book.bookRecord.Type = typeString }

func (book *Book)GetUpdateDate() string { return book.bookRecord.UpdateDate }
func (book *Book)SetUpdateDate(date string) { book.bookRecord.UpdateDate = date }

func (book *Book)GetUpdateChapter() string { return book.bookRecord.UpdateChapter }
func (book *Book)SetUpdateChapter(chapter string) { book.bookRecord.UpdateChapter = chapter }

func (book *Book)GetStatus() database.StatusCode { return book.bookRecord.Status }
func (book *Book)SetStatus(status database.StatusCode) { book.bookRecord.Status = status }

func (book *Book)GetError() error {
	if book.errorRecord == nil {
		return nil
	}
	return book.errorRecord.Error
}
func (book *Book)SetError(err error) {
	if err == nil {
		book.errorRecord = nil
		return
	}
	if book.errorRecord == nil {
		book.errorRecord = new(database.ErrorRecord)
		book.errorRecord.Site = book.bookRecord.Site
		book.errorRecord.Id = book.bookRecord.Id
	}
	book.errorRecord.Error = err
}

func (book *Book)getContentLocation() (location string) {
	location = os.Getenv("ASSETS_LOCATION") + book.config.StorageDirectory + strconv.Itoa(book.bookRecord.Id)
	if book.bookRecord.HashCode != 0 {
		location += "-v" + strconv.Itoa(book.bookRecord.HashCode)
	}
	location += ".txt"
	return
}
func (book *Book)HasContent() bool {
	return utils.Exists(book.getContentLocation())
}
func (book *Book)GetContent() (content string) {
	// defer utils.Recover(func() {})
	if !book.HasContent() {
		return 
	}
	contentBytes, err := ioutil.ReadFile(book.getContentLocation())
	utils.CheckError(err)
	content = string(contentBytes)
	return
}

func (book *Book)validHTML(html string) error {
	if len(html) == 0 {
		return errors.New("load html fail - zero length")
	} else if _, err := strconv.Atoi(html); err == nil {
		return errors.New("load html fail - code " + html)
	}
	return nil
}

func (book *Book)Map() map[string]interface{} {
	site, id, hash := book.GetInfo()
	return map[string]interface{} {
		"site": site,
		"id": id,
		"hash": strconv.FormatInt(int64(hash), 36),
		"title": book.GetTitle(),
		"writer": book.GetWriter(),
		"type": book.GetType(),
		"updateDate": book.GetUpdateDate(),
		"updateChapter": book.GetUpdateChapter(),
		"status": database.StatustoString(book.GetStatus()),
	}
}

func (book *Book)String() string {
	return fmt.Sprintf(
		"%v-%v-%v",
		book.bookRecord.Site,
		strconv.Itoa(book.bookRecord.Id),
		strconv.FormatInt(int64(book.bookRecord.HashCode), 36))
}
