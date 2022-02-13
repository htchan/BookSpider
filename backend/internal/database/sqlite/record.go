package sqlite

import (
	"fmt"
	"html"
	"github.com/htchan/BookSpider/internal/database"
)

func BookInsertStatement(record *database.BookRecord, writerName string) string {
	if record.WriterId < 0 {
		return fmt.Sprintf(
			"insert into books " +
			"(" + database.BOOK_RECORD_FIELDS + ") " +
			"values (\"%v\", %v, %v, \"%v\", " +
			"(select id from writers where name=\"%v\"), " +
			"\"%v\", \"%v\", \"%v\", %v)",
			html.EscapeString(record.Site), record.Id, record.HashCode,
			html.EscapeString(record.Title), html.EscapeString(writerName),
			html.EscapeString(record.Type), html.EscapeString(record.UpdateDate),
			html.EscapeString(record.UpdateChapter), record.Status)
	}
	return fmt.Sprintf(
		"insert into books " +
		"(" + database.BOOK_RECORD_FIELDS + ") " +
		"values (\"%v\", %v, %v, \"%v\", %v, " +
		"\"%v\", \"%v\", \"%v\", %v)",
		html.EscapeString(record.Site), record.Id, record.HashCode,
		html.EscapeString(record.Title), record.WriterId,
		html.EscapeString(record.Type), html.EscapeString(record.UpdateDate),
		html.EscapeString(record.UpdateChapter), record.Status)
	
}

func BookUpdateStatement(record *database.BookRecord, writerName string) string {
	if record.WriterId < 0 {
		return fmt.Sprintf(
			"update books set " +
			"title=\"%v\", writer_id=(select id from writers where name=\"%v\"), " +
			"type=\"%v\", update_date=\"%v\", update_chapter=\"%v\", status=%v " +
			"where site=\"%v\" and id=%v and hash_code=%v",
			html.EscapeString(record.Title), html.EscapeString(writerName), 
			html.EscapeString(record.Type),
			html.EscapeString(record.UpdateDate), html.EscapeString(record.UpdateChapter),
			record.Status,
			html.EscapeString(record.Site), record.Id, record.HashCode)
	}
	return fmt.Sprintf(
		"update books set " +
		"title=\"%v\", writer_id=%v, type=\"%v\", " +
		"update_date=\"%v\", update_chapter=\"%v\", status=%v " +
		"where site=\"%v\" and id=%v and hash_code=%v",
		html.EscapeString(record.Title), record.WriterId, html.EscapeString(record.Type),
		html.EscapeString(record.UpdateDate), html.EscapeString(record.UpdateChapter),
		record.Status,
		html.EscapeString(record.Site), record.Id, record.HashCode)
}

func BookDeleteStatement(record *database.BookRecord) string {
	return ""
}

func WriterInsertStatement(record *database.WriterRecord) string {
	if record.Id < 0 {
		return fmt.Sprintf(
			"insert into writers " +
			"(name) values (\"%v\")",
			html.EscapeString(record.Name))
	} else {
		return fmt.Sprintf(
			"insert into writers " +
			"(id, name) values (%v, \"%v\")",
			record.Id, html.EscapeString(record.Name))

	}
}

func WriterUpdateStatement(record *database.WriterRecord) string {
	return ""
}

func WriterDeleteStatement(record *database.WriterRecord) string {
	return ""
}

func ErrorInsertStatement(record *database.ErrorRecord) string {
	return fmt.Sprintf(
		"insert into errors " +
		"(" + database.ERROR_RECORD_FIELDS + ") " +
		"values (\"%v\", %v, \"%v\")",
		html.EscapeString(record.Site), record.Id,
		html.EscapeString(record.Error.Error()))
}

func ErrorUpdateStatement(record *database.ErrorRecord) string {
	return fmt.Sprintf(
		"update errors set data=\"%v\" " +
		"where site=\"%v\" and id=%v",
		html.EscapeString(record.Error.Error()),
		html.EscapeString(record.Site), record.Id)
}

func ErrorDeleteStatement(record *database.ErrorRecord) string {
	return fmt.Sprintf(
		"delete from errors " +
		"where site=\"%v\" and id=%v",
		html.EscapeString(record.Site), record.Id)

}