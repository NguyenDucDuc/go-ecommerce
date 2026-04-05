package repository

import (
	"context"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/internal/model"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) IProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

func (productRepo *ProductRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	result, err := productRepo.collection.InsertOne(ctx, product)
	if err != nil {
		return &model.Product{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed insert to database")
	}

	product.ID = result.InsertedID.(bson.ObjectID)

	return product, nil
}