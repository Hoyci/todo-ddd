package usecase

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/task"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
)

type UpdateTaskStatusInput struct {
	TaskID string
	Status valueobject.Status
	UserID string
}

type UpdateTaskStatusOutput struct {
	domain.Task
}

type UpdateTaskStatusUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *UpdateTaskStatusUseCase) Execute(input UpdateTaskStatusInput) (*UpdateTaskStatusOutput, error) {
	task, err := uc.TaskRepo.FindByID(input.TaskID, input.UserID)
	if err != nil {
		slog.Error("error trying to find task by id", "taskID", input.TaskID)
		return nil, err
	}

	switch input.Status {
	case valueobject.StatusNew:
		task.SetAsNew()
	case valueobject.StatusInProgress:
		task.SetInProgress()
	case valueobject.StatusCompleted:
		task.SetCompleted()
	}

	err = uc.TaskRepo.Update(task)
	if err != nil {
		slog.Error("error trying to update task status", "taskID", task.ID, "taskStatus", task.Status)
		return nil, err
	}
	return &UpdateTaskStatusOutput{*task}, nil
}
