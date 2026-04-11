package main

import (
	"fmt"
	product "go-ecommerce/common/gen-proto/products"
	"go-ecommerce/common/pkg/rabbitmq"
	pkg_redis "go-ecommerce/common/pkg/redis"
	util "go-ecommerce/common/utils"
	product_config "go-ecommerce/product-service/internal/config"
	"go-ecommerce/product-service/internal/db"
	"go-ecommerce/product-service/internal/module"
	pkg_product_redis "go-ecommerce/product-service/internal/pkg/redis"
	"go-ecommerce/product-service/internal/worker"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	util.LoadEnv()

	cfg := product_config.NewProductServiceConfig()
	// mongodb
	db := db.ConnectDB(cfg.DatabaseConfig.MongoUri, cfg.DatabaseConfig.MongoDBName)
	// redisdb
	rdb := pkg_product_redis.ConnectRedis(cfg.RedisConfig)
	redisService := pkg_redis.NewRedisService(rdb)

	// rabbit mq
	rabbitMQService, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQConfig.Uri)
	if err != nil {
		log.Fatal(err)
	}

	// load module
	productModule := module.NewProductModule(db, redisService, rabbitMQService)


	// rabbit mq worker
	productWorker := worker.NewProductWorker(rabbitMQService, productModule.Service)
	productWorker.Start()

	// gRPC setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(cfg.GrpcPort)))
	if err != nil {
		log.Fatalf("❌ Lỗi lắng nghe port 50002: %v", err)
	}
	// gRPC service setup
	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(grpcServer, productModule.Service)

	reflection.Register(grpcServer)

	log.Println("🚀 Product gRPC Service đang chạy tại port :50002")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Lỗi khởi chạy gRPC Server: %v", err)
	}
}