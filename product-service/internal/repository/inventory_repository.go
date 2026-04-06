package repository

import (
	"context"
	util "go-ecommerce/common/utils"
	"go-ecommerce/product-service/internal/model"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InventoryRepository struct {
	collection *mongo.Collection
}

func NewInventoryRepository(db *mongo.Database) IInventoryRepository {
	return &InventoryRepository{
		collection: db.Collection("inventories"),
	}
}

func (inventoryRepo *InventoryRepository) Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	res, err := inventoryRepo.collection.InsertOne(ctx, inventory)
	if err != nil {
		return &model.Inventory{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed insert into database")
	}

	inventory.ID = res.InsertedID.(bson.ObjectID)
	return inventory, nil
}