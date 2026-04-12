package handler

import (
	"context"
	"go-ecommerce/api-gateway/internal/dto"
	user "go-ecommerce/common/gen-proto/users"
	util "go-ecommerce/common/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserHandler struct {
	userClient user.UserServiceClient
}

func NewUserHandler(client user.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: client,
	}
}

func (handler *UserHandler) CreateUser(c *gin.Context){
	var req user.CreateUserDto
	if err := c.ShouldBindJSON(&req); err != nil {
		util.NewBindingError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := handler.userClient.CreateUser(ctx, &req)
	if err != nil {
		util.NewResponseError(c, err)
		return
	}

	id, _ := bson.ObjectIDFromHex(user.Id)
	rsp := &dto.CreateUserResponse{
		Id:        id, // Chuyển ObjectID sang String
		Otp: user.Otp,
	}

	util.NewResponseData(c, http.StatusCreated, util.Success ,"Create user successfully", rsp)
}