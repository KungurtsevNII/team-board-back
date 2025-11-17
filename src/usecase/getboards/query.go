package getboards

import (
	"github.com/google/uuid"
)

type GetBoardsQuery struct {
	UserID uuid.UUID
}

func NewGetBoardsQuery(userID string) (GetBoardsQuery, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return GetBoardsQuery{}, ErrInvalidUserID
	}

	return GetBoardsQuery{
		UserID: uid,
	}, nil
}
