package repository

import (
	"context"
	"go-ecommerce/product-service/internal/model"
)

type IProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
}

type IInventoryRepository interface {
	Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
}