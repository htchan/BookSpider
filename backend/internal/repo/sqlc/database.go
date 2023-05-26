package repo

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	config "github.com/htchan/BookSpider/internal/config_new"
	_ "github.com/lib/pq"
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
	log.Printf("postgres_database.open; %v", database)
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
	// database.SetMaxIdleConns(5)
	// database.SetMaxOpenConns(10)
	// database.SetConnMaxIdleTime(5 * time.Second)
	// database.SetConnMaxLifetime(5 * time.Second)
	log.Printf("postgres_database.open; %v", database)
	return database, err
}

func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate fail: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate fail: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		log.Printf("migration: %s", err)
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
