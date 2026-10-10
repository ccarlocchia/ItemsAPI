package services

import (
	"context"
	"project/internal/model"
)

type ItemRepository interface {
	GetItemsByIDs(ctx context.Context, ids []string) []model.ItemDTO
}

type ItemsService struct {
	repository ItemRepository
}

func NewItemsService(repository ItemRepository) *ItemsService {
	return &ItemsService{
		repository: repository,
	}
}

func (s *ItemsService) GetItemsByIDs(ctx context.Context, ids []string) []model.ItemDTO {
	return s.repository.GetItemsByIDs(ctx, ids)
}
