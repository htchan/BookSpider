package model

import (
	"github.com/google/go-cmp/cmp"
	"runtime"
	"strings"
	"testing"
)

func TestBookModel_StatusToString(t *testing.T) {
	t.Parallel()

	t.Run("error status", func(t *testing.T) {
		t.Parallel()
		result := StatusToString(Error)
		if result != "error" {
			t.Errorf("status to string return %v", result)
		}
	})

	t.Run("in_progress status", func(t *testing.T) {
		t.Parallel()
		result := StatusToString(InProgress)
		if result != "in_progress" {
			t.Errorf("status to string return %v", result)
		}
	})

	t.Run("end status", func(t *testing.T) {
		t.Parallel()
		result := StatusToString(End)
		if result != "end" {
			t.Errorf("status to string return %v", result)
		}
	})

	t.Run("download status", func(t *testing.T) {
		t.Parallel()
		result := StatusToString(Download)
		if result != "download" {
			t.Errorf("status to string return %v", result)
		}
	})
}

func TestBookModel_SaveBookModel(t *testing.T) {
	t.Parallel()
	site := "save_bk_model"
	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	t.Run("create book with hash 0 if none of the book exist", func(t *testing.T) {
		t.Parallel()
		id, hashCode := 1, 100
		model := BookModel{Site: site, ID: id, HashCode: hashCode}
		err := SaveBookModel(db, &model)
		if err != nil {
			t.Errorf("save book model return error: %v", err)
			return
		}
		if model.HashCode != 0 {
			t.Errorf("save book model does bot update model: %v", model)
			return
		}
		result, err := QueryBookModel(db, model.Site, model.ID, model.HashCode)
		if err != nil || !cmp.Equal(model, result) {
			t.Errorf("saved wrong result: error: %v; model: %v", err, result)
		}
	})

	t.Run("create book with model hash if hash 0 exist", func(t *testing.T) {
		t.Parallel()
		id, hashCode := 2, 100
		model := BookModel{Site: site, ID: id, HashCode: hashCode}
		SaveBookModel(db, &model)
		model.HashCode = 100
		err := SaveBookModel(db, &model)
		if err != nil {
			t.Errorf("save book model return error: %v", err)
			return
		}
		result, err := QueryBookModel(db, model.Site, model.ID, model.HashCode)
		if err != nil || !cmp.Equal(model, result) {
			t.Errorf("saved wrong result: error: %v; model: %v", err, result)
		}
	})

	t.Run("update book with hash 0 if model hash is 0 and it exist", func(t *testing.T) {
		t.Parallel()
		id, hashCode := 3, 0
		model := BookModel{Site: site, ID: id, HashCode: hashCode}
		SaveBookModel(db, &model)
		model.WriterID, model.Status = 1, InProgress
		err := SaveBookModel(db, &model)
		if err != nil {
			t.Errorf("save book model return error: %v", err)
			return
		}
		result, err := QueryBookModel(db, site, id, hashCode)
		if err != nil || !cmp.Equal(model, result) {
			t.Errorf("saved wrong result: error: %v; model: %v", err, result)
		}
	})
}

func TestBookModel_QueryBookModel(t *testing.T) {
	t.Parallel()
	site, id := "query_bk_model", 1
	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	db.Exec(`insert into books 
	(site, id, hash_code, title, writer_id, type, update_date, update_chapter, status) 
	VALUES 
	($1, $2, 0, '', 0, '', '', '', 0), 
	($1, $2, 100, '', 0, '', '', '', 0);`,
		site, id)

	t.Run("query existing book with specific site, id, hash", func(t *testing.T) {
		t.Parallel()
		model := BookModel{Site: site, ID: id, HashCode: 0}
		result, err := QueryBookModel(db, site, id, 0)
		if err != nil || !cmp.Equal(model, result) {
			t.Errorf("query model return wrong result: error: %v; result:%v", err, result)
		}
	})

	t.Run("query latest book with specific site, id", func(t *testing.T) {
		t.Parallel()
		model := BookModel{Site: site, ID: id, HashCode: 100}
		result, err := QueryBookModel(db, site, id, -1)
		if err != nil || !cmp.Equal(model, result) {
			t.Errorf("query model return wrong result: error: %v; result:%v", err, result)
		}
	})

	t.Run("query non existence book", func(t *testing.T) {
		t.Parallel()
		result, err := QueryBookModel(db, site, -123, -1)
		if err == nil {
			t.Errorf("query model return wrong result: error: %v; result:%v", err, result)
		}
	})
}

