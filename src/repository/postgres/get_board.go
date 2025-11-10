package postgres

import (
	"context"
	"errors"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r Repository) GetBoard(ID string) (domain.Board, error) {
	uid, err := uuid.Parse(ID)
	if err != nil {
		return domain.Board{}, domain.ErrInvalidID
	}

	var board domain.Board
	tx, err := r.pool.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return domain.Board{}, err
	}
	defer tx.Rollback(context.TODO())

	row := tx.QueryRow(context.TODO(),
		`SELECT id, name, short_name, created_at, updated_at, deleted_at
		FROM boards
		WHERE id = $1`,
		uid,
	)

	err = row.Scan(
		&board.ID,
		&board.Name,
		&board.ShortName,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Board{}, domain.ErrBoardNotFound
		}
		return domain.Board{}, err
	}

	return board, nil
}
