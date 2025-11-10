package createboard

import "regexp"

type CreateBoardCommand struct {
	Name      string
	ShortName string
}

var (
	ShortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewCreateBoardCommand(name string, shortName string) (*CreateBoardCommand, error) {

	//validation
	if name == "" {
		return nil, EmptyNameErr
	}
	if len(name) > 100 {
		return nil, InvalidNameErr
	}

	if shortName == "" {
		return nil, EmptyShortNameErr
	}
	if !ShortNameRegex.MatchString(shortName) {
		return nil, InvalidShortNameErr
	}

	return &CreateBoardCommand{
		Name:      name,
		ShortName: shortName,
	}, nil
}
