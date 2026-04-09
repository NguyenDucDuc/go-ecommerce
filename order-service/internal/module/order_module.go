package module

import (
	product "go-ecommerce/common/gen-proto/products"
	"go-ecommerce/common/pkg/rabbitmq"
	pkg_redis "go-ecommerce/common/pkg/redis"
	"go-ecommerce/order-service/internal/repository"
	"go-ecommerce/order-service/internal/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderModule struct {
	Repo repository.IOrderRepository
	Service *service.OrderService
}

func NewOrderModule(db *mongo.Database, redisService pkg_redis.IRedisService, productClient product.ProductServiceClient, rabbitMQService rabbitmq.IRabbitMQService) *OrderModule {
	repo := repository.NewOrderRepository(db)
	service := service.NewOrderService(repo, redisService, productClient, rabbitMQService)
	return &OrderModule{
		Repo: repo,
		Service: service,
	}
}