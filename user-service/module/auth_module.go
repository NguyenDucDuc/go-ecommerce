package module

import (
	"go-ecommerce/common/pkg/jwt"
	"go-ecommerce/user-service/internal/repository"
	"go-ecommerce/user-service/internal/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthModule struct {
	AuthService *service.AuthService
}

func NewAuthModule(db *mongo.Database, jwtService jwt.IJwtService) *AuthModule {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtService)
	return &AuthModule{
		AuthService: authService,
	}
}