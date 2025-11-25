package getboard

import "errors"

var (
	ErrBoardIsNotExists = errors.New("board is not exists")
	ErrInvalidID        = errors.New("invalid id")
	ErrBoardNotFound    = errors.New("board not found")
)
