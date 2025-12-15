package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateBoard(t *testing.T) {
	now := time.Now().UTC()
	boardID := uuid.New()

	tests := []struct {
		name        string
		board       *domain.Board
		mockSetup   func(mock pgxmock.PgxPoolIface, board *domain.Board)
		expectedErr error
	}{
		{
			name: "успешное обновление доски без удаления",
			board: &domain.Board{
				ID:        boardID,
				Name:      "Updated Board",
				ShortName: "UB",
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, board *domain.Board) {
				mock.ExpectBegin()

				// goqu инлайнит значения. Используем regex для проверки
				// Ожидаем: UPDATE "boards" SET "deleted_at"=NULL,"name"='Updated Board',"short_name"='UB',"updated_at"='...' WHERE (("id" = 'uuid') AND ("deleted_at" IS NULL))
				mock.ExpectExec(`UPDATE "boards" SET .+"name"='Updated Board',"short_name"='UB'.+ WHERE`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "успешное мягкое удаление доски с каскадным удалением",
			board: &domain.Board{
				ID:        boardID,
				Name:      "Deleted Board",
				ShortName: "DB",
				UpdatedAt: now,
				DeletedAt: &now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, board *domain.Board) {
				mock.ExpectBegin()

				// Обновление самой доски
				mock.ExpectExec(`UPDATE "boards" SET "deleted_at"=.+,"name"='Deleted Board'.+ WHERE`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				// Каскадное обновление задач
				mock.ExpectExec(`UPDATE "tasks" SET "deleted_at"=.+ WHERE`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 5))

				// Каскадное обновление колонок
				mock.ExpectExec(`UPDATE "columns" SET "deleted_at"=.+ WHERE`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 3))

				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "ошибка начала транзакции",
			board: &domain.Board{
				ID:        boardID,
				Name:      "Error Board",
				ShortName: "EB",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, board *domain.Board) {
				mock.ExpectBegin().WillReturnError(errors.New("tx error"))
			},
			expectedErr: errors.New("tx error"),
		},
		{
			name: "ошибка обновления доски",
			board: &domain.Board{
				ID:        boardID,
				Name:      "Update Error",
				ShortName: "UE",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, board *domain.Board) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "boards"`).
					WillReturnError(errors.New("update error"))

				mock.ExpectRollback()
			},
			expectedErr: errors.New("update error"),
		},
		{
			name: "ошибка каскадного обновления задач",
			board: &domain.Board{
				ID:        boardID,
				Name:      "Cascade Error",
				ShortName: "CE",
				DeletedAt: &now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, board *domain.Board) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "boards"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				mock.ExpectExec(`UPDATE "tasks"`).
					WillReturnError(errors.New("tasks update error"))

				mock.ExpectRollback()
			},
			expectedErr: errors.New("tasks update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.board)

			repo := &Repository{pool: mock}
			err = repo.UpdateBoard(context.Background(), tt.board)

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
