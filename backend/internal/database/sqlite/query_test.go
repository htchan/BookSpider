package sqlite

import (
	"testing"
	"os"
	"io"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/utils"
)

func init() {
	source, err := os.Open("../../../assets/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./query_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func Test_Sqlite_DB_QueryBookBySiteIdHash(t *testing.T) {
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

func Test_Sqlite_DB_QueryBooksByTitle(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success if title match", func(t *testing.T) {
		query := db.QueryBooksByTitle("title-1")

		if !query.Next() {
			t.Fatalf("QueryBooksByTitle(\"title-1\") returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryBooksByTitle(\"title-1\") return %v records", n)
		}
	})

	t.Run("fail even partial title match", func(t *testing.T) {
		query := db.QueryBooksByTitle("title-")

		if query.Next() {
			t.Fatalf("QueryBooksByTitle(\"title-\") returns no record")
		}
	})

	t.Run("success if title not match", func(t *testing.T) {
		query := db.QueryBooksByTitle("-writer")

		if query.Next() {
			t.Fatalf("QueryBooksByTitle(\"title-\") returns record")
		}
	})
}

func Test_Sqlite_DB_QueryBooksByPartialTitle(t *testing.T) {
	db := NewSqliteDB("./query_test.db")
	defer db.Close()

	t.Run("success if title match", func(t *testing.T) {
		titles := []string{ "title-1" }
		query := db.QueryBooksByPartialTitle(titles)

		if !query.Next() {
			t.Fatalf("QueryBooksByPartialTitle(\"title-1\") returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 1 {
			t.Fatalf("QueryBooksByPartialTitle(\"title-1\") return %v records", n)
		}
	})

	t.Run("success if partial title match", func(t *testing.T) {
		titles := []string{ "-1", "-3" }
		query := db.QueryBooksByPartialTitle(titles)

		if !query.Next() {
			t.Fatalf("QueryBooksByPartialTitle(\"-1\", \"-3\") returns no record")
		}

		n := 1
		for ; query.Next(); n++ {}

		if n != 3 {
			t.Fatalf("QueryBooksByPartialTitle(\"-1\", \"-3\") return %v records", n)
		}
	})

	t.Run("success if title not match", func(t *testing.T) {
		titles := []string{ "-writer" }
		query := db.QueryBooksByPartialTitle(titles)

		if query.Next() {
			t.Fatalf("QueryBooksByPartialTitle(\"title-\") returns record")
		}
	})
}

func Test_Sqlite_DB_QueryBooksByWriterId(t *testing.T) {
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

func Test_Sqlite_DB_QueryBooksByStatus(t *testing.T) {
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

func Test_Sqlite_DB_QueryWriterById(t *testing.T) {
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

func Test_Sqlite_DB_QueryWriterByName(t *testing.T) {
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

func Test_Sqlite_DB_QueryWriterByPartialName(t *testing.T) {
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

func Test_Sqlite_DB_QueryErrorBySiteId(t *testing.T) {
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