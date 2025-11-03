package repository

import "errors"


var(
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyInput = errors.New("empty input")
)