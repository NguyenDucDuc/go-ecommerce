package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	// Email     string        `bson:"email" json:"email"`
	// Password  string        `bson:"password" json:"password"` // Không trả password về client
	FullName  string        `bson:"full_name" json:"full_name"`
	Roles     []string      `bson:"roles" json:"roles"`
	Address string     `bson:"address" json:"address"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}