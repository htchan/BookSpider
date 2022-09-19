package repo

import (
	"testing"
)

func Test_OpenDatabase(t *testing.T) {
	t.Parallel()
	StubPsqlConn()

	result, err := OpenDatabase("test")

	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if result == nil {
		t.Errorf("got nil database: %v", result)
	}
	result.Close()
}
