package sqlite

import (
	"strings"
	"github.com/htchan/BookSpider/internal/database"
)

func (db *SqliteDB) QueryBookBySiteIdHash(site string, id int, hashCode int) database.Rows {
	rows := new(SqliteBookRows)
	if hashCode > 0 {
		rows._rows, _ = db._db.Query(
			"select " + database.BOOK_RECORD_FIELDS + " from books " +
			"where site=? and id=? and hash_code=?",
			site, id, hashCode)
	} else {
		rows._rows, _ = db._db.Query(
			"select " + database.BOOK_RECORD_FIELDS + " from books " +
			"where site=? and id=? order by hash_code desc",
			site, id)
	}
	return rows
}

func (db *SqliteDB) QueryBooksByTitle(title string) database.Rows {
	rows := new(SqliteBookRows)
	rows._rows, _ = db._db.Query(
		"select " + database.BOOK_RECORD_FIELDS + " from books " +
		"where title=? group by site, id",
		title)
	return rows
}

func (db *SqliteDB) QueryBooksByPartialTitle(titles []string) database.Rows {
	rows := new(SqliteBookRows)
	arguments := make([]interface{}, len(titles))
	for i := range titles {
		arguments[i] = "%" + titles[i] + "%"
	}
	rows._rows, _ = db._db.Query(
		"select " + database.BOOK_RECORD_FIELDS + " from books " +
		"where title like ?" + strings.Repeat(" or title like ?", len(titles) - 1),
		arguments...)
	return rows
}

func (db *SqliteDB) QueryBooksByWriterId(writerId int) database.Rows {
	rows := new(SqliteBookRows)
	rows._rows, _ = db._db.Query(
		"select " + database.BOOK_RECORD_FIELDS + " from books " +
		"where writer_id=? group by site, id",
		writerId)
	return rows
}

func (db *SqliteDB) QueryBooksByStatus(status database.StatusCode) database.Rows {
	rows := new(SqliteBookRows)
	rows._rows, _ = db._db.Query(
		"select " + database.BOOK_RECORD_FIELDS + " from books " +
		"where status=? group by site, id",
		status)
	return rows
}

func (db *SqliteDB) QueryWriterById(id int) database.Rows {
	rows := new(SqliteWriterRows)
	rows._rows, _ = db._db.Query(
		"select " + database.WRITER_RECORD_FIELDS + " from writers " +
		"where id=?",
		id)
	return rows
}

func (db *SqliteDB) QueryWriterByName(name string)  database.Rows {
	rows := new(SqliteWriterRows)
	rows._rows, _ = db._db.Query(
		"select " + database.WRITER_RECORD_FIELDS + " from writers " +
		"where name=?",
		name)
	return rows
}

func (db *SqliteDB) QueryWritersByPartialName(names []string) database.Rows {
	rows := new(SqliteWriterRows)
	arguments := make([]interface{}, len(names))
	for i := range names {
		arguments[i] = "%" + names[i] + "%"
	}
	rows._rows, _ = db._db.Query(
		"select " + database.WRITER_RECORD_FIELDS + " from writers " +
		"where name like ?" + strings.Repeat(" or name like ?", len(names) - 1),
		arguments...)
	return rows
}

func (db *SqliteDB) QueryErrorBySiteId(site string, id int) database.Rows {
	rows := new(SqliteErrorRows)
	rows._rows, _ = db._db.Query(
		"select " + database.ERROR_RECORD_FIELDS + " from errors " +
		"where site=? and id=?",
		site, id)
	return rows
}
