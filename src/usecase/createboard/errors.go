package createboard

import "errors"

// имя слишком короткое или слишком длинное
var ErrInvalidName = errors.New("name must be between 1 and 100 characters")

// короткое имя не соответствует формату
var ErrInvalidShortName = errors.New("short name must be 2–10 characters and contain only letters, numbers, hyphens or underscores")

// имя не может быть пустым
var ErrEmptyName = errors.New("name cannot be empty")

// короткое имя не может быть пустым
var ErrEmptyShortName = errors.New("short name cannot be empty")

// доска уже существует
var ErrBoardIsExists = errors.New("board with this shortname already exists")
