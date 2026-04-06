package service

import (
	"context"
	product "go-ecommerce/common/gen-proto/products"
	util "go-ecommerce/common/utils"
	"go-ecommerce/product-service/internal/db"
	"go-ecommerce/product-service/internal/model"
	"go-ecommerce/product-service/internal/repository"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	repo repository.IProductRepository
	inventoryService *InventoryService
	txManager        db.TransactionManager
	product.UnimplementedProductServiceServer
}


func NewProductService(tx db.TransactionManager, repo repository.IProductRepository, inventoryService *InventoryService) *ProductService{
	return &ProductService{
		repo: repo,
		inventoryService: inventoryService,
		txManager: tx,
	}
}

func (productService *ProductService) CreateProduct(ctx context.Context, input *product.CreateProductDto) (*product.ProductResponse, error) {
	// Khai báo biến để hứng kết quả từ trong closure của transaction
	var res *model.Product

	// Thực hiện toàn bộ logic trong một Transaction
	err := productService.txManager.WithTransaction(ctx, func(sessCtx context.Context) error {
		// 1. Khởi tạo và tạo Product
		productModel := model.Product{
			Name:       input.Name,
			Price:      util.ToDecimal128(input.Price),
			Attributes: input.Attributes.AsMap(),
			Images:     input.Images,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		// QUAN TRỌNG: Phải dùng sessCtx ở đây
		p, err := productService.repo.Create(sessCtx, &productModel)
		if err != nil {
			return err
		}
		res = p


		// 2. Khởi tạo và tạo Inventory
		inventoryModel := model.Inventory{
			ProductID:      res.ID,
			AvailableStock: int32(input.Quantity),
			ReservedStock:  0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// QUAN TRỌNG: Phải dùng sessCtx ở đây
		_, err = productService.inventoryService.Create(sessCtx, &inventoryModel)
		if err != nil {
			return err
		}

		return nil
	})

	// Nếu transaction thất bại (gồm cả việc tạo Product hoặc Inventory lỗi)
	if err != nil {
		return &product.ProductResponse{}, err
	}

	// 3. Mapping kết quả trả về sau khi Transaction thành công (Commit)
	rsp := &product.ProductResponse{
		Id:         res.ID.Hex(),
		Name:       res.Name,
		Price:      res.Price.String(),
		Attributes: util.MapToProtoStruct(res.Attributes),
		Images:     res.Images,
		CreatedAt:  timestamppb.New(res.CreatedAt),
		UpdatedAt:  timestamppb.New(res.UpdatedAt),
	}

	return rsp, nil
}