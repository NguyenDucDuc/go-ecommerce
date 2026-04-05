package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Inventory struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID      bson.ObjectID `bson:"product_id" json:"product_id"`
	AvailableStock int32         `bson:"available_stock" json:"available_stock"`
	ReservedStock  int32         `bson:"reserved_stock" json:"reserved_stock"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}