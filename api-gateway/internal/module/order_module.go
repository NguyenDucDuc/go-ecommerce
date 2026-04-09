package module

import (
	"go-ecommerce/api-gateway/internal/handler"
	"go-ecommerce/api-gateway/internal/routes"
	order "go-ecommerce/common/gen-proto/orders"
	"go-ecommerce/common/pkg/jwt"
)

type OrderModule struct {
	Handler *handler.OrderHandler
	Routes  *routes.OrderRoutes
}

func NewOrderModule(client order.OrderServiceClient, jwtService jwt.IJwtService) *OrderModule {
	handler := handler.NewOrderHandler(client)
	routes := routes.NewOrderRoutes(handler, jwtService)
	return &OrderModule{
		Handler: handler,
		Routes:  routes,
	}
}