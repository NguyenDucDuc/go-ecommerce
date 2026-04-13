package module

import (
	pkg_redis "go-ecommerce/common/pkg/redis"
	"go-ecommerce/user-service/internal/service"
)

type OtpModule struct {
	Service *service.OtpService
}

func NewOtpModule(loginMethodService *service.LoginMethodService, userService *service.UserService, redisService pkg_redis.IRedisService) *OtpModule {
	otpService := service.NewOtpService(loginMethodService, userService, redisService)
	return &OtpModule{
		Service: otpService,
	}
}