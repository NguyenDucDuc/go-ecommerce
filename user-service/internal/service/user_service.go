package service

import (
	"context"
	user "go-ecommerce/common/gen-proto/users"
	"go-ecommerce/common/pkg/rabbitmq"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/internal/model"
	"go-ecommerce/user-service/internal/repository"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	repo repository.IUserRepository
	rabbitService rabbitmq.IRabbitMQService
	user.UnimplementedUserServiceServer
}

func NewUserService(repo repository.IUserRepository, rabbitService rabbitmq.IRabbitMQService) *UserService {
	return &UserService{
		repo: repo,
		rabbitService: rabbitService,
	}
}

func (userService *UserService) CreateUser(ctx context.Context, input *user.CreateUserDto) (*user.UserResponse, error) {
	// check email existed
	existedUser, _ := userService.repo.FindByEmail(ctx, input.Email)
	if existedUser.ID != bson.NilObjectID {
		return &user.UserResponse{}, util.NewAppError(http.StatusConflict, util.ErrConflict, "Email already exist")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return &user.UserResponse{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Failed to hash password")
	}
	userModel := model.User{
		Email: input.Email,
		Password: string(hashPassword),
		Address: input.Address,
		FullName: input.FullName,
		Roles: []string{"CUSTOMER"},
	}

	res, err := userService.repo.Create(ctx, &userModel)
	if err != nil {
		return &user.UserResponse{}, err
	}

	userResponse := &user.UserResponse{
		Id: res.ID.Hex(),
		Email: res.Email,
		Password: res.Password,
		FullName: res.FullName,
		Address: res.Address,
		Roles: res.Roles,
		CreatedAt: timestamppb.New(res.CreatedAt),
	}

	// send mail
	otp, _ := util.GenerateOTP()
	msg := map[string]interface{}{
		"otp": otp,
		"email" : userResponse.Email,
	}
	userService.rabbitService.Publish("topic_exchange","user.created",msg)

	return userResponse, nil
}

func (userService *UserService) FindByEmail(ctx context.Context, input *user.FindByEmailDto) (*user.UserResponse, error) {
	res, err := userService.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return &user.UserResponse{}, err
	}
	rsp := &user.UserResponse{
		Id: res.ID.Hex(),
		Email: res.Email,
		FullName: res.FullName,
		Address: res.Address,
		CreatedAt: timestamppb.New(res.CreatedAt),
	}
	return rsp, nil
}