package createcolumn

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	validBoardID := uuid.New().String()

	testCases := []struct {
		name        string
		boardID     string
		columnName  string
		expectError error
	}{
		{
			name:        "Success: valid inputs",
			boardID:     validBoardID,
			columnName:  "My Column",
			expectError: nil,
		},
		{
			name:        "Failure: invalid board ID format",
			boardID:     "not-a-valid-uuid",
			columnName:  "My Column",
			expectError: ErrInvalidUUID,
		},
		{
			name:        "Failure: empty column name",
			boardID:     validBoardID,
			columnName:  "", 
			expectError: ErrValidationFailed,
		},
		{
			name:        "Failure: column name too long",
			boardID:     validBoardID,
			columnName:  strings.Repeat("a", 101),
			expectError: ErrValidationFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := NewCommand(tc.boardID, tc.columnName)

			if tc.expectError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectError, "Wrong error type")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.boardID, cmd.BoardID.String())
				assert.Equal(t, tc.columnName, cmd.Name)
			}
		})
	}
}
