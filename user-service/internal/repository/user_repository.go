package repository

import (
	"context"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/internal/model"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) IUserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (userRepo *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
	result, err := userRepo.collection.InsertOne(ctx, user)
	if err != nil {
		return &model.User{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Insert to database failed")
	}

	user.ID = result.InsertedID.(bson.ObjectID)
	return user, nil
}

func (userRepo *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	filter := bson.D{{Key: "email", Value: email}}
	err := userRepo.collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return &model.User{}, err
	}

	return &user, nil
}