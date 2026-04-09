package routes

import (
	"go-ecommerce/api-gateway/internal/handler"
	"go-ecommerce/common/middleware"
	"go-ecommerce/common/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type OrderRoutes struct {
	handler *handler.OrderHandler
	jwtService jwt.IJwtService
}


func NewOrderRoutes(handler *handler.OrderHandler, jwtService jwt.IJwtService) *OrderRoutes{
	return &OrderRoutes{
		handler: handler,
		jwtService: jwtService,
	}
}

func (orderRoutes *OrderRoutes) RegisterRoutes(r *gin.RouterGroup) {
	authMiddleware := middleware.AuthMiddleware(orderRoutes.jwtService)
	orders := r.Group("/orders")
	orders.Use(authMiddleware)
	{
		orders.POST("", orderRoutes.handler.Create)
	}
}