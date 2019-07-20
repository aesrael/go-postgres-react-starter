package db

//CreateUsersTable //
func CreateUsersTable() {
	DB.Query(`
	CREATE TABLE IF NOT EXISTS users(
	user_id serial PRIMARY KEY,
	username VARCHAR (50) UNIQUE NOT NULL,
	password VARCHAR (50) NOT NULL,
	email VARCHAR (355) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
	)`,
	)
}
