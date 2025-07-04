package repo

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/htchan/BookSpider/internal/config/v2"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	KEY_PSQL_HOST     = "PSQL_HOST"
	KEY_PSQL_PORT     = "PSQL_PORT"
	KEY_PSQL_USER     = "PSQL_USER"
	KEY_PSQL_PASSWORD = "PSQL_PASSWORD"
	KEY_PSQL_NAME     = "PSQL_NAME"
)

var (
	host   = os.Getenv(KEY_PSQL_HOST)
	port   = os.Getenv(KEY_PSQL_PORT)
	dbName = os.Getenv(KEY_PSQL_NAME)

	user     = os.Getenv(KEY_PSQL_USER)
	password = os.Getenv(KEY_PSQL_PASSWORD)
)

// open database for psql
func OpenDatabase(site string) (*sql.DB, error) {
	conn := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		host, port, user, password, dbName,
	)
	database, err := sql.Open("postgres", conn)
	if err != nil {
		return database, err
	}
	// database.SetMaxIdleConns(5)
	// database.SetMaxOpenConns(10)
	// database.SetConnMaxIdleTime(5 * time.Second)
	// database.SetConnMaxLifetime(5 * time.Second)
	log.Info().Str("site", site).Msg("postgres database opened")
	return database, err
}

func OpenDatabaseByConfig(conf config.DatabaseConfig) (*sql.DB, error) {
	conn := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name,
	)
	database, err := sql.Open("postgres", conn)
	if err != nil {
		return database, err
	}

	database.SetMaxOpenConns(conf.MaxOpenConns)
	database.SetMaxIdleConns(conf.MaxIdleConns)
	database.SetConnMaxIdleTime(conf.ConnMaxIdleTime)
	log.Info().Msg("postgres database opened")
	return database, err
}

func Migrate(conf config.DatabaseConfig, migratePath string) error {
	db, dbErr := OpenDatabaseByConfig(conf)
	defer db.Close()
	if dbErr != nil {
		return fmt.Errorf("load db for migration failed: %v", dbErr)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	defer driver.Close()
	if err != nil {
		return fmt.Errorf("migrate fail: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migratePath),
		"postgres",
		driver,
	)
	defer m.Close()
	if err != nil {
		return fmt.Errorf("migrate fail: %w", err)
	}

	upErr := m.Up()
	if upErr != nil {
		log.Error().Err(upErr).Msg("migration up failed")
	}

	return nil
}

func StubPsqlConn() {
	host = "localhost"
	port = "5432"
	user = "test"
	password = "test"
	dbName = "test_book"
}
