package module

import (
	"go-ecommerce/common/pkg/rabbitmq"
	"go-ecommerce/user-service/internal/repository"
	"go-ecommerce/user-service/internal/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserModule struct {
	Repo repository.IUserRepository
	Service *service.UserService
}

func NewUserModule(db *mongo.Database,loginMethodService *service.LoginMethodService,rabbitService rabbitmq.IRabbitMQService) *UserModule{
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, loginMethodService ,rabbitService)
	return &UserModule{
		Repo: userRepo,
		Service: userService,
	}
}