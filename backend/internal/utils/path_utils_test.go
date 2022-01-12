package utils

import (
	"testing"
)

func Test_Utils_Exist(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		result := Exists("./path_utils_test.go")
		if !result {
			t.Fatalf("utils.Exist return false for existing file")
		}
	})

	t.Run("fail if path does not exist", func(t *testing.T) {
		result := Exists("./not_exist.go")
		if result {
			t.Fatalf("utils.Exist return true for not exist file")
		}
	})
}
