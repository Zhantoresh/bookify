package validator

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

var phoneRegexp = regexp.MustCompile(`^\+?[0-9]{7,20}$`)

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "invalid input data"
	}
	return "invalid input data"
}

func (v *ValidationErrors) Add(field, errText string) {
	*v = append(*v, FieldError{Field: field, Error: errText})
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("must be a valid email address")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("must be at least 8 characters long")
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case 'A' <= r && r <= 'Z':
			hasUpper = true
		case 'a' <= r && r <= 'z':
			hasLower = true
		case '0' <= r && r <= '9':
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("must include upper, lower case letters and digits")
	}
	return nil
}

func ValidateRequired(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("is required")
	}
	return nil
}

func ValidatePhone(phone string) error {
	if phone == "" {
		return nil
	}
	if !phoneRegexp.MatchString(phone) {
		return errors.New("must be a valid phone number")
	}
	return nil
}
