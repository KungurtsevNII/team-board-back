package deletecolumn

import (
	"errors"
)

var (
	ErrDeleteColumnUnknown = errors.New("unknown error deletion column")
	ErrInvalidColumnID  = errors.New("invalid column id")
	ErrColumnNotFound = errors.New("column not found")
	ErrGetColumnUnknown = errors.New("unknown error getting column")
	ErrCheckColumnIsEmptyUnknown = errors.New("unknown error checking if column is empty")
	ErrColumnNotEmpty = errors.New("column is not empty")
)
