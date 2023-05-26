package repo

import (
	"testing"

	config "github.com/htchan/BookSpider/internal/config_new"
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

func Test_OpenDatabaseByConfig(t *testing.T) {
	t.Parallel()

	conf := config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "test",
		Password: "test",
		Name:     "test_book",
	}

	result, err := OpenDatabaseByConfig(conf)

	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if result == nil {
		t.Errorf("got nil database: %v", result)
	}
	result.Close()
}
