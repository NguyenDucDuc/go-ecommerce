package service

import (
	"context"
	user "go-ecommerce/common/gen-proto/users"
	"go-ecommerce/common/pkg/rabbitmq"
	pkg_redis "go-ecommerce/common/pkg/redis"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/internal/model"
	"go-ecommerce/user-service/internal/repository"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	repo repository.IUserRepository
	loginMethodService *LoginMethodService
	rabbitService rabbitmq.IRabbitMQService
	redisService pkg_redis.IRedisService
	user.UnimplementedUserServiceServer
}

func NewUserService(repo repository.IUserRepository, loginMethodService *LoginMethodService ,rabbitService rabbitmq.IRabbitMQService, redisService pkg_redis.IRedisService) *UserService {
	return &UserService{
		repo: repo,
		loginMethodService: loginMethodService,
		rabbitService: rabbitService,
		redisService: redisService,
	}
}

func (userService *UserService) CreateUser(ctx context.Context, input *user.CreateUserDto) (*user.CreateUserResponse, error) {
	// check email existed
	loginMethodExisted, _ := userService.loginMethodService.FindOne(ctx, bson.M{"email": input.Email})
	if loginMethodExisted.ID != bson.NilObjectID {
		if loginMethodExisted.IsActive == true {
			return &user.CreateUserResponse{}, util.NewAppError(http.StatusConflict, util.ErrConflict, "Email already exist")
		} else {
			return &user.CreateUserResponse{}, util.NewAppError(http.StatusConflict, util.ErrConflict, "Email already exist, please validate otp to active account")
		}
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return &user.CreateUserResponse{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed to hash password")
	}

	// login method
	loginMethod := model.LoginMethod{
		Email: input.Email,
		Password: string(hashPassword),
		IsActive: false,
	}

	userModel := model.User{
		Address: input.Address,
		FullName: input.FullName,
		Roles: []string{"CUSTOMER"},
	}
	// create user first
	resUser, err := userService.repo.Create(ctx, &userModel)
	if err != nil {
		return &user.CreateUserResponse{}, err
	}
	// create login method
	loginMethod.UserId = resUser.ID
	_, err = userService.loginMethodService.Create(ctx, &loginMethod)
	if err != nil {
		return &user.CreateUserResponse{}, err
	}

	userResponse := &user.UserResponse{
		Id: resUser.ID.Hex(),
		FullName: resUser.FullName,
		Address: resUser.Address,
		Roles: resUser.Roles,
		CreatedAt: timestamppb.New(resUser.CreatedAt),
	}

	// send mail
	otp, _ := util.GenerateOTP()
	msg := map[string]interface{}{
		"otp": otp,
		"email" : userResponse.Email,
	}
	userService.rabbitService.Publish("topic_exchange","user.created",msg)
	// cache otp to redis
	cacheKey := "otp:" + loginMethod.Email
	userService.redisService.SetString(ctx, cacheKey, otp, 3*time.Minute)

	rsp := &user.CreateUserResponse{
		Id: userResponse.Id,
		Otp: otp,
	}

	return rsp, nil
}

func (userService *UserService) FindByEmail(ctx context.Context, input *user.FindByEmailDto) (*user.UserResponse, error) {
	res, err := userService.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return &user.UserResponse{}, err
	}
	rsp := &user.UserResponse{
		Id: res.ID.Hex(),
		FullName: res.FullName,
		Address: res.Address,
		CreatedAt: timestamppb.New(res.CreatedAt),
	}
	return rsp, nil
}