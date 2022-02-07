package utils

import (
	"testing"
	"errors"
)

func TestUtils_Error(t *testing.T) {
	t.Run("func CheckError", func (t *testing.T) {
		t.Run("panic if error is not nil", func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("utils CheckError does not panic for error")
				}
			}()

			CheckError(errors.New("test error"))
		})

		t.Run("do nothing if error is nil", func(t *testing.T) {
			t.Parallel()
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("utils do nothing even error is not nil")
				}
			}()

			CheckError(nil)
		})
	})

	t.Run("func CheckError", func (t *testing.T) {
		t.Run("recover catch the panic error", func(t *testing.T) {
			t.Parallel()
			defer Recover(func() {})

			CheckError(errors.New("test error"))

			t.Fatalf("Recover does not catch the error")
		})

		t.Run("recover do nothing if there is no error", func(t *testing.T) {
			t.Parallel()
			defer Recover(func() {
				t.Fatalf("Recover does something even there is no error")
			})
		})
	})
}