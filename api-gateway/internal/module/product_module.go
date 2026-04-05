package module

import (
	"go-ecommerce/api-gateway/internal/handler"
	"go-ecommerce/api-gateway/internal/routes"
	product "go-ecommerce/common/gen-proto/products"
)

type ProductModule struct {
	Handler *handler.ProductHandler
	Routes *routes.ProductRoutes
}


func NewProductModule(client product.ProductServiceClient) *ProductModule {
	handler := handler.NewProductHandler(client)
	routes := routes.NewProductRoutes(handler)
	return &ProductModule{
		Handler: handler,
		Routes: routes,
	}
}