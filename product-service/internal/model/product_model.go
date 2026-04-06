package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID         bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name       string          `bson:"name" json:"name"`
	Price      bson.Decimal128 `bson:"price" json:"price"`
	Attributes bson.M          `bson:"attributes" json:"attributes"` // Flexible schema cực gọn với bson.M
	Images     []string        `bson:"images" json:"images"`
	CreatedAt  time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `bson:"updated_at" json:"updated_at"`
}