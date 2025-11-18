package getboard

type GetBoardCommand struct {
	ID string
}

func NewGetBoardCommand(ID string) (GetBoardCommand, error) {
	return GetBoardCommand{
		ID: ID,
	}, nil
}
