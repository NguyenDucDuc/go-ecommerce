package repository

import (
	"context"
	"go-ecommerce/user-service/internal/model"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}


