package routes

import (
	"go-ecommerce/api-gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

type ProductRoutes struct {
	handler *handler.ProductHandler
}

func NewProductRoutes(handler *handler.ProductHandler) *ProductRoutes{
	return &ProductRoutes{
		handler: handler,
	}
}

func (productRoutes *ProductRoutes) RegisterRoutes(r *gin.RouterGroup) {
	products := r.Group("/products")
	{
		products.POST("", productRoutes.handler.Create)
		products.GET("", productRoutes.handler.GetList)
	}
}