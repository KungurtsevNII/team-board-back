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
			name:   "поиск задач без фильтров",
			tags:   []string{},
			query:  "",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(1), "Task 1", now, now, nil).
					AddRow(uuid.New(), boardID, columnID, int64(2), "Task 2", now, now, nil)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name:   "поиск задач по тегам",
			tags:   []string{"urgent", "bug"},
			query:  "",
			limit:  5,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(1), "Bug Task", now, now, nil)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" WHERE tags @> .+ ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name:   "поиск задач по query",
			tags:   []string{},
			query:  "test",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(1), "Test Task", now, now, nil).
					AddRow(uuid.New(), boardID, columnID, int64(2), "Another Test", now, now, nil)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" WHERE \("title" ILIKE '%test%'\) ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name:   "поиск задач по тегам и query",
			tags:   []string{"feature"},
			query:  "auth",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(1), "Auth Feature", now, now, nil)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" WHERE \(tags @> .+ AND \("title" ILIKE '%auth%'\)\) ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name:   "поиск с пагинацией",
			tags:   []string{},
			query:  "",
			limit:  5,
			offset: 10,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(11), "Task 11", now, now, nil)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" ORDER BY "created_at" DESC LIMIT 5 OFFSET 10`).
					WillReturnRows(rows)
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name:   "поиск задач с deleted_at",
			tags:   []string{},
			query:  "",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				deletedTime := now.Add(-24 * time.Hour)
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(1), "Deleted Task", now, now, &deletedTime)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name:   "пустой результат поиска",
			tags:   []string{"nonexistent"},
			query:  "",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"})

				mock.ExpectQuery(`SELECT .+ FROM "tasks" WHERE tags @> .+ ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name:   "ошибка БД при поиске",
			tags:   []string{},
			query:  "",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT .+ FROM "tasks" ORDER BY "created_at" DESC LIMIT`).
					WillReturnError(errors.New("database error"))
			},
			expectedLen: 0,
			expectedErr: errors.New("database error"),
		},
		{
			name:   "поиск с одним тегом",
			tags:   []string{"backend"},
			query:  "",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(1), "Backend Task", now, now, nil)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" WHERE tags @> .+ ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 1,
			expectedErr: nil,
		},
		{
			name:   "поиск с большим лимитом",
			tags:   []string{},
			query:  "",
			limit:  100,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"})
				for i := 1; i <= 50; i++ {
					rows.AddRow(uuid.New(), boardID, columnID, int64(i), "Task", now, now, nil)
				}

				mock.ExpectQuery(`SELECT .+ FROM "tasks" ORDER BY "created_at" DESC LIMIT 100`).
					WillReturnRows(rows)
			},
			expectedLen: 50,
			expectedErr: nil,
		},
		{
			name:   "поиск с регистронезависимым query",
			tags:   []string{},
			query:  "TEST",
			limit:  10,
			offset: 0,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "board_id", "column_id", "number", "title", "created_at", "updated_at", "deleted_at"}).
					AddRow(uuid.New(), boardID, columnID, int64(1), "test task", now, now, nil)

				mock.ExpectQuery(`SELECT .+ FROM "tasks" WHERE \("title" ILIKE '%TEST%'\) ORDER BY "created_at" DESC LIMIT`).
					WillReturnRows(rows)
			},
			expectedLen: 1,
			expectedErr: nil,
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
