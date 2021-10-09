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
	user := config.Config[config.POSTGRES_USER]
	database := config.Config[config.POSTGRES_DB]
	host := config.Config[config.POSTGRES_SERVER_HOST]

	connString := fmt.Sprintf("postgresql://%s@%s:5432/%s?sslmode=disable", user, host, database)

	db, _ := sql.Open("postgres", connString)
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
