package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/hoyci/todo-ddd/pkg/domain"
	taskDomain "github.com/hoyci/todo-ddd/pkg/domain/task"
	userDomain "github.com/hoyci/todo-ddd/pkg/domain/user"
)

type sqliteWork struct {
	tx       *sql.Tx
	userRepo *SQLiteUserRepository
	taskRepo *SQLiteTaskRepository
}

func (w *sqliteWork) UserRepo() userDomain.UserRepository { return w.userRepo }
func (w *sqliteWork) TaskRepo() taskDomain.TaskRepository { return w.taskRepo }

type SQLiteUnitOfWork struct {
	db *sql.DB
}

func NewSQLiteUnitOfWork(db *sql.DB) domain.UnitOfWork {
	return &SQLiteUnitOfWork{db: db}
}

func (uow *SQLiteUnitOfWork) Execute(ctx context.Context, fn func(work domain.Work) error) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("error on begin tx", "error", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	work := &sqliteWork{
		tx:       tx,
		userRepo: NewSQLiteUserRepository(uow.db).WithTx(tx),
		taskRepo: NewSQLiteTaskRepository(uow.db).WithTx(tx),
	}

	if err := fn(work); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
