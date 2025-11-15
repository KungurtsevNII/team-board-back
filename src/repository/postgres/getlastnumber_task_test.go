package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastNumberTask(t *testing.T) {
	tests := []struct {
		name        string
		boardID     uuid.UUID
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expectedNum int64
		expectedErr error
	}{
		{
			name:    "успешное получение последнего number",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"number"}).AddRow(int64(10))
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedNum: 10,
			expectedErr: nil,
		},
		{
			name:    "нет задач в доске - ErrNoRows",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnError(pgx.ErrNoRows)
			},
			expectedNum: 0,
			expectedErr: pgx.ErrNoRows,
		},
		{
			name:    "number равен нулю",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"number"}).AddRow(int64(0))
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedNum: 0,
			expectedErr: nil,
		},
		{
			name:    "большое значение number",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"number"}).AddRow(int64(999999))
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedNum: 999999,
			expectedErr: nil,
		},
		{
			name:    "отрицательное значение number",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"number"}).AddRow(int64(-1))
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedNum: -1,
			expectedErr: nil,
		},
		{
			name:    "общая ошибка БД",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnError(errors.New("database error"))
			},
			expectedNum: 0,
			expectedErr: errors.New("database error"),
		},
		{
			name:    "ошибка сканирования - неверный тип данных",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"number"}).AddRow("not_a_number")
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedNum: 0,
			expectedErr: errors.New("'int64' not supported for value kind 'string'"),
		},
		{
			name:    "несколько задач в доске - возвращается максимальный",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"number"}).AddRow(int64(15))
				mock.ExpectQuery(`SELECT "number" FROM "tasks"`).
					WillReturnRows(rows)
			},
			expectedNum: 15,
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

			number, err := repo.GetLastNumberTask(context.Background(), tt.boardID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				assert.Equal(t, int64(0), number)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedNum, number)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
