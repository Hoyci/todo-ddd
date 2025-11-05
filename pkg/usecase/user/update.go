package user

import (
	"log/slog"
	"time"

	domain "github.com/hoyci/todo-ddd/pkg/domain/user"
)

type UpdateUserInput struct {
	ID    string
	Name  string
	Email string
}

type UpdateUserOutput struct {
	User *domain.User
}

type UpdateUserUseCase struct {
	UserRepo domain.UserRepository
}

func (uc *UpdateUserUseCase) Execute(input UpdateUserInput) (*UpdateUserOutput, error) {
	user, err := uc.UserRepo.FindByID(input.ID)
	if err != nil {
		slog.Error("error finding user to update", "id", input.ID)
		return nil, err
	}

	now := time.Now()
	user.Name = input.Name
	user.Email = input.Email
	user.UpdatedAt = &now

	if err := uc.UserRepo.Update(*user); err != nil {
		slog.Error("error updating user", "id", input.ID)
		return nil, err
	}

	return &UpdateUserOutput{User: user}, nil
}
