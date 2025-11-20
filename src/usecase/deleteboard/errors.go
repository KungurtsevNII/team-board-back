package deleteboard

import "errors"

var (
	ErrBoardIdEmpty   = errors.New("board id is empty")
	ErrBoardIdInvalid = errors.New("board id is invalid")
)
