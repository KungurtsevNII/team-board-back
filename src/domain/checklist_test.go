package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecklist(t *testing.T) {
	items := []ChecklistItem{
		{Title: "Item 1", Completed: false},
		{Title: "Item 2", Completed: true},
	}

	testCases := []struct {
		name  string
		title string
		items []ChecklistItem
	}{
		{
			name:  "creates a new checklist with items",
			title: "My Checklist",
			items: items,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checklist := NewChecklist(tc.title, tc.items)

			assert.Equal(t, tc.title, checklist.Title)
			assert.Equal(t, tc.items, checklist.Items)
			assert.Len(t, checklist.Items, 2)
		})
	}
}

func TestNewChecklistItem(t *testing.T) {
	testCases := []struct {
		name      string
		title     string
		completed bool
	}{
		{
			name:      "creates a new uncompleted checklist item",
			title:     "First Item",
			completed: false,
		},
		{
			name:      "creates a new completed checklist item",
			title:     "Second Item",
			completed: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item := NewChecklistItem(tc.title, tc.completed)

			assert.Equal(t, tc.title, item.Title)
			assert.Equal(t, tc.completed, item.Completed)
		})
	}
}
