package createcolumn

import (
	"context"
	"errors"
	"testing"

	"github.com/KungurtsevNII/team-board-back/src/usecase/createcolumn/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	boardID := uuid.New()
	ctx := context.Background()

	testCases := []struct {
		name          string
		command       Command
		setupMock     func(*mocks.Repo)
		expectedOrder int64
		expectError   error
	}{
		{
			name: "Success: create first column",
			command: Command{
				BoardID: boardID,
				Name:    "First Column",
			},
			setupMock: func(repo *mocks.Repo) {
				repo.On("CheckBoard", mock.Anything, boardID.String()).Return(true).Once()
				repo.On("GetLastOrderNumColumn", mock.Anything, boardID).Return(int64(0), pgx.ErrNoRows).Once()
				repo.On("CreateColumn", mock.Anything, mock.AnythingOfType("*domain.Column")).Return(nil).Once()
			},
			expectedOrder: 0,
			expectError:   nil,
		},
		{
			name: "Success: create subsequent column",
			command: Command{
				BoardID: boardID,
				Name:    "Next Column",
			},
			setupMock: func(repo *mocks.Repo) {
				repo.On("CheckBoard", mock.Anything, boardID.String()).Return(true).Once()
				repo.On("GetLastOrderNumColumn", mock.Anything, boardID).Return(int64(2), nil).Once()
				repo.On("CreateColumn", mock.Anything, mock.AnythingOfType("*domain.Column")).Return(nil).Once()
			},
			expectedOrder: 3,
			expectError:   nil,
		},
		{
			name: "Failure: board not found",
			command: Command{
				BoardID: boardID,
				Name:    "Some Column",
			},
			setupMock: func(repo *mocks.Repo) {
				repo.On("CheckBoard", mock.Anything, boardID.String()).Return(false).Once()
			},
			expectError: ErrBoardIsNotExists,
		},
		{
			name: "Failure: get last order number fails",
			command: Command{
				BoardID: boardID,
				Name:    "Some Column",
			},
			setupMock: func(repo *mocks.Repo) {
				repo.On("CheckBoard", mock.Anything, boardID.String()).Return(true).Once()
				repo.On("GetLastOrderNumColumn", mock.Anything, boardID).Return(int64(0), errors.New("db connection error")).Once()
			},
			expectError: ErrGetLastOrderNumUnknown,
		},
		{
			name: "Failure: create column in repo fails",
			command: Command{
				BoardID: boardID,
				Name:    "Some Column",
			},
			setupMock: func(repo *mocks.Repo) {
				repo.On("CheckBoard", mock.Anything, boardID.String()).Return(true).Once()
				repo.On("GetLastOrderNumColumn", mock.Anything, boardID).Return(int64(0), nil).Once()
				repo.On("CreateColumn", mock.Anything, mock.AnythingOfType("*domain.Column")).Return(errors.New("db unique constraint violated")).Once()
			},
			expectError: ErrCreateColumnUnknown,
		},
		{
			name: "Failure: validation failed (empty name)",
			command: Command{
				BoardID: boardID,
				Name:    "",
			},
			setupMock: func(repo *mocks.Repo) {
				repo.On("CheckBoard", mock.Anything, boardID.String()).Return(true).Once()
				repo.On("GetLastOrderNumColumn", mock.Anything, boardID).Return(int64(0), nil).Once()
			},
			expectError: ErrValidationFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepo(t)
			if tc.setupMock != nil {
				tc.setupMock(repo)
			}

			uc := NewUC(repo)

			column, err := uc.Handle(ctx, tc.command)

			if tc.expectError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectError, "Wrong error type")
				assert.Nil(t, column)
			} else {
				require.NoError(t, err)
				require.NotNil(t, column)
				assert.Equal(t, tc.command.Name, column.Name)
				assert.Equal(t, tc.command.BoardID, column.BoardID)
				assert.Equal(t, tc.expectedOrder, column.OrderNum)
			}
			
			repo.AssertExpectations(t)
		})
	}
}
