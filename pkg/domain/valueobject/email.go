package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var ErrInvalidEmail = errors.New("invalid email format")

func NewEmail(raw string) (Email, error) {
	email := strings.TrimSpace(strings.ToLower(raw))

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: email}, nil
}

func (e Email) String() string {
	return e.value
}
