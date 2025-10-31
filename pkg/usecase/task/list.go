package usecase

import (
	"log/slog"

	"github.com/hoyci/todo-ddd/pkg/domain"
)

type ListTaskOutput struct {
	*domain.Task
}

type ListTaskUseCase struct {
	TaskRepo domain.TaskRepository
}

func (uc *ListTaskUseCase) Execute() ([]ListTaskOutput, error) {
	tasks, err := uc.TaskRepo.List()
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
