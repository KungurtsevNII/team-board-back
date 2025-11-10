package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/jackc/pgx/v5"
)

func (r Repository) CreateBoard(board domain.Board) (string, error) {
	tx, err := r.pool.BeginTx(context.TODO(), pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return "", err
	}
	defer tx.Rollback(context.TODO())

	const sql = `
		INSERT INTO boards (id, name, short_name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(context.TODO(), sql,
		board.ID,
		board.Name,
		board.ShortName,
		board.CreatedAt,
		board.UpdatedAt,
		board.DeletedAt,
	)
	if err != nil {
		return "", err
	}
	if err = tx.Commit(context.TODO()); err != nil {
		return "", err
	}

	return board.ID.String(), nil
}
