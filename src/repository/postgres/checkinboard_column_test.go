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

func TestCheckColumnInBoard(t *testing.T) {
	tests := []struct {
		name         string
		boardID      uuid.UUID
		columnID     uuid.UUID
		mockSetup    func(mock pgxmock.PgxPoolIface)
		expectedBool bool
		expectedErr  error
	}{
		{
			name:     "колонка существует в доске",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(1))
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			expectedBool: true,
			expectedErr:  nil,
		},
		{
			name:     "колонка не существует в доске",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(0))
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			expectedBool: false,
			expectedErr:  nil,
		},
		{
			name:     "несколько совпадений - возвращается true",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(5))
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			expectedBool: true,
			expectedErr:  nil,
		},
		{
			name:     "общая ошибка БД",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnError(errors.New("database error"))
			},
			expectedBool: false,
			expectedErr:  errors.New("database error"),
		},
		{
			name:     "ошибка сканирования - неверный тип данных",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow("not_a_number")
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			expectedBool: false,
			expectedErr:  errors.New("'int64' not supported for value kind 'string'"),
		},
		{
			name:     "пустой результат запроса",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"})
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			expectedBool: false,
			expectedErr:  errors.New("no rows"),
		},
		{
			name:     "отрицательное значение count (невозможно в реальности)",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(-1))
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			expectedBool: false,
			expectedErr:  nil,
		},
		{
			name:     "большое значение count",
			boardID:  uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).AddRow(int64(100000))
				mock.ExpectQuery(`SELECT COUNT`).
					WillReturnRows(rows)
			},
			expectedBool: true,
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock)
			repo := &Repository{pool: mock}

			exists, err := repo.CheckColumnInBoard(context.Background(), tt.boardID, tt.columnID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				assert.False(t, exists)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBool, exists)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
