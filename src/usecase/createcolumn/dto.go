package createcolumn

type CreateColumnCommand struct {
	Title   string
	BoardID int64
}

func NewCreateColumnCommand(title string, boardID int64) (CreateColumnCommand, error) {
	// todo validation

	return CreateColumnCommand{
		Title:   title,
		BoardID: boardID,
	}, nil
}

type GetTaskStatusQuery struct {
	TaskID string
}
