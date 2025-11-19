package createboard

import "regexp"

// todo команда, без имени паккета. во всех командах
type CreateBoardCommand struct {
	Name      string
	ShortName string
}

var (
	ShortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewCreateBoardCommand(name string, shortName string) (CreateBoardCommand, error) {
	if name == "" {
		return CreateBoardCommand{}, EmptyNameErr
	}
	if len(name) > 100 {
		return CreateBoardCommand{}, InvalidNameErr
	}

	if shortName == "" {
		return CreateBoardCommand{}, InvalidShortNameErr
	}
	if !ShortNameRegex.MatchString(shortName) {
		return CreateBoardCommand{}, InvalidShortNameErr
	}

	return CreateBoardCommand{
		Name:      name,
		ShortName: shortName,
	}, nil
}
