package createboard

import "errors"

var (
	ErrInvalidName      = errors.New("name must be between 1 and 100 characters")
	ErrInvalidShortName = errors.New("short name must be 2–10 characters and contain only letters, numbers, hyphens or underscores")
	ErrEmptyName        = errors.New("name cannot be empty")
	ErrBoardIsExists    = errors.New("board with this shortname already exists")
	ErrNewBoardFailed   = errors.New("failed to create new board")
	ErrCreateBoard      = errors.New("unknown error creating board")
)
