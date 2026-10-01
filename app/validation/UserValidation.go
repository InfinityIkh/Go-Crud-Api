package validation

import (
	"errors"
	"regexp"
)

func Validate(username string ,email string) error {
	if username == "" || len(username) < 4 {
		return errors.New("username is required and must contains 4 caractere at least")
	}
	if email == "" {
		return errors.New("email is required")
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+@[a-z]+\.[a-z]+$`)
	if !re.MatchString(email) {
		return errors.New("email is not validated")
	}
	return nil
}