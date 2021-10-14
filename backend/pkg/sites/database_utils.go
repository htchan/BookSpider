package sites

import (
	"strings"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

func (site *Site) OpenDatabase() {
	var err error
	site.database, err = sql.Open("sqlite3", site.databaseLocation)
	utils.CheckError(err)
	site.database.SetMaxIdleConns(10)
	site.database.SetMaxOpenConns(99999)
}
func (site *Site) CloseDatabase() {
	err := site.database.Close()
	utils.CheckError(err)
	site.database = nil
}

func (site Site) query(sql string, args ...interface{}) (*sql.Rows, *sql.Tx) {
	tx, err := site.database.Begin()
	utils.CheckError(err)
	rows, err := tx.Query(sql, args...)
	utils.CheckError(err)
	return rows, tx
}
func closeQuery(rows *sql.Rows, tx *sql.Tx) {
	utils.CheckError(rows.Close())
	tx.Rollback()
}

func (site Site) bookQuery(sql string, args ...interface{}) (*sql.Rows, error) {
	versionField := "version"
	if strings.Contains(sql, "group by") {
		versionField = "max(version)"
	}
	return site.bookLoadTx.Query("select site, num, "+versionField+", name, writer, type, "+
		"date, chapter, end, download, read from books"+sql, args...)
}

func (site Site) bookCount() (int, int) {
	distinctCount, count := -1, -1
	rows, tx := site.query("select count(DISTINCT num), count(*) as c from books")
	if rows.Next() {
		rows.Scan(&distinctCount, &count)
	}
	closeQuery(rows, tx)
	return distinctCount, count
}

func (site Site) errorCount() (int, int) {
	distinctCount, count := -1, -1
	rows, tx := site.query("select count(DISTINCT num), count(*) as c from error")
	for rows.Next() {
		rows.Scan(&distinctCount, &count)
	}
	closeQuery(rows, tx)
	return distinctCount, count
}

func (site Site) downloadCount() (int, int) {
	distinctCount, count := -1, -1
	rows, tx := site.query("select count(DISTINCT num), count(*) from books where download=?", true)
	if rows.Next() {
		rows.Scan(&distinctCount, &count)
	}
	closeQuery(rows, tx)
	return distinctCount, count
}

func (site Site) endCount() (int, int) {
	distinctCount, count := -1, -1
	rows, tx := site.query("select count(DISTINCT num), count(*) from books where end=?", true)
	if rows.Next() {
		rows.Scan(&distinctCount, &count)
	}
	closeQuery(rows, tx)
	return distinctCount, count
}

func (site Site) readCount() (int, int) {
	distinctCount, count := -1, -1
	rows, tx := site.query("select count(DISTINCT num), count(*) from books where read=?", true)
	if rows.Next() {
		rows.Scan(&distinctCount, &count)
	}
	closeQuery(rows, tx)
	return distinctCount, count
}

func (site Site) maxBookId() int {
	id := -1
	rows, tx := site.query("select num from books order by num desc")
	if rows.Next() {
		rows.Scan(&id)
	}
	closeQuery(rows, tx)
	return id
}

func (site Site) maxErrorId() int {
	id := -1
	rows, tx := site.query("select num from error order by num desc")
	if rows.Next() {
		rows.Scan(&id)
	}
	closeQuery(rows, tx)
	return id
}

func (site Site) maxId() int {
	maxErrorId := site.maxErrorId()
	maxBookId := site.maxBookId()
	if maxBookId > maxErrorId {
		return maxBookId
	} else {
		return maxErrorId
	}
}

func (site Site) missingIds() []int {
	maxBookId := site.maxBookId()
	missingIds := make([]int, 0)
	rows, tx := site.query("select num from " +
		"(select num from error union select num from books) order by num")
	var bookId int
	currentId := 1
	for ; rows.Next(); currentId++ {
		rows.Scan(&bookId)
		for ; currentId < bookId; currentId++ {
			missingIds = append(missingIds, currentId)
		}
	}
	for ; currentId < maxBookId; currentId++ {
		missingIds = append(missingIds, currentId)
	}
	closeQuery(rows, tx)
	return missingIds
}

func insertBookStmt(tx *sql.Tx) *sql.Stmt {
	stmt, err := tx.Prepare("insert into books " +
		"(site, num, version, name, writer, type, date, chapter, end, download, read)" +
		" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	utils.CheckError(err)
	return stmt
}

func updateBookStmt(tx *sql.Tx) *sql.Stmt {
	stmt, err := tx.Prepare("update books " +
		"set name=?, writer=?, type=?, date=?, chapter=?, end=?, download=?, read=?" +
		"where site=? and num=? and version=?")
	utils.CheckError(err)
	return stmt
}

func insertErrorStmt(tx *sql.Tx) *sql.Stmt {
	stmt, err := tx.Prepare("insert into error (site, num) values (?, ?)")
	utils.CheckError(err)
	return stmt
}

func deleteErrorStmt(tx *sql.Tx) *sql.Stmt {
	stmt, err := tx.Prepare("delete from error where site=? and num=?")
	utils.CheckError(err)
	return stmt
}

func (site *Site) PrepareStmt() {
	var err error
	site.bookOperateTx, err = site.database.Begin()
	utils.CheckError(err)
	site.insertBookStmt = insertBookStmt(site.bookOperateTx)
	site.updateBookStmt = updateBookStmt(site.bookOperateTx)
	site.insertErrorStmt = insertErrorStmt(site.bookOperateTx)
	site.deleteErrorStmt = deleteErrorStmt(site.bookOperateTx)
}

func (site *Site) CloseStmt() {
	utils.CheckError(site.insertBookStmt.Close())
	utils.CheckError(site.updateBookStmt.Close())
	utils.CheckError(site.insertErrorStmt.Close())
	utils.CheckError(site.deleteErrorStmt.Close())
	utils.CheckError(site.bookOperateTx.Commit())
	site.bookOperateTx = nil
	site.insertBookStmt = nil
	site.updateBookStmt = nil
	site.insertErrorStmt = nil
	site.deleteErrorStmt = nil
}

func (site Site) InsertBook(book books.Book) error {
	_, err := site.insertBookStmt.Exec(book.SiteName, book.Id, book.Version,
		book.Title, book.Writer, book.Type, book.LastUpdate, book.LastChapter,
		book.EndFlag, book.DownloadFlag, book.ReadFlag)
	return err
}

func (site Site) UpdateBook(book books.Book) error {
	_, err := site.updateBookStmt.Exec(book.Title, book.Writer, book.Type,
		book.LastUpdate, book.LastChapter, book.EndFlag, book.DownloadFlag, book.ReadFlag,
		book.SiteName, book.Id, book.Version)
	return err
}

func (site Site) InsertError(book books.Book) error {
	_, err := site.insertErrorStmt.Exec(book.SiteName, book.Id)
	return err
}

func (site Site) DeleteError(book books.Book) error {
	_, err := site.deleteErrorStmt.Exec(book.SiteName, book.Id)
	return err
}
