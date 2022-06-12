package book

// import (
// 	"testing"
// 	"sync"
// 	"fmt"

// 	"github.com/htchan/BookSpider/internal/book/model"
// 	"github.com/htchan/BookSpider/pkg/config"
// )

// func TestBook_ConcurrentSave(t *testing.T) {
// 	t.Parallel()

// 	var wg sync.WaitGroup
	
// 	site, writer := "thread_save_bk", "save_bk_writer"
// 	con := config.BookConfig{}

// 	t.Cleanup(func () {
// 		db.Exec("delete from books where site=$1", site)
// 		db.Exec("delete from errors where site=$1", site)
// 		db.Exec("delete from writers where name like $1", writer + "%")
// 	})

// 	for i := 0; i < 99999; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			expect := Book{
// 				BookModel: model.BookModel{
// 					Site: site, ID: i, HashCode: 0,
// 					Status: model.InProgress,
// 				},
// 				WriterModel: model.WriterModel{Name: fmt.Sprintf(writer + "_%v", i % 100)},
// 				BookConfig: &con,
// 			}
// 			err := expect.Save(db)
// 			if err != nil || expect.BookModel.WriterID == 0 ||
// 			expect.WriterModel.ID == 0 {
// 				t.Errorf("save book return: %v", err)
// 				t.Errorf("save book book model update: %v", expect.BookModel)
// 				t.Errorf("save book writer model update: %v", expect.WriterModel)
// 				t.Errorf("save book error model update: %v", expect.ErrorModel)
// 				return
// 			}
// 		}()
// 	}
// 	wg.Wait()
// }