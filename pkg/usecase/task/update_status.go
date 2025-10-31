package usecase

import (
	"log/slog"

	"github.com/hoyci/todo-ddd/pkg/domain"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
)

type UpdateTaskStatusInput struct {
	Title       string
	Description string
	Priority    valueobject.Priority
}

type UpdateTaskStatusOutput struct {
	domain.Task
}

type UpdateTaskStatusUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *UpdateTaskStatusUseCase) Execute(ID string, status valueobject.Status) (*UpdateTaskStatusOutput, error) {
	task, err := uc.TaskRepo.FindByID(ID)
	if err != nil {
		slog.Error("error trying to find task by id", "taskID", ID)
		return nil, err
	}

	switch status {
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
