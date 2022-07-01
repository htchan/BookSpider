package model

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSummary_Summary(t *testing.T) {
	site := "summary"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where name=$1", "new writer")
		db.Exec("delete from errors where site=$1", site)
	})

	t.Run("works", func(t *testing.T) {
		SaveBookModel(db, &BookModel{Site: site, ID: 1, Status: Download})
		SaveBookModel(db, &BookModel{Site: site, ID: 1, HashCode: 100, Status: InProgress})
		SaveBookModel(db, &BookModel{Site: site, ID: 2, Status: End})
		SaveBookModel(db, &BookModel{Site: site, ID: 3, Status: Download})
		SaveBookModel(db, &BookModel{Site: site, ID: 4, Status: Error})
		SaveWriterModel(db, &WriterModel{Name: "new writer"})
		SaveErrorModel(db, &ErrorModel{Site: site, ID: 4, Err: errors.New("error")})
		want := SummaryResult{
			BookCount:       5,
			WriterCount:     1,
			ErrorCount:      1,
			UniqueBookCount: 4,
			MaxBookID:       4,
			LatestSuccessID: 3,
			StatusCount: map[StatusCode]int{
				Error:      1,
				InProgress: 1,
				End:        1,
				Download:   2,
			},
		}
		got := Summary(db, site)
		if !cmp.Equal(got, want) {
			t.Errorf("summary return %v, want %v", got, want)
		}
	})
}
