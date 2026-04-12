package module

import (
	"go-ecommerce/common/pkg/jwt"
	"go-ecommerce/user-service/internal/service"
)

type AuthModule struct {
	AuthService *service.AuthService
	LoginMethodService *service.LoginMethodService
	UserService *service.UserService
}

func NewAuthModule(loginMethodService *service.LoginMethodService, userService *service.UserService ,jwtService jwt.IJwtService) *AuthModule {
	authService := service.NewAuthService(loginMethodService,userService ,jwtService)
	return &AuthModule{
		AuthService: authService,
	}
}