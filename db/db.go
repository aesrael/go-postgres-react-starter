package db

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	DB_USER     = "israel"
	DB_PASSWORD = "secret"
	DB_NAME     = "goauth"
)

//DB instance
var DB *sql.DB

//Connect to db
func Connect() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)

	db, _ := sql.Open("postgres", dbinfo)
	err := db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}
	DB = db
	// Create "users" table if it doesnt exist
	CreateUsersTable()
}
