package service

import (
	"context"
	"go-ecommerce/product-service/internal/model"
	"go-ecommerce/product-service/internal/repository"
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

func (s *InventoryService) ReserveStock(ctx context.Context, items []*model.OrderItem) error {
        for _, item := range items {
            // Thực hiện trừ stock khả dụng và đẩy vào stock tạm giữ
            err := s.repo.UpdateStock(ctx, item.ProductId.Hex(), int(item.Quantity))
            if err != nil {
                return err // Trả về lỗi để Transaction tự động Rollback các món trước đó
            }
        }
        return nil
}