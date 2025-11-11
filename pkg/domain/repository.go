package domain

import (
	"context"

	taskDomain "github.com/hoyci/todo-ddd/pkg/domain/task"
	userDomain "github.com/hoyci/todo-ddd/pkg/domain/user"
)

type Work interface {
	UserRepo() userDomain.UserRepository
	TaskRepo() taskDomain.TaskRepository
}

type UnitOfWork interface {
	Execute(ctx context.Context, fn func(work Work) error) error
}
