package dto

import "go.mongodb.org/mongo-driver/v2/bson"

type CreateUserDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

type CreateUserResponse struct {
	Id bson.ObjectID `json:"_id"`
	Otp string `json:"otp"`
}

type UserResponseDto struct {
	ID bson.ObjectID `json:"_id"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
	Address string `json:"address"`
	Roles []string `json:"roles"` 
	CreatedAt string   `json:"created_at"`
}
