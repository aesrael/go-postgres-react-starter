package utils

import (
	"goapp/packages/db"
	"regexp"
)

const (
	emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

func ValidateUser(user db.Register) []string {
	errs := []string{}
	emailCheck := regexp.MustCompile(emailRegex).MatchString(user.Email)
	if emailCheck != true {
		errs = append(err, "Invalid email")
	}
	if len(user.Password) < 4 {
		errs = append(err, "Invalid password, Password should be more than 4 characters")
	}
	if len(user.Name) < 1 {
		errs = append(err, "Invalid name, please enter a name")
	}

	return errs
}

// dummy password reset.
func ValidatePasswordReset(resetPassword db.ResetPassword) (bool, string) {
	if len(resetPassword.Password) < 4 {
		return false, "Invalid password, password should be more than 4 characters"
	}
	if resetPassword.Password != resetPassword.ConfirmPassword {
		return false, "Password reset failed, passwords must match"
	}
	return true, ""
}
