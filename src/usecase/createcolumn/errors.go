package createcolumn

import (
	"errors"
)

var (
	// ErrColumnIsExistsErr = errors.New("column is exists")
	ErrBoardIsNotExistsErr = errors.New("board is not exists")
	ErrInvalidUUID = errors.New("invalid uuid")
	ErrValidationFailed = errors.New("validation failed")
	ErrGetLastOrderNumErr = errors.New("failed to get last order num")
	ErrCreateColumnErr = errors.New("failed to create column")
)
