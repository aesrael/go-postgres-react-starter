package db

//CreateUsersTable //
func CreateUsersTable() {
	DB.Query(`
	CREATE TABLE IF NOT EXISTS users(
	user_id serial PRIMARY KEY,
	name VARCHAR (100) UNIQUE NOT NULL,
	password VARCHAR (355) NOT NULL,
	email VARCHAR (355) UNIQUE NOT NULL,
	created_on TIMESTAMP NOT NULL default current_timestamp,
	updated_at TIMESTAMP NOT NULL default current_timestamp
	)`,
	)
}
