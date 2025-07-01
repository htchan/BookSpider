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
		"test-sqlc_v2-repo", conf,
		func() error {
			migrateErr := Migrate(conf, "../../../database/migrations")
			if migrateErr != nil {
				log.Printf("Could not migrate database: %s", migrateErr)
				code = 1

				return migrateErr
			}

			var err error
			testDB, err = OpenDatabaseByConfig(conf)
			if err != nil {
				log.Printf("Could not open database: %s", err)
				code = 1

				return err
			}

			return testDB.Ping()
		},
	)
	if err != nil {
		log.Println("Could not create PSQL Container: %w", err)
		code = 1

		return
	}
	defer close()
	defer testDB.Close()

	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)
	} else {
		code = m.Run()
	}
}
