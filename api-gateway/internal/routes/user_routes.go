package routes

import (
	"go-ecommerce/api-gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

func NewUserRoutes(hdl *handler.UserHandler) *UserRoutes{
	return &UserRoutes{
		handler: hdl,
	}
}

func (userRoute *UserRoutes) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", userRoute.handler.CreateUser)
	}
}