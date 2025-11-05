package usecase

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/task"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
)

type UpdateTaskInput struct {
	TaskID      string
	Title       string
	Description string
	Priority    valueobject.Priority
	UserID      string
}

type UpdateTaskOutput struct {
	domain.Task
}

type UpdateTaskUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *UpdateTaskUseCase) Execute(input UpdateTaskInput) (*UpdateTaskOutput, error) {
	task, err := uc.TaskRepo.FindByID(input.TaskID, input.UserID)
	if err != nil {
		slog.Error("error trying to find task by id", "taskID", input.TaskID)
		return nil, err
	}

	task.Update(input.Title, input.Description, input.Priority)

	err = uc.TaskRepo.Update(task)
	if err != nil {
		slog.Error("error trying to update task", "taskID", task.ID)
		return nil, err
	}
	return &UpdateTaskOutput{*task}, nil
}
