package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchTasks(t *testing.T) {
	now := time.Now()
	boardID := uuid.New()
	columnID := uuid.New()

	baseCols := []string{
		"tasks.id",
		"tasks.board_id",
		"boards.name",
		"boards.short_name",
		"columns.name",
		"tasks.column_id",
		"tasks.number",
		"tasks.title",
		"tasks.created_at",
		"tasks.updated_at",
		"tasks.deleted_at",
	}

	baseFromJoin := `SELECT .+ FROM "tasks" ` +
		`INNER JOIN "boards" ON \("tasks"\."board_id" = "boards"\."id"\) ` +
		`INNER JOIN "columns" ON \("tasks"\."column_id" = "columns"\."id"\) `

	baseSoftDeleteFilters := `WHERE .+` +
		`"tasks"\."deleted_at" IS NULL.+` +
		`"boards"\."deleted_at" IS NULL.+`

	tests := []struct {
		name        string
		tags        []string
		query       string
		limit       uint
		offset      uint
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expectedLen int
		expectedErr error
	}{
		{
			name:  "поиск задач без фильтров",
			tags:  []string{},
			query: "",
			limit: 10, offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows(baseCols).
					AddRow(uuid.New(), boardID, "Board 1", "B1", "Todo", columnID, int64(1), "Task 1", now, now, nil).
					AddRow(uuid.New(), boardID, "Board 1", "B1", "Todo", columnID, int64(2), "Task 2", now, now, nil)

				mock.ExpectQuery(
					baseFromJoin +
						baseSoftDeleteFilters +
						`ORDER BY "tasks"\."created_at" DESC LIMIT 10`,
				).WillReturnRows(rows)
			},
			expectedLen: 2,
		},
		{
			name:  "поиск задач по тегам",
			tags:  []string{"urgent", "bug"},
			query: "",
			limit: 5, offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows(baseCols).
					AddRow(uuid.New(), boardID, "Board 1", "B1", "Todo", columnID, int64(1), "Bug Task", now, now, nil)

				mock.ExpectQuery(
					baseFromJoin +
						`WHERE .+tags @>.+` + 
						`.+\"tasks\"\.\"deleted_at\" IS NULL.+\"boards\"\.\"deleted_at\" IS NULL.+` +
						`ORDER BY "tasks"\."created_at" DESC LIMIT 5`,
				).WillReturnRows(rows)
			},
			expectedLen: 1,
		},
		{
			name:  "поиск задач по query",
			tags:  []string{},
			query: "test",
			limit: 10, offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows(baseCols).
					AddRow(uuid.New(), boardID, "Board 1", "B1", "Todo", columnID, int64(1), "Test Task", now, now, nil).
					AddRow(uuid.New(), boardID, "Board 1", "B1", "Todo", columnID, int64(2), "Another Test", now, now, nil)

				mock.ExpectQuery(
					baseFromJoin +
						`WHERE .+\"title\" ILIKE '%test%'.+` +
						`\"tasks\"\.\"deleted_at\" IS NULL.+\"boards\"\.\"deleted_at\" IS NULL.+` +
						`ORDER BY "tasks"\."created_at" DESC LIMIT 10`,
				).WillReturnRows(rows)
			},
			expectedLen: 2,
		},
		{
			name:  "поиск задач по тегам и query",
			tags:  []string{"feature"},
			query: "auth",
			limit: 10, offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows(baseCols).
					AddRow(uuid.New(), boardID, "Board 1", "B1", "Todo", columnID, int64(1), "Auth Feature", now, now, nil)

				mock.ExpectQuery(
					baseFromJoin +
						`WHERE .+tags @>.+` +
						`.+\"title\" ILIKE '%auth%'.+` +
						`.+\"tasks\"\.\"deleted_at\" IS NULL.+\"boards\"\.\"deleted_at\" IS NULL.+` +
						`ORDER BY "tasks"\."created_at" DESC LIMIT 10`,
				).WillReturnRows(rows)
			},
			expectedLen: 1,
		},
		{
			name:  "поиск с пагинацией",
			tags:  []string{},
			query: "",
			limit: 5, offset: 10,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows(baseCols).
					AddRow(uuid.New(), boardID, "Board 1", "B1", "Todo", columnID, int64(11), "Task 11", now, now, nil)

				mock.ExpectQuery(
					baseFromJoin +
						baseSoftDeleteFilters +
						`ORDER BY "tasks"\."created_at" DESC LIMIT 5 OFFSET 10`,
				).WillReturnRows(rows)
			},
			expectedLen: 1,
		},
		{
			name:  "пустой результат поиска",
			tags:  []string{"nonexistent"},
			query: "",
			limit: 10, offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows(baseCols)

				mock.ExpectQuery(
					baseFromJoin +
						`WHERE .+tags @>.+` +
						`.+\"tasks\"\.\"deleted_at\" IS NULL.+\"boards\"\.\"deleted_at\" IS NULL.+` +
						`ORDER BY "tasks"\."created_at" DESC LIMIT 10`,
				).WillReturnRows(rows)
			},
			expectedLen: 0,
		},
		{
			name:  "ошибка БД при поиске",
			tags:  []string{},
			query: "",
			limit: 10, offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					baseFromJoin +
						baseSoftDeleteFilters +
						`ORDER BY "tasks"\."created_at" DESC LIMIT 10`,
				).WillReturnError(errors.New("database error"))
			},
			expectedLen: 0,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock)

			repo := &Repository{pool: mock}
			tasks, err := repo.SearchTasks(context.Background(), tt.tags, tt.query, tt.limit, tt.offset)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Len(t, tasks, tt.expectedLen)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
