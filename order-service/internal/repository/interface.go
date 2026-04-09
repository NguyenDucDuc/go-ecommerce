package repository

import (
	"context"
	"go-ecommerce/order-service/internal/model"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
}