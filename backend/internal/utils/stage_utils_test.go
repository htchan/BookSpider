package utils

import (
	"testing"
	"os"
)

func Test_WriteStage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		StageFileName = os.Getenv("ASSETS_LOCATION") + "/test-data/test.log"
		defer os.Remove(StageFileName)

		WriteStage("hello")

		b, err := os.ReadFile(StageFileName)
		CheckError(err)

		if string(b) != "hello\n"{
			t.Errorf("book saveContent save such content: %v", string(b))
		}
	})
	
	t.Run("fail to write if error happen", func(t *testing.T) {
		StageFileName = os.Getenv("ASSETS_LOCATION") + "/test-data/fail.log"
		defer os.Remove(StageFileName)
		os.OpenFile(StageFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
		os.Chmod(StageFileName, 0444)
		defer os.Chmod(StageFileName, 0664)

		WriteStage("hello")
		b, err := os.ReadFile(StageFileName)
		CheckError(err)

		if string(b) != ""{
			t.Errorf("book saveContent save such content: %v", string(b))
		}
	})
}