package user

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/user"
)

type DeleteUserInput struct {
	ID string
}

type DeleteUserUseCase struct {
	UserRepo domain.UserRepository
}

func (uc *DeleteUserUseCase) Execute(input DeleteUserInput) error {
	user, err := uc.UserRepo.FindByID(input.ID)
	if err != nil {
		slog.Error("error finding user to delete", "id", input.ID)
		return err
	}

	user.Delete()

	if err := uc.UserRepo.Delete(user.ID, *user.DeletedAt); err != nil {
		slog.Error("error deleting user", "id", input.ID)
		return err
	}
	return nil
}
