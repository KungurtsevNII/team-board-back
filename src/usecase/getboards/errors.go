package getboards

import "errors"

var (
	ErrInvalidUserID = errors.New("invalid user id")
	ErrGetBoards     = errors.New("unknown error getting boards")
)
