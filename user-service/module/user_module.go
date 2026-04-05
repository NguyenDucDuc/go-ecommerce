package module

import (
	"go-ecommerce/user-service/internal/repository"
	"go-ecommerce/user-service/internal/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserModule struct {
	Repo repository.IUserRepository
	Service *service.UserService
}

func NewUserModule(db *mongo.Database) *UserModule{
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	return &UserModule{
		Repo: userRepo,
		Service: userService,
	}
}