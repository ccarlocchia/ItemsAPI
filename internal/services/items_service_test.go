package services

import (
	"context"
	"project/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepository is a test double for ItemRepository
type mockRepository struct {
	expectedIDs []string
	returnItems []model.ItemDTO
	calledWith  []string
}

func (m *mockRepository) GetItemsByIDs(ctx context.Context, ids []string) []model.ItemDTO {
	m.calledWith = ids
	return m.returnItems
}

func TestItemsService_GetItemsByIDs(t *testing.T) {
	tests := []struct {
		name        string
		inputIDs    []string
		mockReturn  []model.ItemDTO
		expectedIDs []string
	}{
		{
			name:     "returns items from repository",
			inputIDs: []string{"1", "2"},
			mockReturn: []model.ItemDTO{
				{ID: "1", Name: "Item One"},
				{ID: "2", Name: "Item Two"},
			},
			expectedIDs: []string{"1", "2"},
		},
		{
			name:        "handles empty ID list",
			inputIDs:    []string{},
			mockReturn:  []model.ItemDTO{},
			expectedIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{returnItems: tt.mockReturn}
			service := NewItemsService(mockRepo)

			ctx := context.Background()
			got := service.GetItemsByIDs(ctx, tt.inputIDs)

			// Check that the repository was called with the correct IDs
			assert.Equal(t, tt.inputIDs, mockRepo.calledWith, "repository should be called with correct IDs")

			// Check that the service returns the repository result
			assert.Equal(t, tt.mockReturn, got, "service should return the expected repository result")
		})
	}
}
