package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewUser(name, email string) (*User, error) {
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     emailVO.String(),
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}, nil
}

func (t *User) Delete() {
	now := time.Now()
	t.UpdatedAt = &now
	t.DeletedAt = &now
}
