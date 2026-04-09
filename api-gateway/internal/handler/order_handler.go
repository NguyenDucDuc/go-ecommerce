package handler

import (
	"context"
	"go-ecommerce/api-gateway/internal/dto"
	order "go-ecommerce/common/gen-proto/orders"
	util "go-ecommerce/common/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderHandler struct {
	client order.OrderServiceClient
}

func NewOrderHandler(client order.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		client: client,
	}
}

func (orderHandler *OrderHandler) Create(c *gin.Context) {
	var input dto.CreateOrderDto
	if err := c.ShouldBindJSON(&input); err != nil {
		util.NewBindingError(c, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	userId, ok := c.Get("user_id")
	if !ok {
		util.NewResponseError(c, util.NewAppError(http.StatusBadRequest, util.ErrBadRequest, "Invalid user_id from token"))
	}

	itemsGrpc := make([]*order.CreateOrderDto_ItemDto, len(input.Items))
	for i, item := range input.Items {
		itemGrpc := &order.CreateOrderDto_ItemDto{
			ProductId: item.ProductId.Hex(),
			Quantity: item.Quantity,
		}
		itemsGrpc[i] = itemGrpc
	}
	inputGrpc := &order.CreateOrderDto{
		UserId: userId.(string),
		ShippingAddress: input.ShippingAddress,
		Items: itemsGrpc,
	}

	res, err := orderHandler.client.CreateOrder(ctx, inputGrpc)
	if err != nil {
		util.NewResponseError(c, err)
		return
	}

	oId, _ := bson.ObjectIDFromHex(res.Id)
	uId, _ := bson.ObjectIDFromHex(res.UserId)
	rsp := &dto.OrderResponse{
		Id: oId,
		OrderCode: res.OrderCode,
		UserId: uId,
		TotalAmount: res.TotalAmount,
		ShippingAddress: res.ShippingAddress,
		Status: res.Status,
		Items: dto.MapToItemResponse(res.Items),
		CreatedAt: res.CreatedAt.AsTime().Format(time.RFC3339),
		UpdatedAt: res.UpdatedAt.AsTime().Format(time.RFC3339),
	}

	util.NewResponseData(c, http.StatusOK, util.Success, "Create order successfully", rsp)
}