package gettask

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	validTaskID := uuid.New().String()

	testCases := []struct {
		name        string
		taskID      string
		expectError error
	}{
		{
			name:        "Success: valid inputs",
			taskID:      validTaskID,
			expectError: nil,
		},
		{
			name:        "Failure: invalid task ID format",
			taskID:      "not-a-valid-uuid",
			expectError: ErrInvalidTaskID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := NewQuery(tc.taskID)

			if tc.expectError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectError, "Wrong error type")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.taskID, cmd.TaskID.String())
			}
		})
	}
}
