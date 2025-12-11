package postgres

import (
	"context"

	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PgxPoolIface interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Begin(ctx context.Context) (pgx.Tx, error)

	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	Close()
	Ping(ctx context.Context) error
}

type Repository struct {
	pool PgxPoolIface
}

func New(storagePath string) (*Repository, error) {
	const op = "storage.postgresql.New"

	log := slog.Default()
	log.With("storagePath", storagePath, "op", op).Info("connecting to postgres")
	pool, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (s *Repository) Close() {
	s.pool.Close()
}
