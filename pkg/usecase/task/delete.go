package usecase

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/task"
)

type DeleteTaskInput struct {
	TaskID string
	UserID string
}

type DeleteTaskOutput struct {
	ID string
}

type DeleteTaskUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *DeleteTaskUseCase) Execute(input DeleteTaskInput) (*DeleteTaskOutput, error) {
	task, err := uc.TaskRepo.FindByID(input.TaskID, input.UserID)
	if err != nil {
		slog.Error("error trying to find task by id", "taskID", input.TaskID)
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
