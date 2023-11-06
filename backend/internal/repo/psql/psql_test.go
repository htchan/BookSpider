package repo

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/htchan/BookSpider/internal/config/v2"
	"github.com/htchan/BookSpider/internal/repo"
	"go.uber.org/goleak"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var code int
	defer func() { os.Exit(code) }()

	conf := config.DatabaseConfig{
		Host:     "localhost",
		Port:     "35432",
		Name:     "sqlc",
		User:     "postgres",
		Password: "postgres",
	}

	close, err := repo.CreatePsqlContainer(
		"test-sqlc-repo", conf,
		func() error {
			migrateErr := Migrate(conf, "../../../migrations")
			if migrateErr != nil {
				log.Printf("Could not open database: %s", migrateErr)

				return migrateErr
			}

			testDB, err := OpenDatabaseByConfig(conf)
			if err != nil {
				return err
			}

			return testDB.Ping()
		},
	)
	if err != nil {
		log.Println("Could not create PSQL Container: %w", err)

		return
	}
	defer close()

	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)
	} else {
		code = m.Run()
	}
}
