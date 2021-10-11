package db

const (
	CreateUserQuery     = `INSERT INTO users(id, name, password, email) VALUES (DEFAULT, $1 , $2, $3);`
	GetUserByIDQuery    = `SELECT * FROM users WHERE id = $1`
	GetUserByEmailQuery = `SELECT * FROM users WHERE email = $1`
)
