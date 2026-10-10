package repositories

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"project/internal/model"
)

type ItemDAO struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	ImageURL       string   `json:"image_url"`
	Description    string   `json:"description"`
	Price          float32  `json:"price"`
	Rating         float32  `json:"rating"`
	Specifications []string `json:"specifications"`
}

type ItemRepository struct {
	itemsMap map[string]ItemDAO
}

func NewRepository() ItemRepository {
	itemsMap := loadItems()
	return ItemRepository{
		itemsMap: itemsMap,
	}
}

func (r *ItemRepository) GetItemsByIDs(ctx context.Context, ids []string) []model.ItemDTO {
	// Receives an array of ids and returns an array containing the items that match those ids

	var selected []model.ItemDTO
	for i, _ := range ids {
		value, found := r.itemsMap[ids[i]]
		if found {
			selected = append(selected, model.ItemDTO{
				ID:             value.ID,
				Name:           value.Name,
				ImageURL:       value.ImageURL,
				Description:    value.Description,
				Price:          value.Price,
				Rating:         value.Rating,
				Specifications: value.Specifications,
			})
		}
	}
	if selected == nil {
		return []model.ItemDTO{} // ensure non-nil
	}

	return selected
}

func loadItems() map[string]ItemDAO {
	paths := []string{
		"mocked_items.json",
		filepath.Join("..", "mocked_items.json"),
	}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		var items []ItemDAO
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&items); err != nil {
			return map[string]ItemDAO{}
		}

		itemsMap := make(map[string]ItemDAO, len(items))
		for _, item := range items {
			itemsMap[item.ID] = item
		}
		return itemsMap
	}

	return map[string]ItemDAO{}
}
