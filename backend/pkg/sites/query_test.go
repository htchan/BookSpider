package sites

import (
	"testing"
	"io"
	"os"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
)

func initQueryTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./query.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupQueryTest() {
	os.Remove("./query.db")
}

var queryConfig = configs.LoadSiteConfigs(os.Getenv("ASSETS_LOCATION") + "/test-data/configs")["test"]

func Test_Sites_Site_Query(t *testing.T) {
	queryConfig.DatabaseLocation = "./query.db"
	site := NewSite("test", queryConfig)
	site.OpenDatabase()
	defer site.CloseDatabase()

	t.Run("func SearchByIdHash", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			book := site.SearchByIdHash(1, "2s")
			
			if book == nil || book.GetTitle() != "title-1" ||
				book.GetWriter() != "writer-1" {
					t.Fatalf("book search by id hash return wrong result: %v", book)
			}
		})

		t.Run("success without specify hash", func(t *testing.T) {
			book := site.SearchByIdHash(1, "")
			
			if book == nil || book.GetTitle() != "title-1" ||
				book.GetWriter() != "writer-1" {
					t.Fatalf("book search by id hash return wrong result: %v", book)
			}
		})

		t.Run("fail for not exist id", func(t *testing.T) {
			book := site.SearchByIdHash(999, "")
			
			if book != nil {
				t.Fatalf("book search by id hash return some result with not exist id: %v", book)
			}
		})
	})

	t.Run("func SearchByWriterId", func(t *testing.T) {
		queryConfig.DatabaseLocation = "./query.db"
		site := NewSite("test", queryConfig)
		site.OpenDatabase()
		defer site.CloseDatabase()

		t.Run("success", func(t *testing.T) {
			books := site.SearchByWriterId(1)
			
			if len(books) != 1 || books[0].GetTitle() != "title-1" ||
				books[0].GetWriter() != "writer-1" {
					t.Fatalf("book search by id hash return wrong result: %v", books)
			}
		})

		t.Run("return 0 result if writer not exist", func(t *testing.T) {
			books := site.SearchByWriterId(999)
			
			if len(books) != 0 {
				t.Fatalf("book search by id hash return wrong result: %v", books)
			}
		})
	})

	t.Run("func SearchByStatus", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			books := site.SearchByStatus(database.End)
			
			if len(books) != 1 || books[0].GetTitle() != "title-3-new" ||
				books[0].GetWriter() != "writer-3" {
					t.Fatalf("book search by id hash return wrong result: %v", books)
			}
		})

		t.Run("return all records match status", func(t *testing.T) {
			books := site.SearchByStatus(database.Error)
			
			if len(books) != 3 {
				t.Fatalf("book search by id hash return wrong result: %v", books)
			}
		})
	})

	t.Run("func SearchByTitleWriter", func(t *testing.T) {
		t.Run("success with partial title", func(t *testing.T) {
			books := site.SearchByTitleWriter("new", "")

			if len(books) != 1 {
				t.Fatalf("site Query book wrong number of books: %v", len(books))
			}
			if book := books[0]; book.GetTitle() != "title-3-new" ||
				book.GetWriter() != "writer-3" || book.GetStatus() != database.End {
					t.Fatalf("book query wrong books")
			}
		})

		t.Run("success with partial writer name", func(t *testing.T) {
			books := site.SearchByTitleWriter("", "-3")

			if len(books) != 1 {
				t.Fatalf("site Query book wrong number of books: %v", len(books))
			}
			if book := books[0]; book.GetTitle() != "title-3-new" ||
				book.GetWriter() != "writer-3" || book.GetStatus() != database.End {
					t.Fatalf("book query wrong books")
			}
		})

		t.Run("return nothing when nothing is searched", func(t *testing.T) {
			books := site.SearchByTitleWriter("", "")

			if len(books) != 0 {
				t.Fatalf("site Query book wrong number of books: %v", len(books))
			}
		})
	})

	t.Run("func RandomSuggestBooks", func(t *testing.T) {
		t.Run("success for all books", func(t *testing.T) {
			books := site.RandomSuggestBook(50, database.Error)

			if len(books) != 6 {
				t.Fatalf("site Random Suggest Books return %v books", len(books))
			}
		})

		t.Run("success for specific status", func(t *testing.T) {
			books := site.RandomSuggestBook(50, database.InProgress)

			if len(books) != 1 {
				t.Fatalf("site Random Suggest Books return %v books", len(books))
			}

			if books[0].GetTitle() != "title-1" {
				t.Fatalf("site Random Suggest Books does not return target status book: %v", books[0].GetTitle())
			}
		})
	})
}