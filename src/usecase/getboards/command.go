package getboards

import (
	"github.com/google/uuid"
)

type GetBoardsCommand struct {
	UserID uuid.UUID
}

func NewGetBoardsCommand(userID string) (GetBoardsCommand, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return GetBoardsCommand{}, ErrInvalidUserID
	}

	return GetBoardsCommand{
		UserID: uid,
	}, nil
}
