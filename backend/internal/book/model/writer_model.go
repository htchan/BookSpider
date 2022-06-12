package model

import (
	"database/sql"
	"errors"
)

type WriterModel struct {
	ID int
	Name string
}

func rowToWriterModel(rows *sql.Rows) (WriterModel, error) {
	var model WriterModel
	err := rows.Scan(&model.ID, &model.Name)
	return model, err
}

func SaveWriterModel(db *sql.DB, model *WriterModel) error {
	insertStatement := `insert into writers
	(name) values ($1) on conflict(name) do nothing;`
	_, err := db.Exec(insertStatement, model.Name)
	if err != nil {
		return err
	}
	queryStatement := `select id, name from writers where name=$1`
	rows, err := db.Query(queryStatement, model.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return errors.New("faile to fetch id")
	}
	*model, err = rowToWriterModel(rows)
	return err
}

func QueryWriterModel(db *sql.DB, id int) (WriterModel, error) {
	queryStatement := `select id, name from writers where id=$1;`
	rows, err := db.Query(queryStatement, id)
	if err != nil {
		return WriterModel{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return WriterModel{}, errors.New("writer model not exist")
	}
	return rowToWriterModel(rows)
}

func QueryAllWriterModels(db *sql.DB) <-chan WriterModel {
	queryStatement := "select id, name from writers"
	writerChan := make(chan WriterModel)
	rows, err := db.Query(queryStatement)
	if err != nil {
		close(writerChan)
		return writerChan
	}
	go func() {
		defer close(writerChan)
		defer rows.Close()
		var (
			writerModel WriterModel
			err error
		)
		for rows.Next() {
			writerModel, err = rowToWriterModel(rows)
			if err == nil {
				writerChan <- writerModel
			}
		}
	} ()
	return writerChan
}