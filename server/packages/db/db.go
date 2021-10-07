package db

import (
	"database/sql"
	"fmt"
	"goapp/packages/config"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

//Connect to db
func ConnectDB() (*sql.DB, error) {
	user := config.Config["DB_USER"]
	password := config.Config["DB_PASSWORD"]
	database := config.Config["DB_NAME"]

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, database)

	db, _ := sql.Open("postgres", dbinfo)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB, dbName string) error {
	dir, err := filepath.Abs(filepath.Join("packages", "db", "migrations"))
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", dir), dbName, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
