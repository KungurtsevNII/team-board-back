package searchtasks

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrSearchTasks = errors.New("search tasks failed")
)
