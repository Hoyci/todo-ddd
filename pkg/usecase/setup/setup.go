package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hoyci/todo-ddd/internal/adapters/db/sqlite"
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
	TxManager domain.TxManager
	UserRepo  domainUser.UserRepository
	TaskRepo  domainTask.TaskRepository
}

func NewOnboardingUseCase(txManager domain.TxManager, userRepo domainUser.UserRepository, taskRepo domainTask.TaskRepository) *SetupOnboardingUseCase {
	return &SetupOnboardingUseCase{
		TxManager: txManager,
		UserRepo:  userRepo,
		TaskRepo:  taskRepo,
	}
}

func (uc *SetupOnboardingUseCase) Execute(input SetupOnboardingInput) error {
	err := uc.TxManager.Do(context.Background(), func(tx domain.Tx) error {
		sqliteTx, ok := tx.(*sqlite.SQLiteTx)
		if !ok {
			return fmt.Errorf("unexpected tx type: %T", tx)
		}

		txUserRepo := uc.UserRepo.(*sqlite.SQLiteUserRepository).WithTx(sqliteTx.Tx())
		txTaskRepo := uc.TaskRepo.(*sqlite.SQLiteTaskRepository).WithTx(sqliteTx.Tx())

		userExists, err := txUserRepo.FindByEmail(input.Email)
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
		if err = txUserRepo.Save(*user); err != nil {
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
		if err = txTaskRepo.Save(task); err != nil {
			return usecase.ErrTaskSaveFailed
		}

		return nil
	})
	return err
}
