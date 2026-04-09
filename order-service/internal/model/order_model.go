package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusShipping  OrderStatus = "SHIPPING"
	StatusCancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	ProductID   bson.ObjectID `bson:"product_id" json:"product_id"`
	ProductName string             `bson:"product_name" json:"product_name"` // Snapshot name
	Price       bson.Decimal128            `bson:"price" json:"price"`             // Snapshot price
	Quantity    int                `bson:"quantity" json:"quantity"`
}

type Order struct {
	ID              bson.ObjectID      `bson:"_id,omitempty" json:"id"`
	OrderCode       string             `bson:"order_code" json:"order_code"`
	UserID          bson.ObjectID      `bson:"user_id" json:"user_id"`
	TotalAmount     bson.Decimal128            `bson:"total_amount" json:"total_amount"` // Decimal trong Mongo ánh xạ tốt với float64 hoặc decimal128
	Status          string        `bson:"status" json:"status"`
	Items     		[]OrderItem        `bson:"order_items" json:"order_items"` // Dùng interface{} hoặc struct riêng cho Object
	ShippingAddress string             `bson:"shipping_address" json:"shipping_address"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}