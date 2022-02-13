package sqlite

import (
	"testing"
	"os"
	"io"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/utils"
)

func initDbQueryTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./query_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupDbQueryTest() {
	os.Remove("./query_test.db")
}

func TestSqlite_DB_QueryBookBySiteIdHash(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success with specific hash", func(t *testing.T) {
		query := db.QueryBookBySiteIdHash("test", 3, 102)

		if !query.Next() {
			t.Fatalf("QueryBookBySiteIdHash(\"test\", 3, 102) return no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryBookBySiteIdHash(\"test\", 3, 102) return %v records", n)
		}
	})

	t.Run("success without specifying hash", func(t *testing.T) {
		query := db.QueryBookBySiteIdHash("test", 3, -1)

		if !query.Next() {
			t.Fatalf("QueryBookBySiteIdHash(\"test\", 3, -1) return no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 2 {
			t.Fatalf("QueryBookBySiteIdHash(\"test\", 3, 102) return %v records", n)
		}
	})

	t.Run("fail if querying not exist book", func(t *testing.T) {
		query := db.QueryBookBySiteIdHash("test", -1, -1)

		if query.Next() {
			record, err := query.Scan()
			t.Fatalf("QueryBookBySiteIdHash(\"test\", -1, -1) return result record: %v, err: %v", record, err)
		}
	})
}

func TestSqlite_DB_QueryBooksByPartialTitleAndWriter(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("return result match any title or writer", func(t *testing.T) {
		titles := []string{ "-new", "-3" }
		writers := []int { 1, 5 }
		query := db.QueryBooksByPartialTitleAndWriter(titles, writers)

		if !query.Next() {
			t.Fatalf("QueryBooksByPartialTitle([\"-new\", \"-3\"], [1, 5]) returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 3 {
			t.Fatalf("QueryBooksByPartialTitle([\"-new\", \"-3\"], [1, 5]) return %v records", n)
		}
	})

	t.Run("return result match any title", func(t *testing.T) {
		titles := []string{ "-new", "-3" }
		writers := []int {}
		query := db.QueryBooksByPartialTitleAndWriter(titles, writers)

		if !query.Next() {
			t.Fatalf("QueryBooksByPartialTitle([\"-new\", \"-3\"], []) returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 2 {
			t.Fatalf("QueryBooksByPartialTitle([\"-new\", \"-3\"], []) return %v records", n)
		}
	})

	t.Run("return result match any writer", func(t *testing.T) {
		titles := []string{}
		writers := []int { 1, 5 }
		query := db.QueryBooksByPartialTitleAndWriter(titles, writers)

		if !query.Next() {
			t.Fatalf("QueryBooksByPartialTitle([], [1, 5]) returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryBooksByPartialTitle([], [1, 5]) return %v records", n)
		}
	})

	t.Run("return no result if no title or writer match", func(t *testing.T) {
		titles := []string{ "-writer" }
		writers := []int{ 5 }
		query := db.QueryBooksByPartialTitleAndWriter(titles, writers)

		if query.Next() {
			t.Fatalf("QueryBooksByPartialTitleAndWriter(\"-writer\", 5) returns record")
		}
	})
}

func TestSqlite_DB_QueryBooksByWriterId(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		query := db.QueryBooksByWriterId(1)

		if !query.Next() {
			t.Fatalf("QueryBooksByWriterId(1) returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryBooksByWriterId(1) return %v records", n)
		}
	})

	t.Run("fail if writer id not exist", func(t *testing.T) {
		query := db.QueryBooksByWriterId(-1)

		if query.Next() {
			t.Fatalf("QueryBooksByWriterId(-1) returns record")
		}
	})
}

func TestSqlite_DB_QueryBooksByStatus(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		query := db.QueryBooksByStatus(database.Error)

		if !query.Next() {
			t.Fatalf("QueryBooksByStatus(database.Error) returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 3 {
			t.Fatalf("QueryBooksByStatus(database.Error) return %v records", n)
		}
	})
}

