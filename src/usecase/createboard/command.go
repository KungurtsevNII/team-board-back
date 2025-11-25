package createboard

import "regexp"

type Command struct {
	Name       string
	ShortName  string
}

var (
	ShortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewCommand(name string, shortName string) (Command, error) {
	if len(name) > 100 || name == "" {
		return Command{}, ErrInvalidName
	}
	if shortName == "" || !ShortNameRegex.MatchString(shortName) {
		return Command{}, ErrInvalidShortName
	}

	return Command{
		Name:      name,
		ShortName: shortName,
	}, nil
}
