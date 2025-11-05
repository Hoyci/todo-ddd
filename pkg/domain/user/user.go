package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewUser(name, email string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}
}

func (t *User) Delete() {
	now := time.Now()
	t.UpdatedAt = &now
	t.DeletedAt = &now
}
