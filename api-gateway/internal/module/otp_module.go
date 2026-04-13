package module

import (
	"go-ecommerce/api-gateway/internal/handler"
	"go-ecommerce/api-gateway/internal/routes"
	"go-ecommerce/common/gen-proto/otp"
)

type OtpModule struct {
	Handler *handler.OtpHandler
	Routes *routes.OtpRoutes
}

func NewOtpModule(client otp.OtpServiceClient) *OtpModule{
	handler := handler.NewOtpHandler(client)
	routes := routes.NewOtpRoutes(handler)
	return  &OtpModule{
		Handler: handler,
		Routes: routes,
	}
}