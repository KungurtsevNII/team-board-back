package getboard

import "github.com/google/uuid"

type Query struct {
	ID string
}

func NewQuery(ID string) (Query, error) {
	_, err := uuid.Parse(ID)
	if err != nil {
		return Query{}, err
	}
	return Query{
		ID: ID,
	}, nil
}
