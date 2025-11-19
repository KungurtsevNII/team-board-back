package postgres

import "errors"

var (
	// todo выпилить всю эту хуйню
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyInput   = errors.New("empty input")
)
