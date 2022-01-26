package sqlite

import (
	"os"
	"fmt"
	"path/filepath"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

func (db *SqliteDB) BackupSchema(f *os.File) (err error) {
	defer utils.Recover(func() {})
	data, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/schema/schema.sql")
	utils.CheckError(err)
	
	f.Write(data)
	f.WriteString("\n\n")
	return
}

func (db *SqliteDB) BackupBooksTable(f *os.File) (err error) {
	defer utils.Recover(func() {})
	rows := new(SqliteBookRows)
	rows._rows, err = db._db.Query(
		"select " + database.BOOK_RECORD_FIELDS + " from books ")
	utils.CheckError(err)
	defer rows.Close()
	
	if rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)

		f.WriteString("insert into books\n(" + database.BOOK_RECORD_FIELDS + ")\nvalues\n")

		f.WriteString(fmt.Sprintf(
			"('%v', %v, %v, '%v', %v, '%v', '%v', '%v', %v)",
			record.(*database.BookRecord).Parameters()...
		))
		defer f.WriteString(";\n\n")
	}

	for rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)

		f.WriteString(fmt.Sprintf(
			",\n('%v', %v, %v, '%v', %v, '%v', '%v', '%v', %v)",
			record.(*database.BookRecord).Parameters()...
		))
	}
	return
}

func (db *SqliteDB) BackupWritersTable(f *os.File) (err error) {
	defer utils.Recover(func() {})
	rows := new(SqliteWriterRows)
	rows._rows, err = db._db.Query(
		"select " + database.WRITER_RECORD_FIELDS + " from writers ")
	utils.CheckError(err)
	defer rows.Close()
	
	if rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)

		f.WriteString("insert into writers\n(" + database.WRITER_RECORD_FIELDS + ")\nvalues\n")

		f.WriteString(fmt.Sprintf(
			"(%v, '%v')",
			record.(*database.WriterRecord).Parameters()...
		))
		defer f.WriteString(";\n\n")
	}

	for rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)

		f.WriteString(fmt.Sprintf(
			",\n(%v, '%v')",
			record.(*database.WriterRecord).Parameters()...
		))
	}
	return
}

func (db *SqliteDB) BackupErrorTable(f *os.File) (err error) {
	defer utils.Recover(func() {})
	rows := new(SqliteErrorRows)
	rows._rows, err = db._db.Query(
		"select " + database.ERROR_RECORD_FIELDS + " from errors")
	utils.CheckError(err)
	defer rows.Close()
	
	if rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)

		f.WriteString("insert into errors\n(" + database.ERROR_RECORD_FIELDS + ")\nvalues\n")

		f.WriteString(fmt.Sprintf(
			"('%v', %v, '%v')",
			record.(*database.ErrorRecord).Parameters()...
		))
		defer f.WriteString(";\n\n")
	}

	for rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)

		f.WriteString(fmt.Sprintf(
			",\n('%v', %v, '%v')",
			record.(*database.ErrorRecord).Parameters()...
		))
	}
	return
}

func (db *SqliteDB) Backup(directory, filename string) (err error) {
	defer utils.Recover(func() {})
	// construct and create the backup path
	os.MkdirAll(directory, os.ModePerm)
	destinationFileName := filepath.Join(directory, filename)
	
	file, err := os.OpenFile(destinationFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664)
	utils.CheckError(err)
	defer file.Close()

	err = db.BackupSchema(file)
	utils.CheckError(err)
	err = db.BackupBooksTable(file)
	utils.CheckError(err)
	err = db.BackupWritersTable(file)
	utils.CheckError(err)
	err = db.BackupErrorTable(file)
	utils.CheckError(err)
	return
}