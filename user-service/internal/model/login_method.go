package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LoginMethod struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"password"` // Không trả password về client
	IsActive bool 			`bson:"is_active" json:"is_active"`
	UserId bson.ObjectID	`bson:"user_id" json:"user_id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}