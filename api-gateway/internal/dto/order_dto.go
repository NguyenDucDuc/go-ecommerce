package dto

import "go.mongodb.org/mongo-driver/v2/bson"

type ItemDto struct {
	ProductId bson.ObjectID `json:"product_id"`
	Quantity int32 `json:"quantity"`
}

type CreateOrderDto struct {
	ShippingAddress string `json:"shipping_address"`
	Items []ItemDto `json:"items"`
}

type OrderItemResponse struct {
	ProductId bson.ObjectID `json:"product_id"`
	ProductName string `json:"product_name"`
	Price string `json:"price"`
}

type OrderResponse struct {
	Id bson.ObjectID `json:"_id"`
	OrderCode string `json:"order_code"`
	UserId bson.ObjectID `json:"user_id"`
	TotalAmount string `json:"total_amount"`
	ShippingAddress string `json:"shipping_address"`
	Items []OrderItemResponse `json:"items"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
