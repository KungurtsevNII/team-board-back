package deletetask

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repo interface {
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd Command) error{
	//В удалениях вроде не возвращают ошибку при not found, поэтому и я не решил не возвращать
	err := uc.repo.DeleteTask(ctx, cmd.TaskID)
	if err != nil {
		return errors.Wrap(ErrDeleteTaskUnknown, err.Error())
	}
	return nil
}
