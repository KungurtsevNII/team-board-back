package deletetask

import (
	"errors"
)

var (
	ErrDeleteTaskUnknown = errors.New("unknown error deletion task")
	ErrInvalidTaskID  = errors.New("invalid task id")
)
