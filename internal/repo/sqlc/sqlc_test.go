package repo

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/htchan/BookSpider/internal/config/v2"
	"github.com/htchan/BookSpider/internal/repo"
	"go.uber.org/goleak"
)

var (
	conf = config.DatabaseConfig{
		Host:     "localhost",
		Port:     "35432",
		Name:     "sqlc",
		User:     "postgres",
		Password: "postgres",
	}
)

func TestMain(m *testing.M) {
	var code int
	defer func() { os.Exit(code) }()

	close, err := repo.CreatePsqlContainer(
		"test-sqlc-repo", conf,
		func() error {
			db, err := OpenDatabaseByConfig(conf)
			defer db.Close()
			if err != nil {
				log.Printf("Could not open database: %s", err)
				code = 1

				return err
			}

			err = db.Ping()
			if err != nil {
				log.Printf("Could not ping database: %s", err)
				code = 1

				return err
			}

			migrateErr := Migrate(conf, "../../../database/migrations")
			if migrateErr != nil {
				log.Printf("Could not migrate database: %s", migrateErr)
				code = 1

				return migrateErr
			}

			return nil
		},
	)
	if err != nil {
		log.Println("Could not create PSQL Container: %w", err)
		code = 1

		return
	}
	defer close()

	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	code = 0

	if *leak {
		goleak.VerifyTestMain(
			m,
			goleak.Cleanup(func(exitCode int) {
				close()
			}),
		)
	} else {
		code = m.Run()
	}
}
