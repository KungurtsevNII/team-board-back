package postgres

import (
	"time"

	"github.com/google/uuid"
)

type ColumnRecord struct {
	ID        uuid.UUID  `db:"id" goqu:"skipupdate"`
	BoardID   uuid.UUID  `db:"board_id"`
	Name      string     `db:"name"`
	OrderNum  int64      `db:"order_num"`
	CreatedAt time.Time  `db:"created_at" goqu:"skipupdate"`
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

type TaskSearchRecord struct {
	ID             uuid.UUID  `db:"tasks.id"`
	BoardID        uuid.UUID  `db:"tasks.board_id"`
	BoardName      string     `db:"boards.name"`
	BoardShortName string     `db:"boards.short_name"`
	ColumnName     string     `db:"columns.name"`
	ColumnID       uuid.UUID  `db:"tasks.column_id"`
	Number         int64      `db:"tasks.number"`
	Title          string     `db:"tasks.title"`
	CreatedAt      time.Time  `db:"tasks.created_at"`
	UpdatedAt      time.Time  `db:"tasks.updated_at"`
	DeletedAt      *time.Time `db:"tasks.deleted_at"`
}

type TaskSearchRecords []TaskSearchRecord //Для поинтера в маппинге

type BoardRecord struct {
	ID        uuid.UUID  `db:"id" goqu:"skipupdate"`
	Name      string     `db:"name"`
	ShortName string    `db:"short_name"`
	CreatedAt time.Time  `db:"created_at" goqu:"skipupdate"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}