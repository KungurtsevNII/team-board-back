package createboard

import "regexp"

// todo команда, без имени паккета. во всех командах
type Command struct {
	Name      string
	ShortName string
}

var (
	ShortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewCommand(name string, shortName string) (Command, error) {
	if name == "" {
		return Command{}, ErrEmptyName
	}
	if len(name) > 100 {
		return Command{}, ErrInvalidName
	}

	if shortName == "" {
		return Command{}, ErrInvalidShortName
	}
	if !ShortNameRegex.MatchString(shortName) {
		return Command{}, ErrInvalidShortName
	}

	return Command{
		Name:      name,
		ShortName: shortName,
	}, nil
}
