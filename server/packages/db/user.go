package db

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type ResetPassword struct {
	ID              int    `json:"id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type Login struct {
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type CreateReset struct {
	Email string `json:"email"`
}

type User struct {
	ID        string `json:"id,omitempty"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) UserExists(dbConn *sql.DB) bool {
	rows, err := dbConn.Query(GetUserByEmailQuery, user.Email)
	if err != nil || !rows.Next() {
		return false
	}

	return true
}
