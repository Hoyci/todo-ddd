package user

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/user"
)

type CreateUserInput struct {
	Name  string
	Email string
}

type CreateUserOutput struct {
	User *domain.User
}

type CreateUserUseCase struct {
	UserRepo domain.UserRepository
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	user, err := domain.NewUser(input.Name, input.Email)
	if err != nil {
		return nil, err
	}

	if err := uc.UserRepo.Save(*user); err != nil {
		slog.Error("error saving user", "email", input.Email)
		return nil, err
	}

	return &CreateUserOutput{User: user}, nil
}
