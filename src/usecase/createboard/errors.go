package createboard

import "errors"

// todo в один блок и привести к формату ошибок. Err....

// имя слишком короткое или слишком длинное
var InvalidNameErr = errors.New("name must be between 1 and 100 characters")

// короткое имя не соответствует формату
var InvalidShortNameErr = errors.New("short name must be 2–10 characters and contain only letters, numbers, hyphens or underscores")

// имя не может быть пустым
var EmptyNameErr = errors.New("name cannot be empty")

// доска уже существует
var BoardIsExistsErr = errors.New("board with this shortname already exists")
