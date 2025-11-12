package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ErrUniqueViolation   = &pgconn.PgError{Code: "23505"}
	ErrForeignKeyViolation = &pgconn.PgError{Code: "23503"}
	ErrNotNullViolation  = &pgconn.PgError{Code: "23502"}
)

func TestCreateColumn(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		column      *domain.Column
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expectedErr error
	}{
		{
			name: "успешное создание колонки",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "To Do",
				OrderNum:  1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "создание колонки с DeletedAt",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Done",
				OrderNum:  2,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: &now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "создание колонки с пустым именем",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "",
				OrderNum:  0,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "создание колонки с длинным именем",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Very Long Column Name That Might Test Database Limits",
				OrderNum:  10,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "отрицательный order_num",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Test",
				OrderNum:  -1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "ошибка уникальности - дубликат ID",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Duplicate",
				OrderNum:  1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnError(ErrUniqueViolation)
			},
			expectedErr: ErrUniqueViolation,
		},
		{
			name: "ошибка внешнего ключа - несуществующий board_id",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Test",
				OrderNum:  1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnError(ErrForeignKeyViolation)
			},
			expectedErr: ErrForeignKeyViolation,
		},
		{
			name: "ошибка NOT NULL constraint",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Test",
				OrderNum:  1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnError(ErrNotNullViolation)
			},
			expectedErr: ErrNotNullViolation,
		},
		{
			name: "ошибка подключения",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Test",
				OrderNum:  1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnError(errors.New("connection refused"))
			},
			expectedErr: errors.New("connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock)

			repo := &Repository{pool: mock}
			err = repo.CreateColumn(context.Background(), tt.column)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
