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

	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
	"github.com/hoyci/todo-ddd/pkg/usecase"
)

type CreateTaskInput struct {
	Title       string
	Description string
	Priority    valueobject.Priority
	UserID      string
}

type CreateTaskOutput struct {
	*domainTask.Task
}

type CreateTaskUseCase struct {
	TxManager domain.TxManager
	UserRepo  domainUser.UserRepository
	TaskRepo  domainTask.TaskRepository
}

func NewCreateTaskUseCase(txManager domain.TxManager, userRepo domainUser.UserRepository, taskRepo domainTask.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		TxManager: txManager,
		UserRepo:  userRepo,
		TaskRepo:  taskRepo,
	}
}

func (uc *CreateTaskUseCase) Execute(input CreateTaskInput) (*CreateTaskOutput, error) {
	var output *CreateTaskOutput
	err := uc.TxManager.Do(context.Background(), func(tx domain.Tx) error {
		sqliteTx, ok := tx.(*sqlite.SQLiteTx)
		if !ok {
			return fmt.Errorf("unexpected tx type: %T", tx)
		}

		txUserRepo := uc.UserRepo.(*sqlite.SQLiteUserRepository).WithTx(sqliteTx.Tx())
		txTaskRepo := uc.TaskRepo.(*sqlite.SQLiteTaskRepository).WithTx(sqliteTx.Tx())

		user, err := txUserRepo.FindByID(input.UserID)
		if errors.Is(err, sql.ErrNoRows) || user.DeletedAt != nil {
			return usecase.ErrUserNotFoundOrDeleted
		}

		task, err := domainTask.NewTask(input.Title, input.Description, user.ID, input.Priority)
		if err != nil {
			return err
		}
		if err := txTaskRepo.Save(task); err != nil {
			slog.Error("error trying to save task", "err", err)
			return usecase.ErrTaskSaveFailed
		}

		output = &CreateTaskOutput{task}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}
