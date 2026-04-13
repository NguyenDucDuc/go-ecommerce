package service

import (
	"context"
	"go-ecommerce/common/gen-proto/otp"
	pkg_redis "go-ecommerce/common/pkg/redis"
	util "go-ecommerce/common/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OtpService struct {
	loginMethodService *LoginMethodService
	userService        *UserService
	redisService 	pkg_redis.IRedisService
	otp.UnimplementedOtpServiceServer
}

func NewOtpService(loginMethodService *LoginMethodService, userService *UserService, redisService pkg_redis.IRedisService) *OtpService {
	return &OtpService{
		loginMethodService: loginMethodService,
		userService:        userService,
		redisService: redisService,
	}
}

func (otpService *OtpService) ValidateCreateAccount(ctx context.Context, input *otp.ValidateCreateAccountDto) (*otp.ValidateResponse, error) {
	loginMethod, err := otpService.loginMethodService.FindOne(ctx, bson.M{"email": input.Email})
	if err != nil {
		return &otp.ValidateResponse{}, err
	}

	if loginMethod.IsActive == true {
		return &otp.ValidateResponse{}, util.NewAppError(http.StatusBadRequest, util.ErrBadRequest, "Account already active")
	}

	// cache otp
	cacheKey := "otp:" + input.Email
	otpCache, err := otpService.redisService.GetString(ctx, cacheKey)
	if err != nil {
		return &otp.ValidateResponse{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Otp is not valid")
	}
	
	if otpCache != input.Otp {
		return &otp.ValidateResponse{}, util.NewAppError(http.StatusBadRequest, util.ErrBadRequest, "Otp is not valid")
	}

	// update status
	err = otpService.loginMethodService.UpdateOne(ctx, bson.M{"email": input.Email}, bson.M{"is_active": true})
	if err != nil {
		return &otp.ValidateResponse{}, err
	}

	return &otp.ValidateResponse{
		IsValid: true,
	}, nil
}