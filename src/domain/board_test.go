package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	testCases := []struct {
		name        string
		boardName   string
		shortName   string
		expectError bool
		errorType   error
	}{
		{
			name:        "Success: creates a new board with a default column",
			boardName:   "My Project",
			shortName:   "MP-1",
			expectError: false,
		},
		{
			name:        "Failure: empty board name",
			boardName:   "",
			shortName:   "VALID",
			expectError: true,
			errorType:   ErrInvalidName,
		},
		{
			name:        "Failure: board name too long",
			boardName:   string(make([]byte, 101)),
			shortName:   "VALID",
			expectError: true,
			errorType:   ErrInvalidName,
		},
		{
			name:        "Failure: empty short name",
			boardName:   "Valid Name",
			shortName:   "",
			expectError: true,
			errorType:   ErrInvalidName,
		},
		{
			name:        "Failure: short name too short",
			boardName:   "Valid Name",
			shortName:   "A",
			expectError: true,
			errorType:   ErrInvalidName,
		},
		{
			name:        "Failure: short name too long",
			boardName:   "Valid Name",
			shortName:   "ABCDEFGHIJK",
			expectError: true,
			errorType:   ErrInvalidName,
		},
		{
			name:        "Failure: short name with invalid characters",
			boardName:   "Valid Name",
			shortName:   "INVALID!",
			expectError: true,
			errorType:   ErrInvalidName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, err := NewBoard(tc.boardName, tc.shortName)

			if tc.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.errorType)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, board.ID)
				assert.Equal(t, tc.boardName, board.Name)
				assert.Equal(t, tc.shortName, board.ShortName)
				assert.WithinDuration(t, time.Now().UTC(), board.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now().UTC(), board.UpdatedAt, time.Second)
				assert.Nil(t, board.DeletedAt)

				require.Len(t, board.Columns, 1, "Board should have one initial column")
				assert.Equal(t, nameOfFirstColumn, board.Columns[0].Name)
				assert.Equal(t, board.ID, board.Columns[0].BoardID)
			}
		})
	}
}

func TestBoard_GetFirstColumn(t *testing.T) {
	boardWithColumn, err := NewBoard("Test Board", "TB")
	require.NoError(t, err)
	boardWithoutColumn := Board{Columns: []Column{}}

	testCases := []struct {
		name        string
		board       Board
		expectError bool
		errorType   error
	}{
		{
			name:        "Success: returns the first column",
			board:       boardWithColumn,
			expectError: false,
		},
		{
			name:        "Failure: returns error when no columns exist",
			board:       boardWithoutColumn,
			expectError: true,
			errorType:   ErrColumnsIsEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			firstCol, err := tc.board.GetFirstColumn()

			if tc.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.errorType)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.board.Columns[0].ID, firstCol.ID)
				assert.Equal(t, nameOfFirstColumn, firstCol.Name)
			}
		})
	}
}

func TestBoard_Delete(t *testing.T) {
	board, err := NewBoard("Test Board", "TB")
	require.NoError(t, err)

	testCases := []struct {
		name          string
		boardToDelete *Board
	}{
		{
			name:          "marks a board as deleted",
			boardToDelete: &board,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Nil(t, tc.boardToDelete.DeletedAt)

			tc.boardToDelete.Delete()

			require.NotNil(t, tc.boardToDelete.DeletedAt)
			assert.WithinDuration(t, time.Now().UTC(), *tc.boardToDelete.DeletedAt, time.Second)
		})
	}
}
