package main

import (
	"fmt"
	product "go-ecommerce/common/gen-proto/products"
	util "go-ecommerce/common/utils"
	"go-ecommerce/product-service/internal/config"
	"go-ecommerce/product-service/internal/db"
	"go-ecommerce/product-service/internal/module"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	util.LoadEnv()

	cfg := config.NewProductServiceConfig()
	db := db.ConnectDB(cfg.DatabaseConfig.MongoUri, cfg.DatabaseConfig.MongoDBName)

	// load module
	productModule := module.NewProductModule(db)

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