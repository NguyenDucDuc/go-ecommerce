package service

import (
	"context"
	"go-ecommerce/user-service/internal/model"
	"go-ecommerce/user-service/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LoginMethodService struct {
	loginMethodRepo repository.ILoginMethodRepository
}

func NewLoginMethodService(repo repository.ILoginMethodRepository) *LoginMethodService {
	return &LoginMethodService{
		loginMethodRepo: repo,
	}
}

func (lmtService *LoginMethodService) Create(ctx context.Context, body *model.LoginMethod) (*model.LoginMethod, error) {
	loginMethod, err := lmtService.loginMethodRepo.Create(ctx, body)
	if err != nil {
		return &model.LoginMethod{}, err
	}

	return loginMethod, nil
}

func (lmtService *LoginMethodService) FindOne(ctx context.Context, filter bson.M) (*model.LoginMethod, error) {
	loginMethod, err := lmtService.loginMethodRepo.FindOne(ctx, filter)
	if err != nil {
		return &model.LoginMethod{}, err
	}

	return loginMethod, nil
}