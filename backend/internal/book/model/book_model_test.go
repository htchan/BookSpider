package model

import (
	"testing"
	"github.com/google/go-cmp/cmp"
	"strings"
	"runtime"
)

func TestBookModel_SaveBookModel(t *testing.T) {
	t.Parallel()
	site := "save_bk_model"
	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	t.Run("create book with hash 0 if none of the book exist", func (t *testing.T) {
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

	t.Run("create book with model hash if hash 0 exist", func (t *testing.T) {
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

	t.Run("update book with hash 0 if model hash is 0 and it exist", func (t *testing.T) {
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
	t.Cleanup(func () {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	db.Exec(`insert into books 
	(site, id, hash_code, title, writer_id, type, update_date, update_chapter, status) 
	VALUES 
	($1, $2, 0, '', 0, '', '', '', 0), 
	($1, $2, 100, '', 0, '', '', '', 0);`,
	site, id)

	t.Run("query existing book with specific site, id, hash", func (t *testing.T) {
		t.Parallel()
		model := BookModel{Site: site, ID: id, HashCode: 0}
		result, err := QueryBookModel(db, site, id, 0)
		if err != nil || !cmp.Equal(model, result) {
			t.Errorf("query model return wrong result: error: %v; result:%v", err, result)
		}
	})

	t.Run("query latest book with specific site, id", func (t *testing.T) {
		t.Parallel()
		model := BookModel{Site: site, ID: id, HashCode: 100}
		result, err := QueryBookModel(db, site, id, -1)
		if err != nil || !cmp.Equal(model, result) {
			t.Errorf("query model return wrong result: error: %v; result:%v", err, result)
		}
	})

	t.Run("query non existence book", func (t *testing.T) {
		t.Parallel()
		result, err := QueryBookModel(db, site, -123, -1)
		if err == nil {
			t.Errorf("query model return wrong result: error: %v; result:%v", err, result)
		}
	})
}

func TestBookModel_prepareQueryTitleWriterStatement(t *testing.T) {
	t.Parallel()

	t.Run("not zero title count and writers count", func (t *testing.T) {
		t.Parallel()

		expect := `select 
		site, id, hash_code, 
		title, writer_id, type, 
		update_date, update_chapter, status 
		from books 
		where site=$1 and ((title like $2) or (writer_id=$3)) 
		order by id, hash_code limit $4 offset $5;`
		result := prepareQueryTitleWriterStatement(1, 1)
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\t", ""), "\n", "")
		result = strings.ReplaceAll(strings.ReplaceAll(result, "\t", ""), "\n", "")
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("zero title count and non zero writers count", func (t *testing.T) {
		t.Parallel()
		
		expect := `select 
		site, id, hash_code, 
		title, writer_id, type, 
		update_date, update_chapter, status 
		from books 
		where site=$1 and ((writer_id=$2)) 
		order by id, hash_code limit $3 offset $4;`
		result := prepareQueryTitleWriterStatement(0, 1)
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\t", ""), "\n", "")
		result = strings.ReplaceAll(strings.ReplaceAll(result, "\t", ""), "\n", "")
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("non zero title count and zero writers count", func (t *testing.T) {
		t.Parallel()
		
		expect := `select 
		site, id, hash_code, 
		title, writer_id, type, 
		update_date, update_chapter, status 
		from books 
		where site=$1 and ((title like $2)) 
		order by id, hash_code limit $3 offset $4;`
		result := prepareQueryTitleWriterStatement(1, 0)
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\t", ""), "\n", "")
		result = strings.ReplaceAll(strings.ReplaceAll(result, "\t", ""), "\n", "")
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("zero title count and writers count", func (t *testing.T) {
		t.Parallel()

		expect := `select 
		site, id, hash_code, 
		title, writer_id, type, 
		update_date, update_chapter, status 
		from books 
		where site=$1 and status=3 ORDER BY RANDOM() limit $2 offset $3;`
		result := prepareQueryTitleWriterStatement(0, 0)
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\t", ""), "\n", "")
		result = strings.ReplaceAll(strings.ReplaceAll(result, "\t", ""), "\n", "")
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})
}

func TestBookModel_prepareQueryTitleWriterArguement(t *testing.T) {
	t.Parallel()

	t.Run("non empty title and writers", func (t *testing.T) {
		t.Parallel()

		titles := []string{"t1", "t2"}
		writers := []int{1, 2}
		expect := []interface{} {"site", "t1", "t2", 1, 2, 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("empty title and non empty writers", func (t *testing.T) {
		t.Parallel()

		titles := []string{}
		writers := []int{1, 2}
		expect := []interface{} {"site", 1, 2, 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("non empty title and empty writers", func (t *testing.T) {
		t.Parallel()

		titles := []string{"t1", "t2"}
		writers := []int{}
		expect := []interface{} {"site", "t1", "t2", 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})

	t.Run("empty title and writers", func (t *testing.T) {
		t.Parallel()

		titles := []string{}
		writers := []int{}
		expect := []interface{} {"site", 3, 4}
		result := prepareQueryTitleWriterArgument("site", titles, writers, 3, 4)
		if !cmp.Equal(expect, result) {
			t.Error(cmp.Diff(expect, result))
		}
	})
}

func TestBookModel_QueryBookModelsByTitleWriter(t *testing.T) {
	t.Parallel()
	site := "query_bks_model"
	model_1 := BookModel{Site: site, ID: 1, HashCode: 0, Title: "abc", WriterID: 1}
	model_2 := BookModel{Site: site, ID: 1, HashCode: 1, Title: "cde", WriterID: 2}
	model_3 := BookModel{Site: site, ID: 1, HashCode: 2, Title: "efg", WriterID: 3, Status: Download}
	SaveBookModel(db, &model_1)
	SaveBookModel(db, &model_2)
	SaveBookModel(db, &model_3)

	t.Cleanup(func () {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	t.Run("providing both titles and writer ids", func (t *testing.T) {
		t.Parallel()
		
		result, err := QueryBooksModelsByTitleWriter(db, site, []string{"a"}, []int{2}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{model_1, model_2}, result) {
			t.Errorf("query book models by title or writers return error: %v", err)
		}
	})

	t.Run("providing titles only", func (t *testing.T) {
		t.Parallel()
		
		result, err := QueryBooksModelsByTitleWriter(db, site, []string{"c"}, []int{}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{model_1, model_2}, result) {
			t.Errorf("query book models by title or writers return error: %v", err)
		}
	})

	t.Run("providing writer ids only", func (t *testing.T) {
		t.Parallel()
		
		result, err := QueryBooksModelsByTitleWriter(db, site, []string{}, []int{2, 1}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{model_1, model_2}, result) {
			t.Errorf("query book models by title or writers return error: %v", err)
		}
	})

	t.Run("not providing titles and writer ids", func (t *testing.T) {
		t.Parallel()
		
		result, err := QueryBooksModelsByTitleWriter(db, site, []string{}, []int{}, 3, 0)
		if err != nil || !cmp.Equal([]BookModel{model_3}, result) {
			t.Errorf("query book models by title or writers return error: %v; result: %v", err, result)
		}
	})

	t.Run("offset all result", func (t *testing.T) {
		t.Parallel()
		
		result, err := QueryBooksModelsByTitleWriter(db, site, []string{}, []int{2, 1}, 3, 3)
		if err != nil || len(result) != 0 {
			t.Errorf("query book models by title or writers return error: %v; result: %v", err, result)
		}
	})

	t.Run("limit zero result", func (t *testing.T) {
		t.Parallel()

		result, err := QueryBooksModelsByTitleWriter(db, site, []string{}, []int{2, 1}, 0, 0)
		if err != nil || len(result) != 0 {
			t.Errorf("query book models by title or writers return error: %v", err)
		}
	})
}

func TestBookModel_QueryAllBookModels(t *testing.T) {
	t.Parallel()
	t.Skip()
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

	t.Cleanup(func () {
		db.Exec("delete from books where site=$1", site)
		runtime.GC()
	})

	t.Run("success", func (t *testing.T) {
		result, err := OrderBookModelsForUpdate(db, site)
		if err != nil || !cmp.Equal([]BookModel{model_1_2, model_2, model_3}, result) {
			t.Errorf("order book models by update date return: error: %v, result: %v", err, result)
		}
	})
}