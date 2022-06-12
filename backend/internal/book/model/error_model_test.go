package model

import (
	"testing"
	"errors"
	"github.com/google/go-cmp/cmp"
	"runtime"
)

func TestErrorModel_SaveErrorModel(t *testing.T) {
	t.Parallel()
	site := "save_err_model"

	t.Cleanup(func () {
		db.Exec("delete from errors where site=$1", site)
		runtime.GC()
	})

	t.Run("create error if model.Err is not nil", func (t *testing.T) {
		t.Parallel()

		model := ErrorModel{Site: site, ID: 1, Err: errors.New("some error")}
		err := SaveErrorModel(db, &model)
		if err != nil {
			t.Errorf("save error model return error: %v", err)
			return
		}
		result, err := QueryErrorModel(db, model.Site, model.ID)
		if err != nil || !cmp.Equal(result.Error(), "some error") {
			t.Errorf("query saved error model return error: %v; result: %v", err, result)
		}
	})

	t.Run("update error if error record already exist", func (t *testing.T) {
		t.Parallel()
		
		model := ErrorModel{Site: site, ID: 2, Err: errors.New("some error")}
		SaveErrorModel(db, &model)
		model.Err = errors.New("another error")
		err := SaveErrorModel(db, &model)
		if err != nil {
			t.Errorf("save error model return error: %v", err)
			return
		}
		result, err := QueryErrorModel(db, model.Site, model.ID)
		if err != nil || !cmp.Equal(result.Error(), "another error") {
			t.Errorf("query saved error model return error: %v; result: %v", err, result)
		}
	})

	t.Run("delete error if model.Err is nil", func (t *testing.T) {
		t.Parallel()
		
		model := ErrorModel{Site: site, ID: 3, Err: errors.New("some error")}
		SaveErrorModel(db, &model)
		model.Err = nil
		err := SaveErrorModel(db, &model)
		if err != nil {
			t.Errorf("save error model return error: %v", err)
			return
		}
		result, err := QueryErrorModel(db, model.Site, model.ID)
		if err == nil {
			t.Errorf("query saved error model return error: %v; result: %v", err, result)
		}
	})
}

func TestErrorModel_QueryErrorModel(t *testing.T) {
	t.Parallel()
	site := "query_err_model"

	t.Cleanup(func () {
		db.Exec("delete from errors where site=$1", site)
		runtime.GC()
	})

	db.Exec(`insert into errors
	(site, id, data)
	values
	($1, 1, 'data_1');`,
	site)

	t.Run("query existing error model", func (t *testing.T) {
		result, err := QueryErrorModel(db, site, 1)
		if err != nil || result.Site != site || result.ID != 1 || result.Error() != "data_1" {
			t.Errorf("query model return wrong result: error: %v; result: %v", err, result)
		}
	})

	t.Run("query non existence error model", func (t *testing.T) {
		result, err := QueryErrorModel(db, site, -123)
		if err == nil {
			t.Errorf("query model return wrong result: error: %v; result: %v", err, result)
		}
	})
}