func TestBookModel_QueryBookModelsByStatus(t *testing.T) {
	t.Parallel()

	site := "bks_model_state"
	model_1 := BookModel{Site: site, ID: 1, HashCode: 0, Status: Error}
	model_2 := BookModel{Site: site, ID: 2, HashCode: 0, Status: InProgress}
	model_3 := BookModel{Site: site, ID: 3, HashCode: 0, Status: Download}
	model_4 := BookModel{Site: site, ID: 4, HashCode: 0, Status: End}
	SaveBookModel(db, &model_1)
	SaveBookModel(db, &model_2)
	SaveBookModel(db, &model_3)
	SaveBookModel(db, &model_4)

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
	})

	t.Run("query error books", func(t *testing.T) {
		results, err := QueryBookModelsByStatus(db, site, Error)

		if err != nil || len(results) != 1 || results[0].Status != Error {
			t.Errorf("query book models by status return %v", results)
		}
	})

	t.Run("query in progress books", func(t *testing.T) {
		results, err := QueryBookModelsByStatus(db, site, InProgress)

		if err != nil || len(results) != 1 || results[0].Status != InProgress {
			t.Errorf("query book models by status return %v", results)
		}
	})

	t.Run("query end books", func(t *testing.T) {
		results, err := QueryBookModelsByStatus(db, site, End)

		if err != nil || len(results) != 1 || results[0].Status != End {
			t.Errorf("query book models by status return %v", results)
		}
	})

	t.Run("query download books", func(t *testing.T) {
		results, err := QueryBookModelsByStatus(db, site, Download)

		if err != nil || len(results) != 1 || results[0].Status != Download {
			t.Errorf("query book models by status return %v", results)
		}
	})
}

func TestBookModel_QueryBookModelsByRandom(t *testing.T) {
	t.Parallel()

	site := "random_model"

	SaveBookModel(db, &BookModel{Site: site, ID: 1, Status: Download})
	SaveBookModel(db, &BookModel{Site: site, ID: 2, Status: Download})

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		bkModels, err := QueryBookModelsByRandom(db, site, 10, 0)
		if err != nil || len(bkModels) != 2 {
			t.Errorf("query book models by random return error: %v", err)
			t.Errorf("query book models by random return models: %v", bkModels)
		}
	})

	t.Run("limit result", func(t *testing.T) {
		t.Parallel()
		bkModels, err := QueryBookModelsByRandom(db, site, 1, 0)
		if err != nil || len(bkModels) != 1 {
			t.Errorf("query book models by random return error: %v", err)
			t.Errorf("query book models by random return models: %v", bkModels)
		}
	})
}

