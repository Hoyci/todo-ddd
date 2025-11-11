package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/hoyci/todo-ddd/pkg/domain"
	domainTask "github.com/hoyci/todo-ddd/pkg/domain/task"
	domainUser "github.com/hoyci/todo-ddd/pkg/domain/user"
	"github.com/hoyci/todo-ddd/pkg/usecase"
)

type SetupOnboardingInput struct {
	Name  string
	Email string
}

type SetupOnboardingUseCase struct {
	UoW domain.UnitOfWork
}

func (uc *SetupOnboardingUseCase) Execute(input SetupOnboardingInput) error {
	return uc.UoW.Execute(context.Background(), func(work domain.Work) error {
		userRepo := work.UserRepo()
		taskRepo := work.TaskRepo()

		userExists, err := userRepo.FindByEmail(input.Email)
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("unexpected error", "error", err)
			return usecase.ErrUnknown
		}

		if userExists != nil {
			if !userExists.DeletedAt.IsZero() {
				return usecase.ErrUserNotFoundOrDeleted
			}
			return usecase.ErrUserAlreadyExists
		}

		user, err := domainUser.NewUser(input.Name, input.Email)
		if err != nil {
			return err
		}
		if err = userRepo.Save(*user); err != nil {
			return usecase.ErrUserSaveFailed
		}

		task, err := domainTask.NewTask(
			"Tarefa de boas-vindas",
			"Essa é uma tarefa criada automaticamente",
			user.ID,
			1,
		)
		if err != nil {
			return err
		}
		if err = taskRepo.Save(task); err != nil {
			return usecase.ErrTaskSaveFailed
		}

		return nil
	})
}
