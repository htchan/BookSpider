package model

import (
	"testing"
	"github.com/google/go-cmp/cmp"
	"runtime"
)

func TestWriterModel_SaveWriterModel(t *testing.T) {
	t.Parallel()

	t.Cleanup(func () {
		db.Exec("delete from writers where id > 0")
		runtime.GC()
	})

	t.Run("generate id when save new writer model", func(t *testing.T) {
		t.Parallel()
		model := WriterModel{Name: "writer name"}
		err := SaveWriterModel(db, &model)
		if err != nil {
			t.Errorf("save writer model return error: %v", err)
			return
		}
		if model.ID == 0 {
			t.Errorf("save writer model not update writer model: %v", model)
			return
		}
		result, err := QueryWriterModel(db, model.ID)
		if err != nil || !cmp.Equal(result, model) {
			t.Errorf("query saved model return: error: %v; result: %v", err, result)
		}
	})

	t.Run("do nothing if writer already exist", func(t *testing.T) {
		t.Parallel()
		model := WriterModel{Name: "writer name"}
		SaveWriterModel(db, &model)
		err := SaveWriterModel(db, &model)
		if err != nil {
			t.Errorf("save writer model return error: %v", err)
			return
		}
		if model.ID == 0 {
			t.Errorf("save writer model not update writer model: %v", model)
			return
		}
		result, err := QueryWriterModel(db, model.ID)
		if err != nil || !cmp.Equal(result, model) {
			t.Errorf("query saved model return: error: %v; result: %v", err, result)
		}
	})
}

func TestWriterModel_QueryWriterModel(t *testing.T) {
	t.Parallel()
	t.Cleanup(func () {
		db.Exec("delete from writers where id=0")
		runtime.GC()
	})

	db.Exec(`insert into writers (id, name) values (0, 'writer_1')`)

	t.Run("query existing writer", func (t *testing.T) {
		t.Parallel()
		result, err := QueryWriterModel(db, 0)
		if err != nil || !cmp.Equal(result.Name, "writer_1") {
			t.Errorf("query model return wrong result: error: %v; result: %v", err, result)
		}
	})

	t.Run("query non existence writer", func (t *testing.T) {
		t.Parallel()
		result, err := QueryWriterModel(db, -123)
		if err == nil {
			t.Errorf("query model return wrong result: error: %v; result: %v", err, result)
		}
	})
}

func TestWriterModel_QueryAllWriterModels(t *testing.T) {
	t.Parallel()
	t.Skip()
}