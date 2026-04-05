package module

import (
	"go-ecommerce/api-gateway/internal/handler"
	"go-ecommerce/api-gateway/internal/routes"
	user "go-ecommerce/common/gen-proto/users"
)

type UserModule struct {
	Routes *routes.UserRoutes
	Handler *handler.UserHandler
}

func NewUserModule(client user.UserServiceClient) *UserModule {
	userHandler := handler.NewUserHandler(client)
	userRoute := routes.NewUserRoutes(userHandler)
	return &UserModule{
		Handler: userHandler,
		Routes: userRoute,
	}
}