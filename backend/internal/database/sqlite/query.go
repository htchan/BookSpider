package sqlite

import (
	"strings"
	"html"
	"github.com/htchan/BookSpider/internal/database"
)

func (db *SqliteDB) QueryBookBySiteIdHash(site string, id int, hashCode int) database.Rows {
	rows := new(SqliteBookRows)
	if hashCode >= 0 {
		rows._rows, _ = db._db.Query(
			"select " + database.BOOK_RECORD_FIELDS + " from books " +
			"where site=? and id=? and hash_code=?",
			html.EscapeString(site), id, hashCode)
	} else {
		rows._rows, _ = db._db.Query(
			"select " + database.BOOK_RECORD_FIELDS + " from books " +
			"where site=? and id=? order by hash_code desc",
			html.EscapeString(site), id)
	}
	return rows
}

func (db *SqliteDB) QueryBooksByPartialTitleAndWriter(titles []string, writers []int) database.Rows {
	rows := new(SqliteBookRows)
	statement := ""
	if len(titles) > 0 {
		statement += "title like ?" + strings.Repeat(" or title like ?", len(titles) - 1)
	}
	if len(writers) > 0 {
		if len(statement) > 0 {
			statement += " or "
		}
		statement += "writer_id in (?" + strings.Repeat(", ?", len(writers) - 1) + ")"
	}
	if statement != "" {
		statement = "select " + database.BOOK_RECORD_FIELDS + " from books " +
			"where " + statement
	} else {
		rows._rows = nil
		return rows
	}
	arguments := make([]interface{}, len(titles) + len(writers))
	for i, title := range titles {
		arguments[i] = "%" + html.EscapeString(title) + "%"
	}
	for i, writer := range writers {
		arguments[len(titles) + i] = writer
	}
	rows._rows, _ = db._db.Query(statement, arguments...)
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
		html.EscapeString(name))
	return rows
}

func (db *SqliteDB) QueryWritersByPartialName(names []string) database.Rows {
	rows := new(SqliteWriterRows)
	if len(names) == 0 {
		rows._rows = nil
		return rows
	}
	arguments := make([]interface{}, len(names))
	for i := range names {
		arguments[i] = "%" + html.EscapeString(names[i]) + "%"
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
		html.EscapeString(site), id)
	return rows
}

func (db *SqliteDB) QueryBooksOrderByUpdateDate() database.Rows {
	rows := new(SqliteBookRows)
	rows._rows, _ = db._db.Query(
		"select " + database.BOOK_RECORD_FIELDS + " from books " +
		"where status != ? group by site, id " +
		"order by max(update_date) desc", database.Error)
	return rows
}

func (db *SqliteDB) QueryBooksWithRandomOrder(n int, status database.StatusCode) database.Rows {
	rows := new(SqliteBookRows)
	if status == database.Error {
		rows._rows, _ = db._db.Query(
			"select " + database.BOOK_RECORD_FIELDS + " from books " +
			"order by random() limit ?", n)
	} else {
		rows._rows, _ = db._db.Query(
			"select " + database.BOOK_RECORD_FIELDS + " from books " +
			"where status=?" +
			"order by random() limit ?", status, n)
	}
	return rows
}