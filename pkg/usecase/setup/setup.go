package usecase

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/hoyci/todo-ddd/internal/adapters/db/sqlite"
	domainTask "github.com/hoyci/todo-ddd/pkg/domain/task"
	domainUser "github.com/hoyci/todo-ddd/pkg/domain/user"
	"github.com/hoyci/todo-ddd/pkg/usecase"
)

type SetupOnboardingInput struct {
	Name  string
	Email string
}

type SetupOnboardingUseCase struct {
	DB *sql.DB
}

func (uc *SetupOnboardingUseCase) Execute(input SetupOnboardingInput) error {
	tx, err := uc.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	txUserRepo := sqlite.NewSQLiteUserRepositoryWithTx(tx)
	txTaskRepo := sqlite.NewSQLiteTaskRepositoryWithTx(tx)

	userExists, err := txUserRepo.FindByEmail(input.Email)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("error finding user by email", "email", input.Email, "error", err)
		return usecase.ErrSearchingUserByEmail
	}

	if userExists != nil {
		return usecase.ErrUserAlreadyExists
	}

	user := domainUser.NewUser(input.Name, input.Email)
	if err = txUserRepo.Save(*user); err != nil {
		return usecase.ErrUserSaveFailed
	}

	task := domainTask.NewTask("Tarefa de boas-vindas", "Essa é uma tarefa criada automaticamente", user.ID, 1)
	if err = txTaskRepo.Save(task); err != nil {
		return usecase.ErrTaskSaveFailed
	}

	if err = tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "err", err)
		return fmt.Errorf("%w: %v", usecase.ErrTransactionCommitFailed, err)
	}

	return nil
}
