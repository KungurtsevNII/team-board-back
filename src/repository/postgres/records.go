package postgres

import (
	"time"

	"github.com/google/uuid"
)

type ColumnRecord struct {
	ID        uuid.UUID  `db:"id"`
	BoardID   uuid.UUID  `db:"board_id"`
	Name      string     `db:"name"`
	OrderNum  int64      `db:"order_num"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

type TaskRecord struct {
	ID          uuid.UUID  `db:"id"`
	BoardID     uuid.UUID  `db:"board_id"`
	ColumnID    uuid.UUID  `db:"column_id"`
	Number      int64      `db:"number"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	Tags        []string   `db:"tags"`
	Checklists  []byte     `db:"checklists"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

