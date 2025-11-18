package gettask

import (
	"errors"
)

var(
	ErrTaskNotFound = errors.New("task not found")
	ErrGetTaskUnknown = errors.New("unknown error getting task")
	ErrInvalidTaskID = errors.New("invalid task id")
)