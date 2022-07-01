package model

import (
	"database/sql"
	"errors"
)

type ErrorModel struct {
	Site string
	ID   int
	Err  error
}

func rowToErrorModel(rows *sql.Rows) (ErrorModel, error) {
	var (
		model  ErrorModel
		errStr string
	)
	err := rows.Scan(&model.Site, &model.ID, &errStr)
	if err == nil {
		model.Err = errors.New(errStr)
	}
	return model, err
}

func SaveErrorModel(db *sql.DB, model *ErrorModel) error {
	if model.Err == nil {
		deleteStatement := `delete from errors where site=$1 and id=$2;`
		_, err := db.Exec(deleteStatement, model.Site, model.ID)
		return err
	}
	upsertStatementsSequence := []string{
		`insert into errors (site, id, data) values ($1, $2, $3);`,
		`update errors set data=$3 where site=$1 and id=$2;`,
	}
	var err error
	for _, statement := range upsertStatementsSequence {
		_, err = db.Exec(statement, model.Site, model.ID, model.Error())
		if err == nil {
			break
		}
	}
	return err
}

func QueryErrorModel(db *sql.DB, site string, id int) (ErrorModel, error) {
	queryStatement := `select site, id, data from errors where site=$1 and id=$2;`
	rows, err := db.Query(queryStatement, site, id)
	if err != nil {
		return ErrorModel{Site: site, ID: id}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return ErrorModel{Site: site, ID: id}, errors.New("error model not exist")
	}
	return rowToErrorModel(rows)
}

func (model ErrorModel) Error() string {
	if model.Err == nil {
		return ""
	}
	return model.Err.Error()
}
