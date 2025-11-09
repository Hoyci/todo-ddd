package domain

import "context"

type Tx interface {
	Commit() error
	Rollback() error
}

type TxManager interface {
	Do(ctx context.Context, fn func(tx Tx) error) error
}
