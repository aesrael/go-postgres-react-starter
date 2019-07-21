package utils

import (
	"go-postgres-jwt-react-starter/db"
	"regexp"
)

const (
	emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

// ValidateUser returns a slice of string of validation errors
func ValidateUser(user db.Register, err []string) []string {
	emailCheck := regexp.MustCompile(emailRegex).MatchString(user.Email)
	if emailCheck != true {
		err = append(err, "Invalid email")
	}
	if len(user.Password) < 4 {
		err = append(err, "Invalid password, Password should be more than 4 characters")
	}
	if len(user.Name) < 1 {
		err = append(err, "Invalid name, please enter a name")
	}

	return err
}
