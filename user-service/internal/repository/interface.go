package repository

import (
	"context"
	"go-ecommerce/user-service/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}

type ILoginMethodRepository interface {
	Create(ctx context.Context, body *model.LoginMethod) (*model.LoginMethod, error)
	FindOne(ctx context.Context, filter bson.M) (*model.LoginMethod, error)
}


