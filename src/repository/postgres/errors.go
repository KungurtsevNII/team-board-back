package postgres

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmptyInput            = errors.New("empty input")
	ErrTaskNotFoundOrDeleted = errors.New("task not found or has been deleted")
)