func TestBookModel_prepareQueryTitleWriterStatement(t *testing.T) {
	t.Parallel()

	t.Run("not zero title count and writers count", func(t *testing.T) {
		t.Parallel()

		expect := `select 
		site, books.id, hash_code, 
		title, writer_id, type, 
		update_date, update_chapter, status 
		from books join writers on books.writer_id=writers.id 
		where site=$1 and ((title like $2) or (name like $3)) 
		order by books.id, hash_code limit $4 offset $5;`
		result := prepareQueryTitleWriterStatement(1, 1)
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\t", ""), "\n", "")
		result = strings.ReplaceAll(strings.ReplaceAll(result, "\t", ""), "\n", "")
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("zero title count and non zero writers count", func(t *testing.T) {
		t.Parallel()

		expect := `select 
		site, books.id, hash_code, 
		title, writer_id, type, 
		update_date, update_chapter, status 
		from books join writers on books.writer_id=writers.id 
		where site=$1 and ((name like $2)) 
		order by books.id, hash_code limit $3 offset $4;`
		result := prepareQueryTitleWriterStatement(0, 1)
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\t", ""), "\n", "")
		result = strings.ReplaceAll(strings.ReplaceAll(result, "\t", ""), "\n", "")
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("non zero title count and zero writers count", func(t *testing.T) {
		t.Parallel()

		expect := `select 
		site, books.id, hash_code, 
		title, writer_id, type, 
		update_date, update_chapter, status 
		from books join writers on books.writer_id=writers.id 
		where site=$1 and ((title like $2)) 
		order by books.id, hash_code limit $3 offset $4;`
		result := prepareQueryTitleWriterStatement(1, 0)
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\t", ""), "\n", "")
		result = strings.ReplaceAll(strings.ReplaceAll(result, "\t", ""), "\n", "")
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})
}

func TestBookModel_prepareQueryTitleWriterArguement(t *testing.T) {
	t.Parallel()

	t.Run("non empty title and writers", func(t *testing.T) {
		t.Parallel()

		titles := []string{"t1", "t2"}
		writers := []string{"w1", "w2"}
		expect := []interface{}{"site", "t1", "t2", "w1", "w2", 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("empty title and non empty writers", func(t *testing.T) {
		t.Parallel()

		titles := []string{}
		writers := []string{"w1", "w2"}
		expect := []interface{}{"site", "w1", "w2", 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("non empty title and empty writers", func(t *testing.T) {
		t.Parallel()

		titles := []string{"t1", "t2"}
		writers := []string{}
		expect := []interface{}{"site", "t1", "t2", 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("empty title and writers", func(t *testing.T) {
		t.Parallel()

		titles := []string{}
		writers := []string{}
		expect := []interface{}{"site", 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})
}

func TestBookModel_QueryBookModelsByTitleWriter(t *testing.T) {
	t.Parallel()
	site := "query_bks_model"
	writer_model_1 := WriterModel{Name: "hij"}
	writer_model_2 := WriterModel{Name: "jkl"}
	writer_model_3 := WriterModel{Name: "lmn"}
	SaveWriterModel(db, &writer_model_1)
	SaveWriterModel(db, &writer_model_2)
	SaveWriterModel(db, &writer_model_3)
	t.Log(writer_model_1, writer_model_2, writer_model_3)
	book_model_1 := BookModel{Site: site, ID: 1, HashCode: 0, Title: "abc", WriterID: writer_model_1.ID}
	book_model_2 := BookModel{Site: site, ID: 1, HashCode: 1, Title: "cde", WriterID: writer_model_2.ID}
	book_model_3 := BookModel{Site: site, ID: 1, HashCode: 2, Title: "efg", WriterID: writer_model_3.ID, Status: Download}
	SaveBookModel(db, &book_model_1)
	SaveBookModel(db, &book_model_2)
	SaveBookModel(db, &book_model_3)
	t.Log(book_model_1, book_model_2, book_model_3)

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where name=$1 or name=$2 or name=$3", "hij", "jkl", "lmn")
		runtime.GC()
	})

	t.Run("providing both titles and writer ids", func(t *testing.T) {
		t.Parallel()

		result, err := QueryBookModelsByTitleWriter(db, site, []string{"abc"}, []string{"jkl"}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{book_model_1, book_model_2}, result) {
			t.Errorf("query book models by title or writers return error: %v", err)
			t.Errorf("query book models by title or writers return err: %v", result)
		}
	})

	t.Run("providing titles only", func(t *testing.T) {
		t.Parallel()

		result, err := QueryBookModelsByTitleWriter(db, site, []string{"c"}, []string{}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{book_model_1, book_model_2}, result) {
			t.Errorf("query book models by title or writers return error: %v", err)
			t.Errorf("query book models by title or writers return err: %v", result)
		}
	})

	t.Run("providing writer ids only", func(t *testing.T) {
		t.Parallel()

		result, err := QueryBookModelsByTitleWriter(db, site, []string{}, []string{"j"}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{book_model_1, book_model_2}, result) {
			t.Errorf("query book models by title or writers return error: %v", err)
			t.Errorf("query book models by title or writers return err: %v", result)
		}
	})

	t.Run("not providing titles and writer ids", func(t *testing.T) {
		t.Parallel()

		result, err := QueryBookModelsByTitleWriter(db, site, []string{}, []string{}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{book_model_3}, result) {
			t.Errorf("query book models by title or writers return error: %v; result: %v", err, result)
		}
	})

	t.Run("offset all result", func(t *testing.T) {
		t.Parallel()

		result, err := QueryBookModelsByTitleWriter(db, site, []string{}, []string{}, 3, 3)
		if err != nil || len(result) != 0 {
			t.Errorf("query book models by title or writers return error: %v; result: %v", err, result)
		}
	})

	t.Run("limit zero result", func(t *testing.T) {
		t.Parallel()

		result, err := QueryBookModelsByTitleWriter(db, site, []string{}, []string{}, 0, 0)
		if err != nil || len(result) != 0 {
			t.Errorf("query book models by title or writers return error: %v", err)
		}
	})
}

func TestBookModel_QueryAllBookModels(t *testing.T) {
	site := "all_models"

	SaveBookModel(db, &BookModel{Site: site, ID: 1, Status: Download})
	SaveBookModel(db, &BookModel{Site: site, ID: 1, HashCode: 1, Status: InProgress})
	SaveBookModel(db, &BookModel{Site: site, ID: 2, Status: End})
	SaveBookModel(db, &BookModel{Site: site, ID: 3, Status: Error})

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
	})

	t.Run("works", func(t *testing.T) {
		bkModelsChan, err := QueryAllBookModels(db, site)
		if err != nil {
			t.Errorf("query all book models return error: %v", err)
		}
		i := 0
		for _ = range bkModelsChan {
			i++
		}
		if i != 4 {
			t.Errorf("query all book models return %d result", i)
		}
	})
}

func TestBookModel_OrderBookModelsForUpdate(t *testing.T) {
	t.Parallel()

	site := "order_bks_model"
	model_1 := BookModel{Site: site, ID: 1, HashCode: 0, UpdateDate: "3"}
	model_1_2 := BookModel{Site: site, ID: 1, HashCode: 2, UpdateDate: "4"}
	model_2 := BookModel{Site: site, ID: 2, HashCode: 0, UpdateDate: "2"}
	model_3 := BookModel{Site: site, ID: 3, HashCode: 0, UpdateDate: "1"}
	SaveBookModel(db, &model_1)
	SaveBookModel(db, &model_1_2)
	SaveBookModel(db, &model_2)
	SaveBookModel(db, &model_3)

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	t.Run("success", func(t *testing.T) {
		bkModelsChan, err := OrderBookModelsForUpdate(db, site)
		if err != nil {
			t.Errorf("query all book models return error: %v", err)
		}
		i := 0
		want := []BookModel{model_1_2, model_2, model_3}
		for book := range bkModelsChan {
			if !cmp.Equal(book, want[i]) {
				t.Errorf("OrderBookModelsForUpdate() at %d position return %v, want %v", i, book, want[i])
			}
			i++
		}
		if i != len(want) {
			t.Errorf("OrderBookModelsForUpdate() return %d result", i)
		}
	})
}

func TestBookModel_OrderBookModelsForDownload(t *testing.T) {
	t.Parallel()

	site := "down_bks_model"
	model_1 := BookModel{Site: site, ID: 1, HashCode: 0, Status: End, UpdateDate: "2"}
	model_1_2 := BookModel{Site: site, ID: 1, HashCode: 2, Status: End, UpdateDate: "1"}
	model_2 := BookModel{Site: site, ID: 2, HashCode: 0, Status: Download, UpdateDate: "3"}
	model_3 := BookModel{Site: site, ID: 3, HashCode: 0, Status: InProgress, UpdateDate: "4"}
	SaveBookModel(db, &model_1)
	SaveBookModel(db, &model_1_2)
	SaveBookModel(db, &model_2)
	SaveBookModel(db, &model_3)

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	t.Run("success", func(t *testing.T) {
		bkModelsChan, err := OrderBookModelsForDownload(db, site)
		if err != nil {
			t.Errorf("query all book models return error: %v", err)
		}
		i := 0
		want := []BookModel{model_1, model_1_2}
		for book := range bkModelsChan {
			if !cmp.Equal(book, want[i]) {
				t.Errorf("OrderBookModelsForDownload() at %d position return %v, want %v", i, book, want[i])
			}
			i++
		}
		if i != len(want) {
			t.Errorf("OrderBookModelsForDownload() return %d result", i)
		}
	})
}
