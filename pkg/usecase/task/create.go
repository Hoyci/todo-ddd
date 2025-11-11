package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hoyci/todo-ddd/pkg/domain"
	domainTask "github.com/hoyci/todo-ddd/pkg/domain/task"

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
	UoW domain.UnitOfWork
}

func (uc *CreateTaskUseCase) Execute(input CreateTaskInput) (*CreateTaskOutput, error) {
	var output *CreateTaskOutput
	err := uc.UoW.Execute(context.Background(), func(work domain.Work) error {
		userRepo := work.UserRepo()
		taskRepo := work.TaskRepo()

		user, err := userRepo.FindByID(input.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return usecase.ErrUserNotFoundOrDeleted
			}
			return err
		}

		if user.DeletedAt != nil {
			return usecase.ErrUserNotFoundOrDeleted
		}

		task, err := domainTask.NewTask(input.Title, input.Description, user.ID, input.Priority)
		if err != nil {
			return err
		}

		if err := taskRepo.Save(task); err != nil {
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
