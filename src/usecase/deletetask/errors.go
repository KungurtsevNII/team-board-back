package deletetask

import (
	"errors"
)

var (
	ErrDeleteTaskUnknown = errors.New("unknown error deletion task")
	ErrInvalidTaskID  = errors.New("invalid task id")
	ErrTaskNotFound = errors.New("task not found")
	ErrGetTaskUnknown = errors.New("unknown error getting task")
)
