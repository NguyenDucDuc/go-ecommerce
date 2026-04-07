package module

import (
	pkg_redis "go-ecommerce/common/pkg/redis"
	"go-ecommerce/product-service/internal/db"
	"go-ecommerce/product-service/internal/repository"
	"go-ecommerce/product-service/internal/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductModule struct {
	Repo repository.IProductRepository
	Service *service.ProductService
}

func NewProductModule(database *mongo.Database, rdbService pkg_redis.IRedisService) *ProductModule {
	tx := db.NewTransactionManager(database.Client())
	inventoryRepo := repository.NewInventoryRepository(database)
	inventoryService := service.NewInventoryService(inventoryRepo)
	productRepo := repository.NewProductRepository(database)
	productService := service.NewProductService(tx,productRepo, inventoryService, rdbService)
	return &ProductModule{
		Repo: productRepo,
		Service: productService,
	}
}