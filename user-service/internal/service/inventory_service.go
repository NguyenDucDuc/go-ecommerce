package service

import (
	"context"
	"go-ecommerce/user-service/internal/model"
	"go-ecommerce/user-service/internal/repository"
)

type InventoryService struct {
	repo repository.IInventoryRepository
}

func NewInventoryService(repo repository.IInventoryRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (inventoryService *InventoryService) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	res, err := inventoryService.repo.Create(ctx, inventory)
	if err != nil {
		return &model.Inventory{}, err
	}

	return res, nil
}