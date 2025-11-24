package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTaskByID(t *testing.T) {
	now := time.Now()

	baseChecklist := []domain.Checklist{
		{
			Title: "Checklist 1",
			Items: []domain.ChecklistItem{
				{Title: "item1", Completed: false},
				{Title: "item2", Completed: true},
			},
		},
	}

	tests := []struct {
		name        string
		taskID      uuid.UUID
		mockSetup   func(mock pgxmock.PgxPoolIface, taskID uuid.UUID)
		expected    *domain.Task
		expectedErr error
	}{
		{
			name:   "успешное получение задачи",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				checklistsJSON, _ := json.Marshal(baseChecklist)

				rows := pgxmock.NewRows([]string{
					"id",
					"board_id",
					"column_id",
					"number",
					"title",
					"description",
					"tags",
					"checklists",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					taskID,                 // id
					uuid.MustParse("11111111-1111-1111-1111-111111111111"), // board_id
					uuid.MustParse("22222222-2222-2222-2222-222222222222"), // column_id
					int64(5),              // number
					"Test task",           // title
					stringPtr("Desc"),     // description
					[]string{"tag1"},      // tags
					checklistsJSON,        // checklists
					now,                   // created_at
					now,                   // updated_at
					nil,                   // deleted_at
				)

				mock.ExpectQuery(`SELECT .* FROM "tasks"`).
					WillReturnRows(rows)
			},
			expected: &domain.Task{
				Number:      5,
				Title:       "Test task",
				Description: stringPtr("Desc"),
				Tags:        []string{"tag1"},
				Checklists:  baseChecklist,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   zeroTimePtr(),
			},
			expectedErr: nil,
		},
		{
			name:   "задача не найдена - ErrNoRows",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectQuery(`SELECT .* FROM "tasks"`).
					WillReturnError(pgx.ErrNoRows)
			},
			expected:    nil,
			expectedErr: pgx.ErrNoRows,
		},
		{
			name:   "ошибка БД при выборке",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectQuery(`SELECT .* FROM "tasks"`).
					WillReturnError(errors.New("database error"))
			},
			expected:    nil,
			expectedErr: errors.New("database error"),
		},
		{
			name:   "некорректный JSON в checklists - ошибка Unmarshal",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				// заведомо битый JSON
				badJSON := []byte(`{"invalid": [}`)

				rows := pgxmock.NewRows([]string{
					"id",
					"board_id",
					"column_id",
					"number",
					"title",
					"description",
					"tags",
					"checklists",
					"created_at",
					"updated_at",
					"deleted_at",
				}).AddRow(
					taskID,
					uuid.MustParse("11111111-1111-1111-1111-111111111111"),
					uuid.MustParse("22222222-2222-2222-2222-222222222222"),
					int64(1),
					"Bad checklist",
					nil,
					[]string{},
					badJSON,
					now,
					now,
					nil,
				)

				mock.ExpectQuery(`SELECT .* FROM "tasks"`).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: errors.New("invalid character"),
		},
		{
			name:   "таска есть в БД, но помечена как удалённая - ErrNoRows",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectQuery(`SELECT .* FROM "tasks"`).
					WillReturnError(pgx.ErrNoRows)
			},
			expected:    nil,
			expectedErr: pgx.ErrNoRows,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.taskID)

			repo := &Repository{pool: mock}

			task, err := repo.GetTaskByID(context.Background(), tt.taskID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				// проверяем, что обёрнутая ошибка содержит текст исходной
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				assert.Nil(t, task)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)

				// заполняем ожидаемые ID, если нужно
				if tt.expected.ID == uuid.Nil {
					tt.expected.ID = tt.taskID
				}
				if tt.expected.BoardID == uuid.Nil {
					tt.expected.BoardID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
				}
				if tt.expected.ColumnID == uuid.Nil {
					tt.expected.ColumnID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
				}

				assert.Equal(t, tt.expected.ID, task.ID)
				assert.Equal(t, tt.expected.BoardID, task.BoardID)
				assert.Equal(t, tt.expected.ColumnID, task.ColumnID)
				assert.Equal(t, tt.expected.Number, task.Number)
				assert.Equal(t, tt.expected.Title, task.Title)
				assert.Equal(t, tt.expected.Description, task.Description)
				assert.Equal(t, tt.expected.Tags, task.Tags)
				assert.Equal(t, tt.expected.Checklists, task.Checklists)
				assert.WithinDuration(t, tt.expected.CreatedAt, task.CreatedAt, time.Second)
				assert.WithinDuration(t, tt.expected.UpdatedAt, task.UpdatedAt, time.Second)
				assert.Equal(t, tt.expected.DeletedAt, task.DeletedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func zeroTimePtr() *time.Time {
    t := time.Time{}
    return &t
}