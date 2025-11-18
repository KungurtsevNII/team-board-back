package createtask

import (
	"errors"
)

var(
	ErrInvalidUUID = errors.New("invalid uuid")
	ErrValidationFailed = errors.New("validation failed")
	ErrCheckColumnInBoardFailed = errors.New("check column in board failed")
	ErrColumnOrBoardIsNotExists = errors.New("column or board is not exists")
	ErrGetLastNumberFailed = errors.New("get last number failed")
	ErrCreateTaskUnknown = errors.New("failed to create task")
)