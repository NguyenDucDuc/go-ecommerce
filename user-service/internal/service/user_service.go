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
	loginMethodService *LoginMethodService
	rabbitService rabbitmq.IRabbitMQService
	user.UnimplementedUserServiceServer
}

func NewUserService(repo repository.IUserRepository, loginMethodService *LoginMethodService ,rabbitService rabbitmq.IRabbitMQService) *UserService {
	return &UserService{
		repo: repo,
		loginMethodService: loginMethodService,
		rabbitService: rabbitService,
	}
}

func (userService *UserService) CreateUser(ctx context.Context, input *user.CreateUserDto) (*user.CreateUserResponse, error) {
	// check email existed
	existedUser, _ := userService.repo.FindByEmail(ctx, input.Email)
	if existedUser.ID != bson.NilObjectID {
		return &user.CreateUserResponse{}, util.NewAppError(http.StatusConflict, util.ErrConflict, "Email already exist")
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