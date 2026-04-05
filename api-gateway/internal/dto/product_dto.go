package dto

import (
	"google.golang.org/protobuf/types/known/structpb"
)

type CreateProductDto struct {
	Name       string `json:"name"`
	Price      string `json:"price"`
	Attributes *structpb.Struct `json:"attributes"`
	Images []string `json:"images"`
	Quantity int64 `json:"quantity"`
}

type ProductResponse struct {
	ID string `json:"_id"`
	Name string `json:"name"`
	Price string `json:"price"`
	Attributes *structpb.Struct `json:"attributes"`
	Images []string `json:"images"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}