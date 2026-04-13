package repository

import (
	"context"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/internal/model"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LoginMethodRepository struct {
	collection *mongo.Collection
}

func NewLoginMethodRepository(db *mongo.Database) ILoginMethodRepository {
	return &LoginMethodRepository{
		collection: db.Collection("login_methods"),
	}
}

func (loginMethodRepo *LoginMethodRepository) Create(ctx context.Context, body *model.LoginMethod) (*model.LoginMethod, error) {
	res, err := loginMethodRepo.collection.InsertOne(ctx, body)
	if err != nil {
		return &model.LoginMethod{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed to insert login method to database")
	}

	body.ID = res.InsertedID.(bson.ObjectID)
	return body, nil
}

func (loginMethodRepo *LoginMethodRepository) FindOne(ctx context.Context, filter bson.M) (*model.LoginMethod, error) {
	var loginMethod model.LoginMethod

	err := loginMethodRepo.collection.FindOne(ctx, filter).Decode(&loginMethod)
	if err != nil {
        if err == mongo.ErrNoDocuments {
            return &model.LoginMethod{}, nil // Không tìm thấy bản ghi nào
        }
        return &model.LoginMethod{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed to find one login method")
    }

	return &loginMethod, nil
}

func (loginMethodRepo *LoginMethodRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
    res, err := loginMethodRepo.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
    if err != nil {
        return util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed to update login method")
    }
    if res.MatchedCount == 0 {
        return nil 
    }
    return nil
}