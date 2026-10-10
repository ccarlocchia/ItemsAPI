package repositories

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"project/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to create a temporary JSON file for loadItems()
func createMockItemsFile(t *testing.T, items []ItemDAO) string {
	t.Helper()
	tmpFile := filepath.Join(t.TempDir(), "mocked_items.json")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("failed to create mock file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(items); err != nil {
		t.Fatalf("failed to encode mock items: %v", err)
	}

	return tmpFile
}

func TestGetItemsByIDs(t *testing.T) {
	// Prepare mock data
	mockItems := []ItemDAO{
		{
			ID:    "1",
			Name:  "Item One",
			Price: 10.5,
		},
		{
			ID:    "2",
			Name:  "Item Two",
			Price: 20.0,
		},
	}

	// Create repository with mock data (bypass loadItems)
	repo := ItemRepository{
		itemsMap: map[string]ItemDAO{
			"1": mockItems[0],
			"2": mockItems[1],
		},
	}

	tests := []struct {
		name     string
		ids      []string
		expected []model.ItemDTO
	}{
		{
			name: "return matching items",
			ids:  []string{"1", "2"},
			expected: []model.ItemDTO{
				{ID: "1", Name: "Item One", Price: 10.5},
				{ID: "2", Name: "Item Two", Price: 20.0},
			},
		},
		{
			name:     "ignore non-existent IDs",
			ids:      []string{"3"},
			expected: []model.ItemDTO{},
		},
		{
			name:     "empty ID list returns empty result",
			ids:      []string{},
			expected: []model.ItemDTO{},
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.GetItemsByIDs(ctx, tt.ids)
			assert.Equal(t, tt.expected, got, "expected and actual results should match")
		})
	}
}

func TestLoadItems(t *testing.T) {
	// Create temporary mock JSON file
	mockItems := []ItemDAO{
		{ID: "1", Name: "Mock Item"},
	}
	mockPath := createMockItemsFile(t, mockItems)
	defer os.Remove(mockPath)

	// Temporarily change working directory to where the mock file lives
	oldWD, _ := os.Getwd()
	defer os.Chdir(oldWD)
	_ = os.Chdir(filepath.Dir(mockPath))

	loaded := loadItems()
	if len(loaded) != 1 {
		t.Fatalf("expected 1 item, got %d", len(loaded))
	}

	if loaded["1"].Name != "Mock Item" {
		t.Errorf("expected item name 'Mock Item', got %q", loaded["1"].Name)
	}
}
