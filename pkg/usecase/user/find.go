package user

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/user"
)

type FindUserInput struct {
	ID string
}

type FindUserOutput struct {
	User *domain.User
}

type FindUserUseCase struct {
	UserRepo domain.UserRepository
}

func (uc *FindUserUseCase) Execute(input FindUserInput) (*FindUserOutput, error) {
	user, err := uc.UserRepo.FindByID(input.ID)
	if err != nil {
		slog.Error("error finding user by id", "id", input.ID)
		return nil, err
	}
	return &FindUserOutput{User: user}, nil
}
