package routes

import (
	"go-ecommerce/api-gateway/internal/handler"

	"github.com/gin-gonic/gin"
)

type OtpRoutes struct {
	handler *handler.OtpHandler
}

func NewOtpRoutes(hdl *handler.OtpHandler) *OtpRoutes{
	return &OtpRoutes{
		handler: hdl,
	}
}

func (otpRoute *OtpRoutes) RegisterRoutes(r *gin.RouterGroup) {
	otp := r.Group("/otp")
	{
		otp.POST("/validate-create-account", otpRoute.handler.ValidateCreateAccount)
	}
}