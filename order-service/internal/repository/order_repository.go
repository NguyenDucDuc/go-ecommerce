package repository

import (
	"context"
	util "go-ecommerce/common/utils"
	"go-ecommerce/order-service/internal/model"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) IOrderRepository{
	return &OrderRepository{
		collection: db.Collection("orders"),
	}
}

func (orderRepo *OrderRepository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	res, err := orderRepo.collection.InsertOne(ctx, order)
	if err != nil {
		return &model.Order{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed to insert order into database")
	}

	order.ID = res.InsertedID.(bson.ObjectID)

	return order, nil
}

func (orderRepo *OrderRepository) UpdateOne(ctx context.Context, id bson.ObjectID, updateData bson.M) error {
    filter := bson.M{"_id": id}

    // Đảm bảo luôn có updated_at
    updateData["updated_at"] = time.Now()

    // Bọc tất cả vào trong $set
    finalUpdate := bson.M{
        "$set": updateData,
    }

    _, err := orderRepo.collection.UpdateOne(ctx, filter, finalUpdate)
    return err
}