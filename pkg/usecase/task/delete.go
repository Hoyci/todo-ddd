package usecase

import (
	"log/slog"

	"github.com/hoyci/todo-ddd/pkg/domain"
)

type DeleteTaskInput struct {
	ID string
}

type DeleteTaskOutput struct {
	ID string
}

type DeleteTaskUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *DeleteTaskUseCase) Execute(input DeleteTaskInput) (*DeleteTaskOutput, error) {
	task, err := uc.TaskRepo.FindByID(input.ID)
	if err != nil {
		slog.Error("error trying to find task by id", "taskID", input.ID)
		return nil, err
	}

	task.Delete()

	err = uc.TaskRepo.Delete(task.ID, *task.DeletedAt)
	if err != nil {
		slog.Error("error trying to delete task by id", "taskID", task.ID)
		return nil, err
	}
	return &DeleteTaskOutput{ID: task.ID}, nil
}
