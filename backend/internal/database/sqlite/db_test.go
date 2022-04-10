package sqlite

import (
	"os"
	"io"
	"testing"
	"errors"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

func initDbTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./db_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupDbTest() {
	os.Remove("./db_test.db")
}

func TestSqlite_DB_Constructor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := NewSqliteDB("./db_test.db", 100)
		defer db.Close()

		if db._db == nil {
			t.Fatalf("DB construct failed")
		}
	})
}

func TestSqlite_DB_Summary(t *testing.T) {
	db := NewSqliteDB("./db_test.db", 100)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		actual := db.Summary("test")
		
		if actual.BookCount != 6 || actual.WriterCount != 3 ||
			actual.ErrorCount != 3 || actual.UniqueBookCount != 5 ||
			actual.MaxBookId != 5 || actual.LatestSuccessId != 3 ||
			actual.StatusCount[database.Error] != 3 ||
			actual.StatusCount[database.InProgress] != 1 ||
			actual.StatusCount[database.End] != 1 ||
			actual.StatusCount[database.Download] != 1 {
			t.Fatalf(
				"DB Summary() failed\nactual: %v",
				actual)
		}
	})

	t.Run("return all 0 when failed", func(t *testing.T) {
		db := NewSqliteDB("./unknown.db", 100)
		defer db.Close()

		actual := db.Summary("test")
	
		if actual.BookCount != 0 || actual.WriterCount != 0 ||
			actual.ErrorCount != 0 || actual.UniqueBookCount != 0 ||
			actual.MaxBookId != 0 || actual.LatestSuccessId != 0 ||
			actual.StatusCount[database.Error] != 0 ||
			actual.StatusCount[database.InProgress] != 0 ||
			actual.StatusCount[database.End] != 0 ||
			actual.StatusCount[database.Download] != 0 {
			t.Fatalf(
				"DB Summary() failed\nactual: %v",
				actual)
		}
	})
}

func TestSqlite_DB_execute(t *testing.T) {
	db := NewSqliteDB("./db_test.db", 2)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		db.execute("abc")
		if db.statementCount != 1 || db.statements[0] != "abc" {
			t.Fatalf("DB.execute(abc) does not update statement Count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("commit when it reach the max statements", func(t *testing.T) {
		db.execute("abc")
		if db.statementCount != 0 {
			t.Fatalf("DB.execute(abc) does not update statement Count: %v, statements: %v", db.statementCount, db.statements)
		}
	})
}

