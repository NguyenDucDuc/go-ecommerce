package module

import (
	"go-ecommerce/user-service/internal/repository"
	"go-ecommerce/user-service/internal/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductModule struct {
	Repo repository.IProductRepository
	Service *service.ProductService
}

func NewProductModule(db *mongo.Database) *ProductModule {
	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepo)
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, inventoryService)
	return &ProductModule{
		Repo: productRepo,
		Service: productService,
	}
}