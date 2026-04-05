package handler

import (
	"context"
	"go-ecommerce/api-gateway/internal/dto"
	product "go-ecommerce/common/gen-proto/products"
	util "go-ecommerce/common/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	client product.ProductServiceClient
}

func NewProductHandler(client product.ProductServiceClient) *ProductHandler {
	return &ProductHandler{
		client: client,
	}
}

func (handler *ProductHandler) Create(c *gin.Context) {
	var input dto.CreateProductDto
	if err := c.ShouldBindJSON(&input); err != nil {
		util.NewBindingError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	createProductDto := &product.CreateProductDto{
		Name: input.Name,
		Price: input.Price,
		Images: input.Images,
		Attributes: input.Attributes,
		Quantity: input.Quantity,
	}
	res, err := handler.client.CreateProduct(ctx, createProductDto)
	if err != nil {
		util.NewResponseError(c, err)
	}

	rsp := dto.ProductResponse{
		ID: res.Id,
		Name: res.Name,
		Price: res.Price,
		Images: res.Images,
		CreatedAt: res.CreatedAt.AsTime().Format(time.RFC3339),
		UpdatedAt: res.UpdatedAt.AsTime().Format(time.RFC3339),
	}

	util.NewResponseData(c, http.StatusOK, util.Success, "Create product successfully", rsp)
}