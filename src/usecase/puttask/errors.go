package puttask

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrTaskNotFound = errors.New("task not found")
	ErrPutTaskUnknown = errors.New("unknown error while putting task")
	ErrColumnNotFound = errors.New("column not found")
)
