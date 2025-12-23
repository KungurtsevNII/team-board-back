package gettask

import (
	"github.com/google/uuid"
)

type GetTaskQuery struct {
	TaskID uuid.UUID 
}

func NewQuery(userID string) (GetTaskQuery, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return GetTaskQuery{}, ErrInvalidTaskID
	}

	return GetTaskQuery{
		TaskID: uid,
	}, nil
}
