package db

type Register struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Login struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
