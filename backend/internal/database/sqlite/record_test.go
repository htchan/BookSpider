package sqlite

import (
	"testing"
	"errors"
	"github.com/htchan/BookSpider/internal/database"
)

func TestSqlite_Record_BookInsertStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 100,
			Title: "title-123",
			WriterId: 123,
			Type: "type",
			UpdateDate: "update-date",
			UpdateChapter: "update-chapter",
			Status: database.InProgress,
		}
		statement := BookInsertStatement(record, "writer-1")

		if statement != "insert into books " +
			"(site, id, hash_code, title, writer_id, " +
			"type, update_date, update_chapter, status) " +
			"values (\"test\", 1, 100, \"title-123\", 123, \"type\", " +
			"\"update-date\", \"update-chapter\", 1)" {
				t.Errorf("BookInsertStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with escape content", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test\"",
			Id: 1,
			HashCode: 100,
			Title: "title-123\"",
			WriterId: 123,
			Type: "type\"",
			UpdateDate: "update-date\"",
			UpdateChapter: "update-chapter\"",
			Status: database.InProgress,
		}
		statement := BookInsertStatement(record, "writer-1")

		if statement != "insert into books " +
			"(site, id, hash_code, title, writer_id, " +
			"type, update_date, update_chapter, status) " +
			"values (\"test&#34;\", 1, 100, \"title-123&#34;\", 123, \"type&#34;\", " +
			"\"update-date&#34;\", \"update-chapter&#34;\", 1)" {
				t.Errorf("BookInsertStatement return wrong string: %v", statement)
		}
	})

	t.Run("success if writer id is negative", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 100,
			Title: "title-123",
			WriterId: -1,
			Type: "type",
			UpdateDate: "update-date",
			UpdateChapter: "update-chapter",
			Status: database.InProgress,
		}
		statement := BookInsertStatement(record, "writer-1")

		if statement != "insert into books " +
			"(site, id, hash_code, title, writer_id, " +
			"type, update_date, update_chapter, status) " +
			"values (\"test\", 1, 100, \"title-123\", " +
			"(select id from writers where name=\"writer-1\"), " +
			"\"type\", " +
			"\"update-date\", \"update-chapter\", 1)" {
				t.Errorf("BookInsertStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with escape content if writer id is negative", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test\"",
			Id: 1,
			HashCode: 100,
			Title: "title-123\"",
			WriterId: -1,
			Type: "type\"",
			UpdateDate: "update-date\"",
			UpdateChapter: "update-chapter\"",
			Status: database.InProgress,
		}
		statement := BookInsertStatement(record, "writer-1\"")

		if statement != "insert into books " +
			"(site, id, hash_code, title, writer_id, " +
			"type, update_date, update_chapter, status) " +
			"values (\"test&#34;\", 1, 100, \"title-123&#34;\", " +
			"(select id from writers where name=\"writer-1&#34;\"), " +
			"\"type&#34;\", " +
			"\"update-date&#34;\", \"update-chapter&#34;\", 1)" {
				t.Errorf("BookInsertStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_BookUpdateStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 100,
			Title: "title-123",
			WriterId: 123,
			Type: "type",
			UpdateDate: "update-date",
			UpdateChapter: "update-chapter",
			Status: database.InProgress,
		}
		statement := BookUpdateStatement(record, "writer-1")

		if statement != "update books set " +
			"title=\"title-123\", writer_id=123, type=\"type\", " +
			"update_date=\"update-date\", update_chapter=\"update-chapter\", status=1 " +
			"where site=\"test\" and id=1 and hash_code=100" {
				t.Errorf("BookUpdateStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with escape content", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test\"",
			Id: 1,
			HashCode: 100,
			Title: "title-123\"",
			WriterId: 123,
			Type: "type\"",
			UpdateDate: "update-date\"",
			UpdateChapter: "update-chapter\"",
			Status: database.InProgress,
		}
		statement := BookUpdateStatement(record, "writer-1")

		if statement != "update books set " +
			"title=\"title-123&#34;\", writer_id=123, type=\"type&#34;\", " +
			"update_date=\"update-date&#34;\", update_chapter=\"update-chapter&#34;\", status=1 " +
			"where site=\"test&#34;\" and id=1 and hash_code=100" {
				t.Errorf("BookUpdateStatement return wrong string: %v", statement)
		}
	})

	t.Run("success when writer id is negative", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 100,
			Title: "title-123",
			WriterId: -1,
			Type: "type",
			UpdateDate: "update-date",
			UpdateChapter: "update-chapter",
			Status: database.InProgress,
		}
		statement := BookUpdateStatement(record, "writer-1")

		if statement != "update books set " +
			"title=\"title-123\", writer_id=(select id from writers where name=\"writer-1\"), type=\"type\", " +
			"update_date=\"update-date\", update_chapter=\"update-chapter\", status=1 " +
			"where site=\"test\" and id=1 and hash_code=100" {
				t.Errorf("BookUpdateStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with escape content when writer id is negative", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test\"",
			Id: 1,
			HashCode: 100,
			Title: "title-123\"",
			WriterId: -1,
			Type: "type\"",
			UpdateDate: "update-date\"",
			UpdateChapter: "update-chapter\"",
			Status: database.InProgress,
		}
		statement := BookUpdateStatement(record, "writer-1\"")

		if statement != "update books set " +
			"title=\"title-123&#34;\", writer_id=(select id from writers where name=\"writer-1&#34;\"), type=\"type&#34;\", " +
			"update_date=\"update-date&#34;\", update_chapter=\"update-chapter&#34;\", status=1 " +
			"where site=\"test&#34;\" and id=1 and hash_code=100" {
				t.Errorf("BookUpdateStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_BookDeleteStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 100,
			Title: "title-123",
			WriterId: 123,
			Type: "type",
			UpdateDate: "update-date",
			UpdateChapter: "update-chapter",
			Status: database.InProgress,
		}
		statement := BookDeleteStatement(record)

		if statement != "" {
			t.Errorf("BookDeleteStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_WriterInsertStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.WriterRecord{
			Id: 100,
			Name: "test",
		}
		statement := WriterInsertStatement(record)

		if statement != "insert into writers (id, name) values (100, \"test\")" {
				t.Errorf("BookInsertStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with negative id", func(t *testing.T) {
		record := &database.WriterRecord{
			Id: -1,
			Name: "test",
		}
		statement := WriterInsertStatement(record)

		if statement != "insert into writers (name) values (\"test\")" {
				t.Errorf("BookInsertStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with escape content", func(t *testing.T) {
		record := &database.WriterRecord{
			Id: 100,
			Name: "test\"",
		}
		statement := WriterInsertStatement(record)

		if statement != "insert into writers (id, name) values (100, \"test&#34;\")"  {
				t.Errorf("BookInsertStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_WriterUpdateStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.WriterRecord{
			Id: 100,
			Name: "test",
		}
		statement := WriterUpdateStatement(record)

		if statement != "" {
			t.Errorf("WriterUpdateStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_WriterDeleteStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.WriterRecord{
			Id: 100,
			Name: "test",
		}
		statement := WriterDeleteStatement(record)

		if statement != "" {
			t.Errorf("WriterDeleteStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_ErrorInsertStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.ErrorRecord{
			Site: "test",
			Id: 100,
			Error: errors.New("test-error"),
		}
		statement := ErrorInsertStatement(record)

		if statement != "insert into errors (site, id, data) " +
			"values (\"test\", 100, \"test-error\")" {
				t.Errorf("ErrorInsertStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with escape content", func(t *testing.T) {
		record := &database.ErrorRecord{
			Site: "test\"",
			Id: 100,
			Error: errors.New("test-error\""),
		}
		statement := ErrorInsertStatement(record)

		if statement != "insert into errors (site, id, data) " +
			"values (\"test&#34;\", 100, \"test-error&#34;\")" {
				t.Errorf("ErrorInsertStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_ErrorUpdateStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.ErrorRecord{
			Site: "test",
			Id: 100,
			Error: errors.New("test-error"),
		}
		statement := ErrorUpdateStatement(record)

		if statement != "update errors set data=\"test-error\" " +
			"where site=\"test\" and id=100" {
				t.Errorf("ErrorUpdateStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with erscape content", func(t *testing.T) {
		record := &database.ErrorRecord{
			Site: "test\"",
			Id: 100,
			Error: errors.New("test-error\""),
		}
		statement := ErrorUpdateStatement(record)

		if statement != "update errors set data=\"test-error&#34;\" " +
			"where site=\"test&#34;\" and id=100" {
				t.Errorf("ErrorUpdateStatement return wrong string: %v", statement)
		}
	})
}

func TestSqlite_Record_ErrorDeleteStatement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		record := &database.ErrorRecord{
			Site: "test",
			Id: 100,
			Error: errors.New("test-error"),
		}
		statement := ErrorDeleteStatement(record)

		if statement != "delete from errors where site=\"test\" and id=100" {
			t.Errorf("ErrorDeleteStatement return wrong string: %v", statement)
		}
	})

	t.Run("success with escape content", func(t *testing.T) {
		record := &database.ErrorRecord{
			Site: "test\"",
			Id: 100,
			Error: errors.New("test-error\""),
		}
		statement := ErrorDeleteStatement(record)

		if statement != "delete from errors where site=\"test&#34;\" and id=100" {
			t.Errorf("ErrorDeleteStatement return wrong string: %v", statement)
		}
	})
}