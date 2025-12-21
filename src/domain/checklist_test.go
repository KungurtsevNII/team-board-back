package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecklist(t *testing.T) {
	t.Run("creates a new checklist with items", func(t *testing.T) {
		items := []ChecklistItem{
			{Title: "Item 1", Completed: false},
			{Title: "Item 2", Completed: true},
		}
		title := "My Checklist"

		checklist := NewChecklist(title, items)

		assert.Equal(t, title, checklist.Title)
		assert.Equal(t, items, checklist.Items)
		assert.Len(t, checklist.Items, 2)
	})
}

func TestNewChecklistItem(t *testing.T) {
	t.Run("creates a new checklist item", func(t *testing.T) {
		title := "First Item"
		completed := false

		item := NewChecklistItem(title, completed)

		assert.Equal(t, title, item.Title)
		assert.Equal(t, completed, item.Completed)
	})

	t.Run("creates a new completed checklist item", func(t *testing.T) {
		title := "Second Item"
		completed := true

		item := NewChecklistItem(title, completed)

		assert.Equal(t, title, item.Title)
		assert.Equal(t, completed, item.Completed)
	})
}
