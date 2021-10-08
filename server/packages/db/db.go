package db

import (
	"database/sql"
	"fmt"
	"goapp/packages/config"
	"path/filepath"

	"github.com/apex/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	log.Info("running db migrations, to disable set RUN_MIGRATION=false")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	dir, _ := filepath.Abs("../packages/db/migrations")

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
