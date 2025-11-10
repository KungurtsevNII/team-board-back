package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r Repository) CheckBoard(id string) bool {
	uid, err := uuid.Parse(id)
	if err != nil {
		return false
	}

	tx, err := r.pool.BeginTx(context.TODO(), pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return false
	}
	defer tx.Rollback(context.TODO())

	var exists bool
	err = tx.QueryRow(context.TODO(),
		`SELECT EXISTS (SELECT 1 FROM boards WHERE id = $1 AND deleted_at IS NULL)`,
		uid,
	).Scan(&exists)

	if err != nil {
		return false
	}

	_ = tx.Commit(context.TODO())

	return exists
}
