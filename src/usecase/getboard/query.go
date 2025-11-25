package getboard

import "github.com/google/uuid"

type Query struct {
	ID uuid.UUID
}

func NewQuery(ID string) (Query, error) {
	uid, err := uuid.Parse(ID)
	if err != nil {
		return Query{}, err
	}
	return Query{
		ID: uid,
	}, nil
}
