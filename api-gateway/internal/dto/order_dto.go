package dto

import (
	order "go-ecommerce/common/gen-proto/orders"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ItemDto struct {
	ProductId bson.ObjectID `json:"product_id"`
	Quantity int32 `json:"quantity"`
}

type CreateOrderDto struct {
	ShippingAddress string `json:"shipping_address"`
	Items []ItemDto `json:"items"`
}

type ItemResponse struct {
	ProductId bson.ObjectID `json:"product_id"`
	ProductName string `json:"product_name"`
	Price string `json:"price"`
}

type OrderResponse struct {
	Id bson.ObjectID `json:"_id"`
	OrderCode string `json:"order_code"`
	UserId bson.ObjectID `json:"user_id"`
	TotalAmount string `json:"total_amount"`
	Status string `json:"status"`
	ShippingAddress string `json:"shipping_address"`
	Items []*ItemResponse `json:"items"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func MapToItemResponse(items []*order.OrderItem) []*ItemResponse {
	rsp := make([]*ItemResponse, len(items))
	for i, item := range items {
		pId, _ := bson.ObjectIDFromHex(item.ProductId)
		itemRsp := &ItemResponse{
			ProductId: pId,
			ProductName: item.ProductName,
			Price: item.Price,
		}
		rsp[i] = itemRsp
	}

	return rsp
}