func TestSqlite_DB_Commit(t *testing.T) {
	db := NewSqliteDB("./db_test.db", 100)
	defer db.Close()

	t.Run("success for insert normal book record", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 10,
			HashCode: 10,
			Title: "title-10",
			WriterId: 1,
			Type: "type-10",
			UpdateDate: "update-date-10",
			UpdateChapter: "update-chapter-10",
			Status: database.InProgress,
		}
		db.execute(BookInsertStatement(bookRecord, "writer-1"))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit insert record return error: %v", err)
		}

		rows := db.QueryBookBySiteIdHash("test", 10, 10)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query insert record return error: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 10 ||
			actualRecord.HashCode != 10 || actualRecord.Title != "title-10" ||
			actualRecord.WriterId != 1 || actualRecord.Type != "type-10" ||
			actualRecord.UpdateDate != "update-date-10" ||
			actualRecord.UpdateChapter != "update-chapter-10" ||
			actualRecord.Status != database.InProgress {
				t.Fatalf("DB.Commit fail for insert normal book record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for insert book record with negative writer id", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 11,
			HashCode: 11,
			Title: "title-11",
			WriterId: -1,
			Type: "type-11",
			UpdateDate: "update-date-11",
			UpdateChapter: "update-chapter-11",
			Status: database.InProgress,
		}
		db.execute(BookInsertStatement(bookRecord, "writer-1"))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit insert record return error: %v", err)
		}

		rows := db.QueryBookBySiteIdHash("test", 11, 11)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query insert record return error: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 11 ||
			actualRecord.HashCode != 11 || actualRecord.Title != "title-11" ||
			actualRecord.WriterId != 1 || actualRecord.Type != "type-11" ||
			actualRecord.UpdateDate != "update-date-11" ||
			actualRecord.UpdateChapter != "update-chapter-11" ||
			actualRecord.Status != database.InProgress {
				t.Fatalf("DB.Commit fail for insert book record with negative writer id - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for update normal book record", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 10,
			HashCode: 10,
			Title: "title-10-new",
			WriterId: 1,
			Type: "type-10-new",
			UpdateDate: "update-date-10-new",
			UpdateChapter: "update-chapter-10-new",
			Status: database.InProgress,
		}
		db.execute(BookUpdateStatement(bookRecord, "writer-1"))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit update record return error: %v", err)
		}

		rows := db.QueryBookBySiteIdHash("test", 10, 10)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query update record return error: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 10 ||
			actualRecord.HashCode != 10 || actualRecord.Title != "title-10-new" ||
			actualRecord.WriterId != 1 || actualRecord.Type != "type-10-new" ||
			actualRecord.UpdateDate != "update-date-10-new" ||
			actualRecord.UpdateChapter != "update-chapter-10-new" ||
			actualRecord.Status != database.InProgress {
				t.Fatalf("DB.Commit fail for update normal book record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for update book record with negative writer id", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 11,
			HashCode: 11,
			Title: "title-11-new",
			WriterId: -1,
			Type: "type-11-new",
			UpdateDate: "update-date-11-new",
			UpdateChapter: "update-chapter-11-new",
			Status: database.InProgress,
		}
		db.execute(BookUpdateStatement(bookRecord, "writer-3"))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit update record return error: %v", err)
		}

		rows := db.QueryBookBySiteIdHash("test", 11, 11)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query update record return error: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 11 ||
			actualRecord.HashCode != 11 || actualRecord.Title != "title-11-new" ||
			actualRecord.WriterId != 3 || actualRecord.Type != "type-11-new" ||
			actualRecord.UpdateDate != "update-date-11-new" ||
			actualRecord.UpdateChapter != "update-chapter-11-new" ||
			actualRecord.Status != database.InProgress {
				t.Fatalf("DB.Commit fail for update book record with negative writer id - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("do nothing for delete book record", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 11,
			HashCode: 11,
			Title: "title-11-new",
			WriterId: 1,
			Type: "type-11-new",
			UpdateDate: "update-date-11-new",
			UpdateChapter: "update-chapter-11-new",
			Status: database.InProgress,
		}
		db.execute(BookDeleteStatement(bookRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit delete record return error: %v", err)
		}

		rows := db.QueryBookBySiteIdHash("test", 11, 11)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query delete record return error: %v", err)
		}
		actualRecord := record.(*database.BookRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 11 ||
			actualRecord.HashCode != 11 || actualRecord.Title != "title-11-new" ||
			actualRecord.WriterId != 3 || actualRecord.Type != "type-11-new" ||
			actualRecord.UpdateDate != "update-date-11-new" ||
			actualRecord.UpdateChapter != "update-chapter-11-new" ||
			actualRecord.Status != database.InProgress {
				t.Fatalf("DB.Commit fail for insert normal book record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for insert normal writer record", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: 10,
			Name: "writer-10",
		}
		db.execute(WriterInsertStatement(writerRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit insert record return error: %v", err)
		}

		rows := db.QueryWriterById(10)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query insert record return error: %v", err)
		}
		actualRecord := record.(*database.WriterRecord)
		if actualRecord.Id != 10 || actualRecord.Name != "writer-10" {
			t.Fatalf("DB.Commit fail for insert normal writer record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for insert writer record without writer id", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: -1,
			Name: "writer-11",
		}
		db.execute(WriterInsertStatement(writerRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit insert record return error: %v", err)
		}

		rows := db.QueryWriterById(11)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query insert record return error: %v", err)
		}
		actualRecord := record.(*database.WriterRecord)
		if actualRecord.Id != 11 || actualRecord.Name != "writer-11" {
			t.Fatalf("DB.Commit fail for insert normal writer record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("do nothing for update writer record", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: 10,
			Name: "writer-10-new",
		}
		db.execute(WriterUpdateStatement(writerRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit update record return error: %v", err)
		}

		rows := db.QueryWriterById(10)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query update record return error: %v", err)
		}
		actualRecord := record.(*database.WriterRecord)
		if actualRecord.Id != 10 || actualRecord.Name != "writer-10" {
			t.Fatalf("DB.Commit fail for update normal writer record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("do nothing for delete wrtier record", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: 10,
			Name: "writer-10",
		}
		db.execute(WriterDeleteStatement(writerRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit delete record return error: %v", err)
		}

		rows := db.QueryWriterById(10)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query delete record return error: %v", err)
		}
		actualRecord := record.(*database.WriterRecord)
		if actualRecord.Id != 10 || actualRecord.Name != "writer-10" {
			t.Fatalf("DB.Commit fail for delete normal writer record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for insert error record", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 10,
			Error: errors.New("error-10"),
		}
		db.execute(ErrorInsertStatement(errorRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit insert record return error: %v", err)
		}

		rows := db.QueryErrorBySiteId("test", 10)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query insert record return error: %v", err)
		}
		actualRecord := record.(*database.ErrorRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 10 ||
			actualRecord.Error.Error() != "error-10" {
			t.Fatalf("DB.Commit fail for insert error record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for update error record", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 10,
			Error: errors.New("error-10-new"),
		}
		db.execute(ErrorUpdateStatement(errorRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit update record return error: %v", err)
		}

		rows := db.QueryErrorBySiteId("test", 10)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query update record return error: %v", err)
		}
		actualRecord := record.(*database.ErrorRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 10 ||
			actualRecord.Error.Error() != "error-10-new" {
			t.Fatalf("DB.Commit fail for update error record - err: %v, record: %v", err, actualRecord)
		}
	})

	t.Run("success for delete error record", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 10,
			Error: errors.New("error-10-new"),
		}
		db.execute(ErrorDeleteStatement(errorRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit delete record return error: %v", err)
		}

		rows := db.QueryErrorBySiteId("test", 10)
		defer rows.Close()
		_, err = rows.Scan()
		if err == nil {
			t.Fatalf("query delete record not return error")
		}
	})

	t.Run("do nothing to the empty string", func(t *testing.T) {
		db.execute("")
		db.execute("")
		db.execute("")
		db.execute("")
		db.execute("")
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit insert record return error: %v", err)
		}
	})

	t.Run("do nothing to the repeated string", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 11,
			Error: errors.New("error-11"),
		}
		db.execute(ErrorInsertStatement(errorRecord))
		db.execute(ErrorInsertStatement(errorRecord))
		err := db.Commit()
		if err != nil {
			t.Fatalf("commit insert record return error: %v", err)
		}

		rows := db.QueryErrorBySiteId("test", 11)
		defer rows.Close()
		record, err := rows.Scan()
		if err != nil {
			t.Fatalf("query insert record return error: %v", err)
		}
		actualRecord := record.(*database.ErrorRecord)
		if actualRecord.Site != "test" || actualRecord.Id != 11 ||
			actualRecord.Error.Error() != "error-11" {
			t.Fatalf("DB.Commit fail for insert error record - err: %v, record: %v", err, actualRecord)
		}
	})
}

func TestSqlite_DB_Close(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := NewSqliteDB("./db_test.db", 100)
		db.Close()
		if db._db != nil {
			t.Fatalf("DB Close() failed")
		}
	})
}

func TestSqlite_DB_interface(t *testing.T) {
	var db database.DB
	db = NewSqliteDB("./db_test.db", 100)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		switch db.(type) {
		case *SqliteDB:
		case interface{}:
			t.Fatalf("NewSqliteDB cannot use db interface")
		}
	})

	t.Run("success in query", func(t *testing.T) {
		query := db.QueryBookBySiteIdHash("test", 1, 100)
		switch query.(type) {
		case *SqliteBookRows:
		case interface{}:
			t.Fatalf("NewSqliteDB query does not return *BookRecord")
		}
	})
}

func TestMain(m *testing.M) {
	initDbTest()
	initDbBackupTest()
	initDbCreateTest()
	initDbQueryTest()
	initDbUpdateTest()
	initRowTest()

	code := m.Run()

	cleanupRowTest()
	cleanupDbUpdateTest()
	cleanupDbQueryTest()
	cleanupDbCreateTest()
	cleanupDbBackupTest()
	cleanupDbTest()

	os.Exit(code)
}