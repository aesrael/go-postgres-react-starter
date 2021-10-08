package db

import (
	"database/sql"
	"fmt"
	"goapp/packages/config"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func ConnectDB() (*sql.DB, error) {
	user := config.Config["DB_USER"]
	database := config.Config["DB_NAME"]

	dbinfo := fmt.Sprintf("user=%s dbname=%s  sslmode=disable", user, database)

	db, _ := sql.Open("postgres", dbinfo)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	serverPath, _ := filepath.Abs("../")

	dir := fmt.Sprintf("%s/packages/db/migrations", serverPath)
	if err != nil {
		return err
	}
	fmt.Print(fmt.Sprintf("file://%s", dir))
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", dir), dbName, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
