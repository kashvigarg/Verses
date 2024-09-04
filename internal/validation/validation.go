package validate

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)

	if !match {
		return errors.New("email is not valid")
	}

	return nil
}

func ValidateUsername(username string) error {
	regexstring := `^[a-zA-Z_][a-zA-Z0-9._%+-]{0,8}$`

	is_valid, _ := regexp.MatchString(regexstring, username)

	if !is_valid {
		return errors.New("username is not valid")
	}

	return nil
}

// since regexp doesn't suport lookahead we need to use multiple operations

func ValidatePassword(passwd string) error {

	lowercase := regexp.MustCompile(`[a-z]`).MatchString(passwd)

	uppercase := regexp.MustCompile(`[A-Z]`).MatchString(passwd)

	digit := regexp.MustCompile(`\d`).MatchString(passwd)

	allowedChars := regexp.MustCompile(`^[a-zA-Z\d@$#&_]{8,}$`).MatchString(passwd)

	if !lowercase || !uppercase || !digit || !allowedChars {
		return errors.New("password is not valid")
	}
	return nil
}
