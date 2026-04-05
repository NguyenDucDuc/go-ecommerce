package service

import (
	"context"
	product "go-ecommerce/common/gen-proto/products"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/internal/model"
	"go-ecommerce/user-service/internal/repository"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	repo repository.IProductRepository
	inventoryService *InventoryService
	product.UnimplementedProductServiceServer
}


func NewProductService(repo repository.IProductRepository, inventoryService *InventoryService) *ProductService{
	return &ProductService{
		repo: repo,
		inventoryService: inventoryService,
	}
}

func (productService *ProductService) CreateProduct(ctx context.Context, input *product.CreateProductDto) (*product.ProductResponse, error) {
	// create product
	productModel := model.Product{
		Name: input.Name,
		Price: util.ToDecimal128(input.Price),
		Attributes: input.Attributes.AsMap(),
		Images: input.Images,
	}
	res, err := productService.repo.Create(ctx, &productModel)
	if err != nil {
		return &product.ProductResponse{}, err
	}
	// create inventory
	inventoryModel := model.Inventory{
		ProductID: res.ID,
		AvailableStock: int32(input.Quantity),
		ReservedStock: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = productService.inventoryService.Create(ctx, &inventoryModel)
	if err != nil {
		return &product.ProductResponse{}, err
	}

	
	rsp := &product.ProductResponse{
		Id: res.ID.Hex(),
		Name: res.Name,
		Price: res.Price.String(),
		Attributes: util.MapToProtoStruct(res.Attributes),
		CreatedAt: timestamppb.New(res.CreatedAt),
		UpdatedAt: timestamppb.New(res.UpdatedAt),
	}
	return rsp, nil
}