func TestSqlite_DB_QueryWriterById(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		query := db.QueryWriterById(1)

		if !query.Next() {
			t.Fatalf("QueryWriterById(1) returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryWriterById(1) return %v records", n)
		}
	})

	t.Run("fail if writer not exist", func(t *testing.T) {
		query := db.QueryWriterById(-1)

		if query.Next() {
			t.Fatalf("QueryWriterById(-1) returns record")
		}
	})
}

func TestSqlite_DB_QueryWriterByName(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		query := db.QueryWriterByName("writer-1")

		if !query.Next() {
			t.Fatalf("QueryWriterByName(\"title\") returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryWriterByName(\"title\") return %v records", n)
		}
	})

	t.Run("fail even partial match", func(t *testing.T) {
		query := db.QueryWriterByName("writer")

		if query.Next() {
			t.Fatalf("QueryWriterByName(\"writer\") returns no record")
		}
	})

	t.Run("fail if writer not exist", func(t *testing.T) {
		query := db.QueryWriterByName("title")

		if query.Next() {
			t.Fatalf("QueryWriterByName(\"title\") returns record")
		}
	})
}

func TestSqlite_DB_QueryWriterByPartialName(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success if full match", func(t *testing.T) {
		names := []string{ "writer-1" }
		query := db.QueryWritersByPartialName(names)

		if !query.Next() {
			t.Fatalf("QueryWritersByPartialName(\"writer-1\") returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryWritersByPartialName(\"writer-1\") return %v records", n)
		}
	})

	t.Run("success if partial match", func(t *testing.T) {
		names := []string{ "writer" }
		query := db.QueryWritersByPartialName(names)

		if !query.Next() {
			t.Fatalf("QueryWritersByPartialName(\"writer\") returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 3 {
			t.Fatalf("QueryWritersByPartialName(\"writer\") return %v records", n)
		}
	})

	t.Run("fail if writer not exist", func(t *testing.T) {
		names := []string{ "title" }
		query := db.QueryWritersByPartialName(names)

		if query.Next() {
			t.Fatalf("QueryWritersByPartialName(\"title\") returns record")
		}
	})}

func TestSqlite_DB_QueryErrorBySiteId(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		query := db.QueryErrorBySiteId("test", 2)

		if !query.Next() {
			t.Fatalf("QueryErrorBySiteId(\"test\", 2) returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryErrorBySiteId(\"test\", 2) return %v records", n)
		}
	})

	t.Run("fail if error not exist", func(t *testing.T) {
		query := db.QueryErrorBySiteId("test", 1)

		if query.Next() {
			t.Fatalf("QueryErrorBySiteId(\"test\", 1) returns  record")
		}
	})
}

func TestSqlite_DB_QueryBooksOrderByUpdateDate(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success with desc order update date", func(t *testing.T) {
		rows := db.QueryBooksOrderByUpdateDate()
		n := 0
		expectTitle := []string { "title-1", "title-3", "title-3-new" }
		for ; rows.Next(); n++ {
			record, err := rows.ScanCurrent()
			if err != nil || record.(*database.BookRecord).Title != expectTitle[n] {
				t.Logf("invalid result at position %v: record: %v, err: %v", n, record, err)
				defer t.Fatal()
			}
		}
	})
}

func TestSqlite_DB_QueryBooksWithRandomOrder(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success query all books with status = error", func(t *testing.T) {
		rows := db.QueryBooksWithRandomOrder(4, database.Error)
		n := 0
		for ; rows.Next(); n++ {}
		if n != 4 {
			t.Fatalf("invalid count: n: %v", n)
		}
	})
	t.Run("success query specific in progress books", func(t *testing.T) {
		rows := db.QueryBooksWithRandomOrder(1, database.InProgress)
		record, err := rows.Scan()
		if err != nil || record.(*database.BookRecord).Title != "title-1" || rows.Next() {
			t.Fatalf("invalid result: result: %v, err: %v", record, err)
		}
	})
	t.Run("success query books even n > total books in db", func(t *testing.T) {
		rows := db.QueryBooksWithRandomOrder(10, database.Error)
		n := 0
		for ; rows.Next(); n++ {}
		if n != 6 {
			t.Fatalf("invalid count: n: %v", n)
		}
	})
}