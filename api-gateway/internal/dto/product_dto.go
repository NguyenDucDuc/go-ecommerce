package dto

import (
	product "go-ecommerce/common/gen-proto/products"
	"time"

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

type OrderBy string 
type Sort string 

var (
	CreatedAt OrderBy = "created_at"
	UpdatedAt OrderBy = "updated_at"

	Desc Sort = "desc"
	Asc Sort = "asc"
)

type GetListProductDto struct {
	Page int64 `form:"page" json:"page"`
	Limit int64 `form:"limit" json:"limit"`
	OrderBy OrderBy `form:"order_by" json:"order_by"`
	Sort Sort `form:"sort" json:"sort"`
}

type GetListProductResponse struct {
	Items []*ProductResponse `json:"items"`
	Total int `json:"total"`
	Page int `json:"page"`
	Limit int `json:"limit"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}

func MapToListProductResponse(responseGRPC *product.ListProductResponse) *GetListProductResponse{
	itemsGrpc := responseGRPC.Items
	prodRsp := make([]*ProductResponse, len(itemsGrpc))
	for i := range itemsGrpc {
		r := &ProductResponse{
			ID: itemsGrpc[i].Id,
			Name: itemsGrpc[i].Name,
			Price: itemsGrpc[i].Price,
			Attributes: itemsGrpc[i].Attributes,
			Images: itemsGrpc[i].Images,
			CreatedAt: itemsGrpc[i].CreatedAt.AsTime().Format(time.RFC3339),
			UpdatedAt: itemsGrpc[i].UpdatedAt.AsTime().Format(time.RFC3339),
		}
		prodRsp[i] = r
	}

	return &GetListProductResponse{
		Items: prodRsp,
		Total: int(responseGRPC.Total),
		Page: int(responseGRPC.Page),
		Limit: int(responseGRPC.Limit),
		HasNext: responseGRPC.HasPrev,
		HasPrev: responseGRPC.HasPrev,
	}
}