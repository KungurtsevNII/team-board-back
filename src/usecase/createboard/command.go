package createboard

import "regexp"

type Command struct {
	Name       string
	ShortName  string
	ColumnName string
}

var (
	ShortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewCommand(name string, shortName string, columnName string) (Command, error) {
	if len(name) > 100 || name == "" {
		return Command{}, ErrInvalidName
	}
	if shortName == "" || !ShortNameRegex.MatchString(shortName) {
		return Command{}, ErrInvalidShortName
	}
	if len(columnName) > 100 || columnName == "" {
		return Command{}, ErrInvalidColumnName
	}

	return Command{
		Name:      name,
		ShortName: shortName,
		ColumnName: columnName,
	}, nil
}
