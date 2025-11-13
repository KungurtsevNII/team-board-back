package createcolumn

import (
	"errors"
)

var (
	// ErrColumnIsExistsErr = errors.New("column is exists")
	ErrBoardIsNotExists = errors.New("board is not exists")
	ErrInvalidUUID = errors.New("invalid uuid")
	ErrValidationFailed = errors.New("validation failed")
	ErrGetLastOrderNumUnknown = errors.New("failed to get last order num")
	ErrCreateColumnUnknown = errors.New("failed to create column")
)
