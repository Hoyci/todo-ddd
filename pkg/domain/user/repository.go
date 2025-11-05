package domain

import (
	"time"
)

type UserRepository interface {
	Save(user User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	List() ([]*User, error)
	Update(user User) error
	Delete(id string, timestamp time.Time) error
}
