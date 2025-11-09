package sqlite

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/hoyci/todo-ddd/pkg/domain"
)

type SQLiteTx struct {
	tx *sql.Tx
}

func (t *SQLiteTx) Commit() error   { return t.tx.Commit() }
func (t *SQLiteTx) Rollback() error { return t.tx.Rollback() }
func (t *SQLiteTx) Tx() *sql.Tx {
	return t.tx
}

type SQLiteTxManager struct {
	db *sql.DB
}

func NewSQLiteTxManager(db *sql.DB) *SQLiteTxManager {
	return &SQLiteTxManager{db: db}
}

func (m *SQLiteTxManager) Do(ctx context.Context, fn func(tx domain.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("error on begin tx", "error", err)
		return err
	}

	sqliteTx := &SQLiteTx{tx: tx}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(sqliteTx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
