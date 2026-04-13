package handler

import (
	"context"
	"go-ecommerce/api-gateway/internal/dto"
	"go-ecommerce/common/gen-proto/otp"
	util "go-ecommerce/common/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	client otp.OtpServiceClient
}

func NewOtpHandler(client otp.OtpServiceClient) *OtpHandler {
	return &OtpHandler{
		client: client,
	}
}

func (otpHandler *OtpHandler) ValidateCreateAccount(c *gin.Context) {
	var input dto.ValidateCreateAccount
	if err := c.ShouldBindJSON(&input); err != nil {
		util.NewBindingError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	inputGrpc := &otp.ValidateCreateAccountDto{
		Email: input.Email,
		Otp: input.Otp,
	}
	res, err := otpHandler.client.ValidateCreateAccount(ctx, inputGrpc)
	if err != nil {
		util.NewResponseError(c, err)
		return
	}

	util.NewResponseData(c, http.StatusOK, util.Success, "Validate otp successfully", res)
}