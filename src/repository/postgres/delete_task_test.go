package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteTask(t *testing.T) {
	fixedTime := time.Date(2025, 11, 21, 15, 1, 57, 38200000, time.UTC)

	tests := []struct {
		name        string
		taskID      uuid.UUID
		mockSetup   func(mock pgxmock.PgxPoolIface, taskID uuid.UUID)
		expectedErr error
	}{
		{
			name:   "успешное мягкое удаление задачи",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name:   "задача не найдена - 0 строк обновлено",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedErr: nil,
		},
		{
			name:   "несколько задач обновлено",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 3))
			},
			expectedErr: nil,
		},
		{
			name:   "ошибка БД при обновлении",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnError(errors.New("database connection error"))
			},
			expectedErr: errors.New("database connection error"),
		},
		{
			name:   "ошибка нарушения внешнего ключа",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnError(&pgconn.PgError{Code: "23503", Message: "foreign key violation"})
			},
			expectedErr: &pgconn.PgError{Code: "23503", Message: "foreign key violation"},
		},
		{
			name:   "ошибка нарушения ограничения NOT NULL",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnError(&pgconn.PgError{Code: "23502", Message: "not null violation"})
			},
			expectedErr: &pgconn.PgError{Code: "23502", Message: "not null violation"},
		},
		{
			name:   "задача уже удалена - повторное обновление",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name:   "ошибка таймаута контекста",
			taskID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedErr: context.DeadlineExceeded,
		},
		{
			name:   "пустой UUID задачи",
			taskID: uuid.Nil,
			mockSetup: func(mock pgxmock.PgxPoolIface, taskID uuid.UUID) {
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.taskID)

			repo := &Repository{pool: mock}
			
			originalTimeNow := timeNow
			timeNow = func() time.Time { return fixedTime }
			defer func() { timeNow = originalTimeNow }()

			err = repo.DeleteTask(context.Background(), tt.taskID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				if pgErr, ok := tt.expectedErr.(*pgconn.PgError); ok {
					if actualPgErr, ok := err.(*pgconn.PgError); ok {
						assert.Equal(t, pgErr.Code, actualPgErr.Code)
					} else {
						assert.ErrorContains(t, err, pgErr.Code)
					}
				} else {
					assert.ErrorContains(t, err, tt.expectedErr.Error())
				}
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

var timeNow = time.Now