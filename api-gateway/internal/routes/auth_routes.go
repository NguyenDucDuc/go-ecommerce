package routes

import (
	"go-ecommerce/api-gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	handler *handler.AuthHandler
}

func NewAuthRoutes(handler *handler.AuthHandler) *AuthRoutes{
	return &AuthRoutes{
		handler: handler,
	}
}

func (authRoutes *AuthRoutes) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", authRoutes.handler.Login)
	}
}