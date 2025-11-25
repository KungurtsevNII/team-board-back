package createboard

import "errors"

var (
	ErrInvalidName      = errors.New("name must be between 1 and 100 characters")
	ErrInvalidShortName = errors.New("short name must be 2–10 characters and contain only letters, numbers, hyphens or underscores")
	ErrBoardIsExists    = errors.New("board with this shortname already exists")
	ErrCreateBoard      = errors.New("unknown error creating board")
	ErrInvalidColumnName = errors.New("column name must be between 1 and 100 characters")
	ErrGetLastOrderNumUnknown = errors.New("failed to get last order num")
	ErrCreateColumnUnknown    = errors.New("failed to create column")
	ErrValidationFailed       = errors.New("validation failed")
)
