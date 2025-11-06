package usecase

import "errors"

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrUserSaveFailed          = errors.New("failed to save user")
	ErrSearchingUserByEmail    = errors.New("error seaching user by email")
	ErrUserNotFound            = errors.New("user not found")
	ErrTaskSaveFailed          = errors.New("failed to save task")
	ErrTransactionCommitFailed = errors.New("failed to commit transaction")
)
