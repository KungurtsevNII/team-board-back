package createboard

import "regexp"

type CreateBoardCommand struct {
	Name      string
	ShortName string
}

var (
	ShortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewCreateBoardCommand(name string, shortName string) (CreateBoardCommand, error) {
	if name == "" {
		return CreateBoardCommand{}, ErrEmptyName
	}
	if len(name) > 100 {
		return CreateBoardCommand{}, ErrInvalidName
	}

	if shortName == "" {
		return CreateBoardCommand{}, ErrInvalidShortName
	}
	if !ShortNameRegex.MatchString(shortName) {
		return CreateBoardCommand{}, ErrInvalidShortName
	}

	return CreateBoardCommand{
		Name:      name,
		ShortName: shortName,
	}, nil
}
