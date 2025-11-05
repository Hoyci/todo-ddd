package usecase

import (
	"log/slog"

	domain "github.com/hoyci/todo-ddd/pkg/domain/task"
)

type ListTaskOutput struct {
	*domain.Task
}

type ListTaskUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *ListTaskUseCase) Execute(userID string) ([]ListTaskOutput, error) {
	tasks, err := uc.TaskRepo.List(userID)
	if err != nil {
		slog.Error("error trying to list tasks")
		return nil, err
	}

	var output []ListTaskOutput
	for _, t := range tasks {
		output = append(output, ListTaskOutput{Task: t})
	}

	return output, nil
}
