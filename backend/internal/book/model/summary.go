package model

import (
	"database/sql"
)

type SummaryResult struct {
	BookCount       int
	WriterCount     int
	ErrorCount      int
	UniqueBookCount int
	MaxBookID       int
	LatestSuccessID int
	StatusCount     map[StatusCode]int
}

func Summary(db *sql.DB, name string) SummaryResult {
	var (
		summary SummaryResult
	)

	rowsBook, err := db.Query("select count(*), count(distinct (site||id)), max(id) from books where site=$1", name)
	if err == nil && rowsBook.Next() {
		defer rowsBook.Close()
		rowsBook.Scan(&summary.BookCount, &summary.UniqueBookCount, &summary.MaxBookID)
	}

	rowsLatestSuccessID, err := db.Query("select max(id) from books where status<>$1 and site=$2", Error, name)
	if err == nil && rowsLatestSuccessID.Next() {
		defer rowsLatestSuccessID.Close()
		rowsLatestSuccessID.Scan(&summary.LatestSuccessID)
	}

	rowsWriter, err := db.Query("select count(distinct writers.id) from books join writers on books.writer_id=writers.id where site=$1", name)
	if err == nil && rowsWriter.Next() {
		defer rowsWriter.Close()
		rowsWriter.Scan(&summary.WriterCount)
	}

	rowsError, err := db.Query("select count(distinct id) from errors where site=$1", name)
	if err == nil && rowsError.Next() {
		defer rowsError.Close()
		rowsError.Scan(&summary.ErrorCount)
	}

	summary.StatusCount = make(map[StatusCode]int)
	rowsStatus, err := db.Query("select status, count(*) from books where site=$1 group by status", name)
	var (
		statusKey   StatusCode
		statusValue int
	)
	if err == nil {
		defer rowsStatus.Close()
	}
	for rowsStatus.Next() {
		rowsStatus.Scan(&statusKey, &statusValue)
		summary.StatusCount[statusKey] = statusValue
	}

	return summary
}
