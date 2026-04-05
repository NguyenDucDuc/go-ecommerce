package handler

import (
	"context"
	"go-ecommerce/api-gateway/internal/dto"
	"go-ecommerce/common/gen-proto/auth"
	util "go-ecommerce/common/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthHandler struct {
	client auth.AuthServiceClient
}

func NewAuthHandler(client auth.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func (authHandler *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginDto
	if err := c.ShouldBindJSON(&input); err != nil {
		util.NewBindingError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := authHandler.client.Login(ctx, &auth.LoginDto{Email: input.Email, Password: input.Password})
	if err != nil {
		util.NewResponseError(c, err)
		return
	}

	id, _ := bson.ObjectIDFromHex(res.User.Id)

	rsp := dto.LoginResponse{
		User: dto.UserResponseDto{
			ID: id,
			Email: res.User.Email,
			FullName: res.User.FullName,
			Address: res.User.Address,
			Roles: res.User.Roles,
			CreatedAt: res.User.CreatedAt.AsTime().Format(time.RFC3339),
		},
		AccessToken: res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	util.NewResponseData(c, http.StatusOK, util.Success, "Login successfully", rsp)
}