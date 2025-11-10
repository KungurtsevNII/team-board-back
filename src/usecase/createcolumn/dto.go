package createcolumn

type CreateColumnCommand struct {
	Title   string
	BoardID string
}

func NewCreateColumnCommand(title string, boardID string) (CreateColumnCommand, error) {
	// todo validation

	return CreateColumnCommand{
		Title:   title,
		BoardID: boardID,
	}, nil
}

type GetTaskStatusQuery struct {
	TaskID string
}
