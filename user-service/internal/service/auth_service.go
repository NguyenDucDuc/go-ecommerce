package service

import (
	"context"
	"go-ecommerce/common/gen-proto/auth"
	user "go-ecommerce/common/gen-proto/users"
	"go-ecommerce/common/pkg/jwt"
	util "go-ecommerce/common/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthService struct {
	loginMethodService *LoginMethodService
	userService *UserService
	jwtService jwt.IJwtService
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(loginMethodService *LoginMethodService, userService *UserService, jwtService jwt.IJwtService) *AuthService {
	return &AuthService{
		loginMethodService: loginMethodService,
		userService: userService,
		jwtService: jwtService,
	}
}

func (authService *AuthService) Login(ctx context.Context, in *auth.LoginDto) (*auth.LoginResponse, error) {
	filter := bson.M{"email": in.Email}
	loginMethod, err := authService.loginMethodService.FindOne(ctx, filter)
	if err != nil {
		return &auth.LoginResponse{}, err
	}

	if loginMethod.IsActive == false {
		return &auth.LoginResponse{}, util.NewAppError(http.StatusForbidden, util.ErrForbidden, "Account is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(loginMethod.Password), []byte(in.Password))
	if err != nil {
		return &auth.LoginResponse{}, util.NewAppError(http.StatusBadRequest, util.ErrBadRequest, "Email or password not valid")
	}

	userFound, err := authService.userService.FindByEmail(ctx, &user.FindByEmailDto{Email: in.Email})

	accessToken, err := authService.jwtService.GenerateAccessToken(userFound.Id, loginMethod.Email, userFound.Roles)
	if err != nil {
		return &auth.LoginResponse{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Jwt error")
	}

	refreshToken, err := authService.jwtService.GenerateRefreshToken(userFound.Id)
	if err != nil {
		return &auth.LoginResponse{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Jwt error")
	}

	rsp := &auth.LoginResponse{
		User: &user.UserResponse{
			Id: userFound.Id,
			Email: loginMethod.Email,
			Password: loginMethod.Password,
			FullName: userFound.FullName,
			Address: userFound.Address,
			Roles: userFound.Roles,
			CreatedAt: timestamppb.New(userFound.CreatedAt.AsTime()),
		},
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return rsp, nil

}