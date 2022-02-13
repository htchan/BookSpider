package sites

import (
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"strings"
	"strconv"
)

func (site *Site)rowsToRecords(rows database.Rows) (records []database.Record) {
	defer rows.Close()
	records = make([]database.Record, 0)
	for rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)
		records = append(records, record)
	}
	return
}

func (site *Site)SearchByIdHash(id int, hash string) (book *books.Book) {
	defer utils.Recover(func() {})
	err := site.OpenDatabase()
	utils.CheckError(err)
	defer site.database.Close()

	hashCode, err := strconv.ParseInt(hash, 36, 64)
	if err != nil {
		hashCode = -1
	}

	book = books.LoadBook(site.database, site.Name, id, int(hashCode), site.config.BookMeta)
	return
}

func (site *Site)SearchByWriterId(writerId int) (bookArray []*books.Book) {
	// defer utils.Recover(func() {})
	err := site.OpenDatabase()
	utils.CheckError(err)
	defer site.database.Close()
	
	rows := site.database.QueryBooksByWriterId(writerId)
	bookArray = make([]*books.Book, 0)

	for _, record := range site.rowsToRecords(rows) {
		bookArray = append(
			bookArray,
			books.LoadBookByRecord(
				site.database, record.(*database.BookRecord), site.config.BookMeta))
	}
	return
}

func (site *Site)SearchByStatus(status database.StatusCode) (bookArray []*books.Book) {
	defer utils.Recover(func() {})
	err := site.OpenDatabase()
	utils.CheckError(err)
	defer site.database.Close()
	
	rows := site.database.QueryBooksByStatus(status)
	bookArray = make([]*books.Book, 0)

	for _, record := range site.rowsToRecords(rows) {
		bookArray = append(
			bookArray,
			books.LoadBookByRecord(
				site.database, record.(*database.BookRecord), site.config.BookMeta))
	}
	return
}

func (site *Site)SearchByTitleWriter(titleSearch string, writerSearch string) (bookArray []*books.Book) {
	defer utils.Recover(func() {})
	err := site.OpenDatabase()
	utils.CheckError(err)
	defer site.database.Close()
	var titles, writerNames []string
	if titleSearch != "" {
		titles = strings.Split(titleSearch, " ")
	}
	if writerSearch != "" {
		writerNames = strings.Split(writerSearch, " ")
	}
	rows := site.database.QueryWritersByPartialName(writerNames)
	writerIds := make([]int, 0)

	for rows.Next() {
		record, _ := rows.ScanCurrent()
		writerIds = append(writerIds, record.(*database.WriterRecord).Id)
	}
	rows.Close()

	rows = site.database.QueryBooksByPartialTitleAndWriter(titles, writerIds)
	bookArray = make([]*books.Book, 0)

	for _, record := range site.rowsToRecords(rows) {
		bookArray = append(
			bookArray,
			books.LoadBookByRecord(
				site.database, record.(*database.BookRecord), site.config.BookMeta))
	}
	return
}

func (site *Site) RandomSuggestBook(n int, status database.StatusCode) (bookArray []*books.Book) {defer utils.Recover(func() {})
	err := site.OpenDatabase()
	utils.CheckError(err)
	defer site.database.Close()
	rows := site.database.QueryBooksWithRandomOrder(n, status)
	bookArray = make([]*books.Book, 0)

	for _, record := range site.rowsToRecords(rows) {
		bookArray = append(bookArray, books.LoadBookByRecord(
			site.database, record.(*database.BookRecord), site.config.BookMeta))
	}
	return
}