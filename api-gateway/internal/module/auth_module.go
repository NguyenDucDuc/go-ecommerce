package module

import (
	"go-ecommerce/api-gateway/internal/handler"
	"go-ecommerce/api-gateway/internal/routes"
	"go-ecommerce/common/gen-proto/auth"
)

type AuthModule struct {
	Handler *handler.AuthHandler
	Routes *routes.AuthRoutes
}

func NewAuthModule(client auth.AuthServiceClient) *AuthModule{
	handler := handler.NewAuthHandler(client)
	routes := routes.NewAuthRoutes(handler)
	return &AuthModule{
		Handler: handler,
		Routes: routes,
	}
}