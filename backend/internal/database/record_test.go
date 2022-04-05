package database

import (
	"testing"
	"errors"
)

func Test_record(t *testing.T) {
	t.Run("struct BookRecord", func(t *testing.T) {
		record := BookRecord{
			Site: "site-name",
			Id: 0,
			HashCode: 100,
			Title: "book-title",
			WriterId: 0,
			Type: "book-type",
			UpdateDate: "book-update-date",
			UpdateChapter: "book-update-chapter",
			Status: Error,
		}

		t.Run("func Parameters", func(t *testing.T) {
			result := record.Parameters()

			if result[0].(string) != "site-name" ||
				result[1].(int) != 0 ||
				result[2].(int) != 100 ||
				result[3].(string) != "book-title" ||
				result[4].(int) != 0 ||
				result[5].(string) != "book-type" ||
				result[6].(string) != "book-update-date" ||
				result[7].(string) != "book-update-chapter" ||
				result[8].(StatusCode) != Error {
				t.Fatalf("got result: %v", result)
			}
		})
		
		t.Run("func String", func(t *testing.T) {
			result := record.String()

			if result != "site-name-0-2s" {
				t.Fatalf("got result: %v", result)
			}
		})

		t.Run("func Equal", func(t *testing.T) {
			t.Run("same", func(t *testing.T) {
				success := record.Equal(record)
				if !success {
					t.Fatalf("got result: %v", success)
				}
			})

			t.Run("different", func(t *testing.T) {
				record2 := record
				record2.Status = Download
				failure := record.Equal(record2)
				if failure {
					t.Fatalf("got result: %v", failure)
				}
			})
		})
	})

	t.Run("struct WriterRecord", func(t *testing.T) {
		record := WriterRecord{
			Id: 0,
			Name: "writer-name",
		}
		
		t.Run("func Parameters", func(t *testing.T) {
			result := record.Parameters()

			if result[0].(int) != 0 ||
				result[1].(string) != "writer-name" {
				t.Fatalf("got result: %v", result)
			}
		})

		t.Run("func Equal", func(t *testing.T) {
			t.Run("same", func(t *testing.T) {
				success := record.Equal(record)
				if !success {
					t.Fatalf("got result: %v", success)
				}
			})

			t.Run("different", func(t *testing.T) {
				record2 := record
				record2.Name = "other-writer-name"
				failure := record.Equal(record2)
				if failure {
					t.Fatalf("got result: %v", failure)
				}
			})
		})
	})

	t.Run("struct ErrorRecord", func(t *testing.T) {
		record := ErrorRecord{
			Site: "site-name",
			Id: 0,
			Error: errors.New("error-desc"),
		}
		
		t.Run("func Parameters", func(t *testing.T) {
			result := record.Parameters()

			if result[0].(string) != "site-name" ||
				result[1].(int) != 0 ||
				result[2].(string) != "error-desc" {
				t.Fatalf("got result: %v", result)
			}
		})

		t.Run("func Equal", func(t *testing.T) {
			t.Run("same", func(t *testing.T) {
				success := record.Equal(record)
				if !success {
					t.Fatalf("got result: %v", success)
				}
			})

			t.Run("different", func(t *testing.T) {
				record2 := record
				record2.Id = 999
				failure := record.Equal(record2)
				if failure {
					t.Fatalf("got result: %v", failure)
				}
			})
		})
	})
}