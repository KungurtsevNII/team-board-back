package postgres

import "time"

type ColumnRecord struct {
	ID        string     `db:"id"`
	BoardID   string     `db:"board_id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"primaryKey"`
	UpdatedAt time.Time  `db:"primaryKey"`
}