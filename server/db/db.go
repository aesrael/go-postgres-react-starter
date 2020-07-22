package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-postgres-jwt-react-starter/server/config"
)

//DB instance
var DB *sql.DB

//Connect to db
func Connect() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.DB_USER, config.DB_PASSWORD, config.DB_NAME)

	db, _ := sql.Open("postgres", dbinfo)
	err := db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}
	DB = db
	// Create "users" table if it doesnt exist
	CreateUsersTable()
}
