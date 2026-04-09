package main

import (
	"fmt"
	order "go-ecommerce/common/gen-proto/orders"
	product "go-ecommerce/common/gen-proto/products"
	"go-ecommerce/common/pkg/rabbitmq"
	pkg_redis "go-ecommerce/common/pkg/redis"
	util "go-ecommerce/common/utils"
	order_config "go-ecommerce/order-service/internal/config"
	"go-ecommerce/order-service/internal/db"
	"go-ecommerce/order-service/internal/module"
	pkg_order_redis "go-ecommerce/order-service/pkg/redis"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	util.LoadEnv()
	cfg := order_config.NewProductServiceConfig()

	// setup grpc client
	prodConn, err := grpc.NewClient(cfg.GrpcConfig.GrpcProductUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Không thể kết nối tới Product Service: %v", err)
	}
	defer prodConn.Close()
	productClient := product.NewProductServiceClient(prodConn)

	// mongodb
	db := db.ConnectDB(cfg.DatabaseConfig.MongoUri, cfg.DatabaseConfig.MongoDBName)
	// redisdb
	rdb := pkg_order_redis.ConnectRedis(cfg.RedisConfig)
	redisService := pkg_redis.NewRedisService(rdb)
	// rabbit mq
	rabbitService, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQConfig.Uri)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitService.Close()

	// load module
	orderModule := module.NewOrderModule(db, redisService, productClient, rabbitService)

	// rabbit mq worker
	go func() {
        log.Println("[*] Order Service đang đợi phản hồi từ Inventory...")
        
        // Nghe khi trừ kho thành công
        go rabbitService.Consume(
            "order_inventory_success_queue", 
            "inventory.success", 
            "order_exchange", 
            orderModule.Service.HandleInventorySuccess, // Hàm này bạn viết trong service
        )

        // Nghe khi trừ kho thất bại
        go rabbitService.Consume(
            "order_inventory_failed_queue", 
            "inventory.failed", 
            "order_exchange", 
            orderModule.Service.HandleInventoryFailed, // Hàm này bạn viết trong service
        )
    }()

	// gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(cfg.GrpcConfig.GrpcPort)))
	if err != nil {
		log.Fatalf("❌ Lỗi lắng nghe port 50003: %v", err)
	}
	// gRPC service setup
	grpcServer := grpc.NewServer()
	order.RegisterOrderServiceServer(grpcServer, orderModule.Service)

	reflection.Register(grpcServer)

	log.Println("🚀 Order gRPC Service đang chạy tại port :50003")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Lỗi khởi chạy gRPC Server: %v", err)
	}
}