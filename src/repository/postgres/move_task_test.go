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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoveTaskColumn(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name         string
		taskID       uuid.UUID
		columnID     uuid.UUID
		mockSetup    func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID)
		expectedTask *domain.Task
		expectedErr  error
	}{
		{
			name:     "успешное перемещение задачи",
			taskID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			columnID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID) {
				checklistsJSON, _ := json.Marshal([]domain.Checklist{
					{Title: "Test Checklist"},
				})

				rows := pgxmock.NewRows([]string{
					"id", "board_id", "column_id", "number", "title",
					"description", "tags", "checklists", "created_at",
					"updated_at", "deleted_at",
				}).AddRow(
					taskID,
					uuid.MustParse("33333333-3333-3333-3333-333333333333"),
					columnID,
					int64(5),
					"Test Task",
					stringPtr("Description"),
					[]string{"tag1"},
					checklistsJSON,
					now,
					now,
					nil,
				)

				// Убираем WithArgs - принимаем любые аргументы
				mock.ExpectQuery(`UPDATE "tasks"`).
					WillReturnRows(rows)
			},
			expectedTask: &domain.Task{
				ID:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				BoardID:  uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				ColumnID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				Number:   5,
				Title:    "Test Task",
			},
			expectedErr: nil,
		},
		{
			name:     "задача не найдена - ErrNoRows",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID) {
				// Убираем WithArgs
				mock.ExpectQuery(`UPDATE "tasks"`).
					WillReturnError(pgx.ErrNoRows)
			},
			expectedTask: nil,
			expectedErr:  pgx.ErrNoRows,
		},
		{
			name:     "задача была удалена (deleted_at IS NOT NULL)",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID) {
				// Убираем WithArgs
				mock.ExpectQuery(`UPDATE "tasks"`).
					WillReturnError(pgx.ErrNoRows)
			},
			expectedTask: nil,
			expectedErr:  pgx.ErrNoRows,
		},
		{
			name:     "нарушение внешнего ключа - колонка не существует",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID) {
				// Убираем WithArgs
				mock.ExpectQuery(`UPDATE "tasks"`).
					WillReturnError(&pgconn.PgError{
						Code:    "23503",
						Message: "foreign key violation",
					})
			},
			expectedTask: nil,
			expectedErr:  &pgconn.PgError{Code: "23503"},
		},
		{
			name:     "общая ошибка БД",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID) {
				// Убираем WithArgs
				mock.ExpectQuery(`UPDATE "tasks"`).
					WillReturnError(errors.New("database connection lost"))
			},
			expectedTask: nil,
			expectedErr:  errors.New("database connection lost"),
		},
		{
			name:     "timeout при обновлении",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID) {
				// Убираем WithArgs
				mock.ExpectQuery(`UPDATE "tasks"`).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedTask: nil,
			expectedErr:  context.DeadlineExceeded,
		},
		{
			name:     "некорректный JSON в checklists после UPDATE",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID, columnID uuid.UUID) {
				badJSON := []byte(`{"invalid": [}`)

				rows := pgxmock.NewRows([]string{
					"id", "board_id", "column_id", "number", "title",
					"description", "tags", "checklists", "created_at",
					"updated_at", "deleted_at",
				}).AddRow(
					taskID,
					uuid.New(),
					columnID,
					int64(1),
					"Task",
					nil,
					[]string{},
					badJSON,
					now,
					now,
					nil,
				)

				// Убираем WithArgs
				mock.ExpectQuery(`UPDATE "tasks"`).
					WillReturnRows(rows)
			},
			expectedTask: nil,
			expectedErr:  errors.New("invalid character"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.taskID, tt.columnID)

			repo := &Repository{pool: mock}
			task, err := repo.MoveTaskColumn(context.Background(), tt.taskID, tt.columnID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				if pgErr, ok := tt.expectedErr.(*pgconn.PgError); ok {
					assert.Contains(t, err.Error(), pgErr.Code)
				} else {
					assert.ErrorContains(t, err, tt.expectedErr.Error())
				}
				assert.Nil(t, task)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)
				assert.Equal(t, tt.expectedTask.ID, task.ID)
				assert.Equal(t, tt.expectedTask.BoardID, task.BoardID)
				assert.Equal(t, tt.expectedTask.ColumnID, task.ColumnID)
				assert.Equal(t, tt.expectedTask.Number, task.Number)
				assert.Equal(t, tt.expectedTask.Title, task.Title)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
