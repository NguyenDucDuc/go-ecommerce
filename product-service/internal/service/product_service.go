package service

import (
	"context"
	"fmt"
	product "go-ecommerce/common/gen-proto/products"
	pkg_redis "go-ecommerce/common/pkg/redis"
	util "go-ecommerce/common/utils"
	"go-ecommerce/product-service/internal/db"
	"go-ecommerce/product-service/internal/model"
	"go-ecommerce/product-service/internal/repository"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	txManager        db.TransactionManager
	repo repository.IProductRepository
	inventoryService *InventoryService
	redisService pkg_redis.IRedisService
	product.UnimplementedProductServiceServer
}


func NewProductService(tx db.TransactionManager, repo repository.IProductRepository, inventoryService *InventoryService, rdbService pkg_redis.IRedisService) *ProductService{
	return &ProductService{
		repo: repo,
		inventoryService: inventoryService,
		redisService: rdbService,
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

func (productService *ProductService) GetListProduct(ctx context.Context, input *product.GetListProductDto) (*product.ListProductResponse, error) {
	if input.Page <= 0 {
        input.Page = 1
    }
    if input.Limit <= 0 || input.Limit > 100 { // Giới hạn tối đa để tránh kéo quá nhiều data
        input.Limit = 10
    }
	skip := (input.Page - 1) * input.Limit
	var itemsCache *product.ListProductResponse
	cacheKey := fmt.Sprintf("products:list:page:%d:limit:%d:order:%s:sort:%s", 
		input.Page, 
		input.Limit, 
		input.OrderBy, 
		input.Sort,
	)

	err := productService.redisService.GetJSON(ctx, cacheKey, &itemsCache)
	if err == nil {
		log.Println("Cache product successfully")
		return itemsCache, nil
	}


	items, total ,err := productService.repo.FindAll(ctx, int(skip), int(input.Limit), input.OrderBy, input.Sort)
	if err != nil {
		return &product.ListProductResponse{}, err
	}

	prodRsp := make([]*product.ProductResponse, len(items))
	for i := range items {
		r := &product.ProductResponse{
			Id: items[i].ID.Hex(),
			Name: items[i].Name,
			Price: items[i].Price.String(),
			Attributes: util.MapToProtoStruct(items[i].Attributes),
			InventoryInfo: &product.InventoryResponse{
				Id: items[i].InventoryInfo.ID.Hex(),
				AvailableStock: items[i].InventoryInfo.AvailableStock,
				ReservedStock: items[i].InventoryInfo.ReservedStock,
			},
			Images: items[i].Images,
			CreatedAt:  timestamppb.New(items[i].CreatedAt),
			UpdatedAt:  timestamppb.New(items[i].UpdatedAt),
		}
		prodRsp[i] = r
	}

	rsp := &product.ListProductResponse{
		Items: prodRsp,
		Total: int64(total),
		Page: input.Page,
		Limit: input.Limit,
		HasNext: int64(input.Page) * int64(input.Limit) < int64(total),
		HasPrev: int64(input.Limit) > 1,
	}
	// cache redis
	productService.redisService.SetJSON(ctx, cacheKey, rsp, 10 * time.Minute)
	return rsp, nil
}