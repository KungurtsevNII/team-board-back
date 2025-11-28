package movetask

import (
	"errors"
)

var (
	ErrInvalidUUID      = errors.New("invalid uuid")
	ErrValidationFailed = errors.New("validation failed")
	ErrTaskNotFound     = errors.New("task not found")
	ErrMoveTaskUnknown  = errors.New("failed to move task")
)
