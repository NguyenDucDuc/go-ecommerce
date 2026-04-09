package repository

import (
	"context"
	"go-ecommerce/order-service/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	UpdateOne(ctx context.Context, id bson.ObjectID, updateData bson.M) error
}