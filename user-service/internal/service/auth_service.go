package service

import (
	"context"
	"go-ecommerce/common/gen-proto/auth"
	user "go-ecommerce/common/gen-proto/users"
	"go-ecommerce/common/pkg/jwt"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/internal/repository"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthService struct {
	repo repository.IUserRepository
	jwtService jwt.IJwtService
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(repo repository.IUserRepository, jwtService jwt.IJwtService) *AuthService {
	return &AuthService{
		repo: repo,
		jwtService: jwtService,
	}
}

func (authService *AuthService) Login(ctx context.Context, in *auth.LoginDto) (*auth.LoginResponse, error) {
	userFound, err := authService.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		return &auth.LoginResponse{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(in.Password))
	if err != nil {
		return &auth.LoginResponse{}, util.NewAppError(http.StatusBadRequest, util.ErrBadRequest, "Email or password not valid")
	}

	accessToken, err := authService.jwtService.GenerateAccessToken(userFound.ID.Hex(), userFound.Email, userFound.Roles)
	if err != nil {
		return &auth.LoginResponse{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Jwt error")
	}

	refreshToken, err := authService.jwtService.GenerateRefreshToken(userFound.ID.Hex())
	if err != nil {
		return &auth.LoginResponse{}, util.NewAppError(http.StatusInternalServerError, util.ErrInternalServer, "Jwt error")
	}

	rsp := &auth.LoginResponse{
		User: &user.UserResponse{
			Id: userFound.ID.Hex(),
			Email: userFound.Email,
			Password: userFound.Password,
			FullName: userFound.FullName,
			Address: userFound.Address,
			Roles: userFound.Roles,
			CreatedAt: timestamppb.New(userFound.CreatedAt),
		},
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return rsp, nil

}