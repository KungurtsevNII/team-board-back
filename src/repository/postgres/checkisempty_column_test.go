package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckColumnIsEmpty(t *testing.T) {
	tests := []struct {
		name        string
		columnID    uuid.UUID
		mockSetup   func(mock pgxmock.PgxPoolIface, columnID uuid.UUID)
		expectedRes bool
		expectErr   bool
	}{
		{
			name:     "колонка пустая (count_0)",
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, columnID uuid.UUID) {
				// goqu генерирует SELECT COUNT(*) FROM "tasks" WHERE ...,
				// в тестах достаточно частичного паттерна.
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(0))
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedRes: true,
			expectErr:   false,
		},
		{
			name:     "колонка не пустая (count_>0)",
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, columnID uuid.UUID) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(5))
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedRes: false,
			expectErr:   false,
		},
		{
			name:     "ошибка выполнения запроса",
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, columnID uuid.UUID) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "tasks"`).
					WillReturnError(errors.New("db error"))
			},
			expectedRes: false,
			expectErr:   true,
		},
		{
			name:     "ошибка сканирования результата",
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, columnID uuid.UUID) {
				// Некорректный тип колонки, чтобы спровоцировать scan error
				rows := pgxmock.NewRows([]string{"count"}).AddRow("not-an-int")
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedRes: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.columnID)

			repo := &Repository{pool: mock}
			res, err := repo.CheckColumnIsEmpty(context.Background(), tt.columnID)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
