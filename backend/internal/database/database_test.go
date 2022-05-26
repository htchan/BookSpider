package database

import (
	"testing"
	"time"
)

func Test_database(t *testing.T) {
	t.Run("func GenerateHash", func (t *testing.T) {
		result := GenerateHash()
		time.Sleep(1 * time.Second)
		compareResult := GenerateHash()
		if compareResult - result < 1000 || compareResult - result > 1005 {
			t.Errorf("got difference: %v", compareResult - result)
		}
	})

	t.Run("func StatustoString", func (t *testing.T) {
		t.Run("Error status", func (t *testing.T) {
			result := StatustoString(Error)
			if result != "error" {
				t.Errorf("got result: %v", result)
			}
		})

		t.Run("InProgress status", func (t *testing.T) {
			result := StatustoString(InProgress)
			if result != "in_progress" {
				t.Errorf("got result: %v", result)
			}
		})

		t.Run("End status", func (t *testing.T) {
			result := StatustoString(End)
			if result != "end" {
				t.Errorf("got result: %v", result)
			}
		})

		t.Run("Download status", func (t *testing.T) {
			result := StatustoString(Download)
			if result != "download" {
				t.Errorf("got result: %v", result)
			}
		})
	})
}