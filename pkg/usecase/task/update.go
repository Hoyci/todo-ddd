package usecase

import (
	"log/slog"

	"github.com/hoyci/todo-ddd/pkg/domain"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
)

type UpdateTaskInput struct {
	ID          string
	Title       string
	Description string
	Priority    valueobject.Priority
}

type UpdateTaskOutput struct {
	domain.Task
}

type UpdateTaskUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *UpdateTaskUseCase) Execute(input UpdateTaskInput) (*UpdateTaskOutput, error) {
	task, err := uc.TaskRepo.FindByID(input.ID)
	if err != nil {
		slog.Error("error trying to find task by id", "taskID", input.ID)
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
