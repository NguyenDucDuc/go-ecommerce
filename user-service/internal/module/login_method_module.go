package module

import (
	"go-ecommerce/user-service/internal/repository"
	"go-ecommerce/user-service/internal/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LoginMethodModule struct {
	Repo repository.ILoginMethodRepository
	Service *service.LoginMethodService
}

func NewLoginMethodModule(db *mongo.Database) *LoginMethodModule {
	repo := repository.NewLoginMethodRepository(db)
	service := service.NewLoginMethodService(repo)
	return &LoginMethodModule{
		Repo: repo,
		Service: service,
	}
}