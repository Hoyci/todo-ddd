package usecase

import "errors"

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrUserSaveFailed          = errors.New("failed to save user")
	ErrSearchingUserByEmail    = errors.New("failed to search user by email")
	ErrSearchingUserByID       = errors.New("failed to search user by ID")
	ErrUserNotFoundOrDeleted   = errors.New("user not found or deleted")
	ErrUserNotFound            = errors.New("user not found")
	ErrTaskSaveFailed          = errors.New("failed to save task")
	ErrTransactionCommitFailed = errors.New("failed to commit transaction")
	ErrUnknown                 = errors.New("unexpected error")
)
