package deleteboard

import "errors"

var (
	ErrBoardIdEmpty     = errors.New("board id is empty")
	ErrBoardIdInvalid   = errors.New("board id is invalid")
	ErrBoardDoesntExist = errors.New("board doesn't exist")
	ErrBoardDeleteUnknown = errors.New("board delete unknown error")
)
