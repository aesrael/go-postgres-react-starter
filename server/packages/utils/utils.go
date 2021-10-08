package utils

import (
	"goapp/packages/db"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

func ValidateUser(user db.User) []string {
	errs := []string{}
	isValidEmail := regexp.MustCompile(emailRegex).MatchString(user.Email)
	if !isValidEmail {
		errs = append(errs, "Invalid email")
	}
	if len(user.Password) < 4 {
		errs = append(errs, "Invalid password, Password should be more than 4 characters")
	}
	if len(user.Name) < 1 {
		errs = append(errs, "Invalid name, please enter a name")
	}

	return errs
}

func ValidatePasswordReset(resetPassword db.ResetPassword) (bool, string) {
	if len(resetPassword.Password) < 4 {
		return false, "Invalid password, password should be more than 4 characters"
	}
	if resetPassword.Password != resetPassword.ConfirmPassword {
		return false, "Password reset failed, passwords must match"
	}
	return true, ""
}

func GetHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
