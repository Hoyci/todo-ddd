package usecase

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/task"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
)

type CreateTaskInput struct {
	Title       string
	Description string
	Priority    valueobject.Priority
	UserID      string
}

type CreateTaskOutput struct {
	*domain.Task
}

type CreateTaskUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *CreateTaskUseCase) Execute(input CreateTaskInput) (*CreateTaskOutput, error) {
	task := domain.NewTask(input.Title, input.Description, input.UserID, input.Priority)
	err := uc.TaskRepo.Save(task)
	if err != nil {
		slog.Error("error trying to save task")
		return nil, err
	}
	return &CreateTaskOutput{task}, nil
}
