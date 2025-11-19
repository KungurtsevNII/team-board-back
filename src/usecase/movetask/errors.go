package movetask

import (
	"errors"
)

var (
	ErrInvalidUUID      = errors.New("invalid uuid")
	ErrValidationFailed = errors.New("validation failed")
	ErrTaskNotFound     = errors.New("task not found")
	ErrColumnNotFound   = errors.New("column not found")
	ErrColumnNotInBoard = errors.New("column does not belong to task's board")
	ErrMoveTaskUnknown  = errors.New("failed to move task")
)
