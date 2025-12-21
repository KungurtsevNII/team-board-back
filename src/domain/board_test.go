package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	t.Run("Success: creates a new board with a default column", func(t *testing.T) {
		name := "My Project"
		shortName := "MP-1"

		board, err := NewBoard(name, shortName)

		require.NoError(t, err)

		assert.NotEmpty(t, board.ID)
		assert.Equal(t, name, board.Name)
		assert.Equal(t, shortName, board.ShortName)
		assert.WithinDuration(t, time.Now().UTC(), board.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), board.UpdatedAt, time.Second)
		assert.Nil(t, board.DeletedAt)

		require.Len(t, board.Columns, 1, "Board should have one initial column")
		assert.Equal(t, nameOfFirstColumn, board.Columns[0].Name)
		assert.Equal(t, board.ID, board.Columns[0].BoardID)
	})

	validationTestCases := []struct {
		name        string
		boardName   string
		shortName   string
		expectError error
	}{
		{
			name:        "Failure: empty board name",
			boardName:   "",
			shortName:   "VALID",
			expectError: ErrInvalidName,
		},
		{
			name:        "Failure: board name too long",
			boardName:   string(make([]byte, 101)),
			shortName:   "VALID",
			expectError: ErrInvalidName,
		},
		{
			name:        "Failure: empty short name",
			boardName:   "Valid Name",
			shortName:   "",
			expectError: ErrInvalidName,
		},
		{
			name:        "Failure: short name too short",
			boardName:   "Valid Name",
			shortName:   "A",
			expectError: ErrInvalidName,
		},
		{
			name:        "Failure: short name too long",
			boardName:   "Valid Name",
			shortName:   "ABCDEFGHIJK",
			expectError: ErrInvalidName,
		},
		{
			name:        "Failure: short name with invalid characters",
			boardName:   "Valid Name",
			shortName:   "INVALID!",
			expectError: ErrInvalidName,
		},
	}

	for _, tc := range validationTestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewBoard(tc.boardName, tc.shortName)
			
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.expectError)
		})
	}
}

func TestBoard_GetFirstColumn(t *testing.T) {
	t.Run("Success: returns the first column", func(t *testing.T) {
		board, err := NewBoard("Test Board", "TB")
		require.NoError(t, err)

		firstCol, err := board.GetFirstColumn()

		require.NoError(t, err)
		assert.Equal(t, board.Columns[0].ID, firstCol.ID)
		assert.Equal(t, nameOfFirstColumn, firstCol.Name)
	})

	t.Run("Failure: returns error when no columns exist", func(t *testing.T) {
		board := Board{Columns: []Column{}}

		_, err := board.GetFirstColumn()

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrColumnsIsEmpty)
	})
}

func TestBoard_Delete(t *testing.T) {
	t.Run("marks a board as deleted", func(t *testing.T) {
		board, err := NewBoard("Test Board", "TB")
		require.NoError(t, err)
		require.Nil(t, board.DeletedAt)

		board.Delete()

		require.NotNil(t, board.DeletedAt)
		assert.WithinDuration(t, time.Now().UTC(), *board.DeletedAt, time.Second)
	})
}
