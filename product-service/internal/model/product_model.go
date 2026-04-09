package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID         bson.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name       string          `bson:"name" json:"name"`
	Price      bson.Decimal128 `bson:"price" json:"price"`
	Attributes bson.M          `bson:"attributes" json:"attributes"`
	Images     []string        `bson:"images" json:"images"`
	CreatedAt  time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `bson:"updated_at" json:"updated_at"`
}

type InventoryInfo struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	AvailableStock int64 `bson:"available_stock" json:"available_stock"`
	ReservedStock int64 `bson:"reversed_stock" json:"reversed_stock"`
}

type ProductWithInventory struct {
	ID         bson.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name       string          `bson:"name" json:"name"`
	Price      bson.Decimal128 `bson:"price" json:"price"`
	Attributes bson.M          `bson:"attributes" json:"attributes"`
	Images     []string        `bson:"images" json:"images"`
	InventoryInfo *InventoryInfo `bson:"inventory_info" json:"inventory_info"`
	CreatedAt  time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `bson:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	ProductId bson.ObjectID `json:"product_id"`
	Quantity int32 `json:"quantity"`
}