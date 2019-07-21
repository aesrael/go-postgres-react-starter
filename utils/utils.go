package utils

import (
	"go-postgre-jwt-boilerplate/db"
	"regexp"
)

const (
	emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

func ValidateUser(user db.Register, err []string) []string {
	emailCheck := regexp.MustCompile(emailRegex).MatchString(user.Email)
	if emailCheck != true {
		err[0] = "Invalid email"
	}
	if len(user.Password) < 4 {
		err[1] = "Invalid password"
	}
	if len(user.Name) < 1 {
		err[2] = "Invalid name"
	}

	return err
}
