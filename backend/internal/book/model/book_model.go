package model

import (
	"database/sql"
	"fmt"
	"strings"
	"errors"
	"time"
)

type StatusCode int

const (
	Error = iota
	InProgress
	End
	Download
)
var StatusCodeMap = map[string]StatusCode{
	"ERROR": Error,
	"INPROGRESS": InProgress,
	"END": End,
	"DOWNLOAD": Download,
}

func StatusToString(status StatusCode) string {
	statusList := []string{ "error", "in_progress", "end", "download" }
	return statusList[status]
}

type BookModel struct {
	Site string
	ID, HashCode int
	Title string
	WriterID int
	Type string
	UpdateDate string
	UpdateChapter string
	Status StatusCode
}

func GenerateHash() int {
	return int(time.Now().UnixMilli())
}

func rowToBookModel(rows *sql.Rows) (BookModel, error) {
	var model BookModel
	err := rows.Scan(
		&model.Site, &model.ID, &model.HashCode,
		&model.Title, &model.WriterID, &model.Type,
		&model.UpdateDate, &model.UpdateChapter, &model.Status,
	)
	return model, err
}

func SaveBookModel(db *sql.DB, model *BookModel) error {
	upsertStatementsSequence := []string{
		`INSERT INTO books (
			site, id, hash_code,
			title, writer_id, type,
			update_date, update_chapter, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		`INSERT INTO books (
			site, id, hash_code,
			title, writer_id, type,
			update_date, update_chapter, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		`UPDATE books SET
			title=$4, writer_id=$5, type=$6,
			update_date=$7, update_chapter=$8, status=$9
		WHERE site=$1 and id=$2 and hash_code=$3;`,
	}
	var err error
	for i, statement := range upsertStatementsSequence {
		hashCode := model.HashCode
		if i == 0 {
			hashCode = 0
		}
		_, err = db.Exec(
			statement,
			model.Site, model.ID, hashCode,
			model.Title, model.WriterID, model.Type,
			model.UpdateDate, model.UpdateChapter, model.Status,
		)
		if err == nil {
			if i == 0 {
				model.HashCode = 0
			}
			break
		}
	}
	return err
}

func QueryBookModel(db *sql.DB, site string, id int, hashCode int) (BookModel, error) {
	queryStatement := `select
	site, id, hash_code,
	title, writer_id, type,
	update_date, update_chapter, status
	from books where site=$1 and id=$2`
	var (
		rows *sql.Rows
		err error
	)
	if hashCode < 0 {
		rows, err = db.Query(queryStatement + " order by hash_code desc limit 1", site, id)
	} else {
		rows, err = db.Query(queryStatement + " and hash_code=$3", site, id, hashCode)
	}
	if err != nil {
		return BookModel{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return BookModel{}, errors.New("book model not exist")
	}
	return rowToBookModel(rows)
}

func prepareQueryTitleWriterStatement(titlesCount, writersCount int) string {
	queryStatement := `select 
	site, id, hash_code, 
	title, writer_id, type, 
	update_date, update_chapter, status 
	from books 
	where site=$1`
	conditions := make([]string, titlesCount + writersCount)
	for i := 0; i < titlesCount; i++ {
		conditions[i] = fmt.Sprintf("(title like $%v)", i + 2)
	}
	for i := titlesCount; i < titlesCount + writersCount; i++ {
		conditions[i] = fmt.Sprintf("(writer_id=$%v)", i + 2)
	}
	if len(conditions) == 0 {
		return queryStatement + " and status=3 ORDER BY RANDOM() limit $2 offset $3;"
	}
	return fmt.Sprintf(
		"%s and (%s) order by id, hash_code limit $%v offset $%v;",
		queryStatement, strings.Join(conditions, " or "),
		titlesCount + writersCount + 2, titlesCount + writersCount + 3,
	)
}

func prepareQueryTitleWriterArgument(site string, titles []string, writerIDs []int, limit, offset int) []interface{} {
	result := make([]interface{}, 0, len(titles) + len(writerIDs) + 3)
	result = append(result, site)
	for _, title := range titles {
		result = append(result, title)
	}
	for _, writerID := range writerIDs {
		result = append(result, writerID)
	} 
	result = append(result, limit, offset)
	return result
}

func QueryBooksModelsByTitleWriter(
	db *sql.DB, site string, titles []string, writerIDs []int, limit int, offset int,
) ([]BookModel, error) {
	queryStatement := prepareQueryTitleWriterStatement(len(titles), len(writerIDs))
	for i := range titles {
		titles[i] = fmt.Sprintf("%%%s%%", titles[i])
	}
	rows, err := db.Query(queryStatement, prepareQueryTitleWriterArgument(site, titles, writerIDs, limit, offset)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var models []BookModel
	var finalErr error
	for rows.Next() {
		model, err := rowToBookModel(rows)
		if err != nil {
			finalErr = err
		}
		models = append(models, model)
	}
	return models, finalErr
}

func QueryAllBookModels(db *sql.DB) <-chan BookModel {
	queryStatement := "select site, id, hash_code, title, writer_id, type, update_date, update_chapter, status from books"
	bookChan := make(chan BookModel)
	rows, err := db.Query(queryStatement)
	if err != nil {
		close(bookChan)
		return bookChan
	}
	go func() {
		defer close(bookChan)
		defer rows.Close()
		var (
			bookModel BookModel
			err error
		)
		for rows.Next() {
			bookModel, err = rowToBookModel(rows)
			if err == nil {
				bookChan <- bookModel
			}
		}
	} ()
	return bookChan
}

func OrderBookModelsForUpdate(db *sql.DB, site string) ([]BookModel, error) {
	queryStatement := `select distinct on (site, id)
	site, id, hash_code,
	title, writer_id, type,
	update_date, update_chapter, status
	from books
	where site=$1
	order by site, id, update_date desc`

	rows, err := db.Query(queryStatement, site)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var models []BookModel
	var finalErr error
	for rows.Next() {
		model, err := rowToBookModel(rows)
		if err != nil {
			finalErr = err
		}
		models = append(models, model)
	}
	return models, finalErr
}