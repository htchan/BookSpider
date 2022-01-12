package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/htchan/BookSpider/internal/database"
	"errors"
)

type SqliteBookRows struct {
	_rows *sql.Rows
}

type SqliteWriterRows struct {
	_rows *sql.Rows
}

type SqliteErrorRows struct {
	_rows *sql.Rows
}

func (rows SqliteBookRows) Scan() (database.Record, error) {
	if rows._rows == nil {
		return nil, errors.New("empty rows")
	}
	record := new(database.BookRecord)
	err := rows._rows.Scan(
		&record.Site, &record.Id, &record.HashCode,
		&record.Title, &record.WriterId, &record.Type,
		&record.UpdateDate, &record.UpdateChapter,
		&record.Status)
	return record, err
}

func (rows *SqliteBookRows) Next() bool {
	if rows._rows == nil {
		return false
	}
	result := rows._rows.Next()
	if result == false {
		rows._rows = nil
	}
	return result
}

func (rows *SqliteBookRows) Close() error {
	if rows._rows == nil {
		return errors.New("error rows closed")
	}
	err := rows._rows.Close()
	rows._rows = nil
	return err
}

func (rows *SqliteWriterRows) Scan() (database.Record, error) {
	if rows._rows == nil {
		return nil, errors.New("empty rows")
	}
	record := new(database.WriterRecord)
	err := rows._rows.Scan(&record.Id, &record.Name)
	return record, err
}

func (rows *SqliteWriterRows) Next() bool {
	if rows._rows == nil {
		return false
	}
	result := rows._rows.Next()
	if result == false {
		rows._rows = nil
	}
	return result
}

func (rows *SqliteWriterRows) Close() error {
	if rows._rows == nil {
		return errors.New("error rows closed")
	}
	err := rows._rows.Close()
	rows._rows = nil
	return err
}

func (rows *SqliteErrorRows) Scan() (database.Record, error) {
	if rows._rows == nil {
		return nil, errors.New("empty rows")
	}
	record := new(database.ErrorRecord)
	var errorString string
	err := rows._rows.Scan(&record.Site, &record.Id, &errorString)
	record.Error = errors.New(errorString)
	return record, err
}

func (rows *SqliteErrorRows) Next() bool {
	if rows._rows == nil {
		return false
	}
	result := rows._rows.Next()
	if result == false {
		rows._rows = nil
	}
	return result
}

func (rows *SqliteErrorRows) Close() error {
	if rows._rows == nil {
		return errors.New("error rows closed")
	}
	err := rows._rows.Close()
	rows._rows = nil
	return err
}