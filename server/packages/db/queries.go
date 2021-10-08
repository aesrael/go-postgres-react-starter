package db

const (
	CheckUserExists         = `SELECT true from users WHERE email = $1`
	LoginQuery              = `SELECT * from users WHERE email = $1`
	CreateUserQuery         = `INSERT INTO users(id, name, password, email) VALUES (DEFAULT, $1 , $2, $3);`
	UpdateUserPasswordQuery = `UPDATE users SET password = $2 WHERE id = $1`
	GetUserByIDQuery        = `SELECT * FROM users WHERE id = $1`
	GetUserByEmailQuery     = `SELECT * FROM users WHERE email = $1`
)
