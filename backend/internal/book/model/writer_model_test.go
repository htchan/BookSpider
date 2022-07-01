package model

import (
	"github.com/google/go-cmp/cmp"
	"runtime"
	"testing"
)

func TestWriterModel_SaveWriterModel(t *testing.T) {
	t.Parallel()

	writerName := "writer name"

	t.Cleanup(func() {
		db.Exec("delete from writers where name=$1", writerName)
		runtime.GC()
	})

	t.Run("generate id when save new writer model", func(t *testing.T) {
		t.Parallel()
		model := WriterModel{Name: writerName}
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
		model := WriterModel{Name: writerName}
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

	writerName := "writer_1"

	t.Cleanup(func() {
		db.Exec("delete from writers where name=$1", writerName)
		runtime.GC()
	})

	db.Exec(`insert into writers (id, name) values (-1, $1)`, writerName)

	t.Run("query existing writer", func(t *testing.T) {
		t.Parallel()
		result, err := QueryWriterModel(db, -1)
		if err != nil || !cmp.Equal(result.Name, writerName) {
			t.Errorf("query model return wrong result: error: %v; result: %v", err, result)
		}
	})

	t.Run("query non existence writer", func(t *testing.T) {
		t.Parallel()
		result, err := QueryWriterModel(db, -123)
		if err == nil {
			t.Errorf("query model return wrong result: error: %v; result: %v", err, result)
		}
	})
}

func TestWriterModel_QueryAllWriterModels(t *testing.T) {
	SaveWriterModel(db, &WriterModel{Name: "all_writer_1"})
	SaveWriterModel(db, &WriterModel{Name: "all_writer_2"})

	t.Cleanup(func() {
		db.Exec("delete from writers where name like $1", "all_writer%")
	})

	t.Run("works", func(t *testing.T) {
		writerModelsChan, err := QueryAllWriterModels(db)
		if err != nil {
			t.Errorf("query all writer models return error: %v", err)
		}
		i := 0
		for m := range writerModelsChan {
			t.Log(m)
			i++
		}
		if i != 3 {
			t.Errorf("query all writer models return %d models", i)
		}
	})
}
