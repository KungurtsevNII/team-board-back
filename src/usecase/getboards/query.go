package getboards

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Query struct {
	UserID uuid.UUID
}

func NewQuery(userID string) (Query, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return Query{}, errors.Wrap(ErrInvalidUserID, err.Error())
	}

	return Query{
		UserID: uid,
	}, nil
}
