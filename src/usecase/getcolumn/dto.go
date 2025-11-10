package getcolumn

type GetColumnCommand struct {
	ID string
}

func NewGetColumnCommand(ID string) (GetColumnCommand, error) {
	// todo validation

	return GetColumnCommand{
		ID: ID,
	}, nil
}

