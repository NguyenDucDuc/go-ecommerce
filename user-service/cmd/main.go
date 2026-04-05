package main

import (
	"fmt"
	"go-ecommerce/common/gen-proto/auth"
	product "go-ecommerce/common/gen-proto/products"
	user "go-ecommerce/common/gen-proto/users"
	"go-ecommerce/common/pkg/jwt"
	util "go-ecommerce/common/utils"
	"go-ecommerce/user-service/config"
	"go-ecommerce/user-service/db"
	"go-ecommerce/user-service/module"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	util.LoadEnv()

	cfg := config.NewUserServiceConfig()
	db := db.ConnectDB(cfg.DatabaseConfig.MongoUri, cfg.DatabaseConfig.MongoDBName)

	// init jwt
	jwtService := jwt.NewJWTService(cfg.JwtConfig)
	// load module
	userModule := module.NewUserModule(db)
	authModule := module.NewAuthModule(db, jwtService)
	productModule := module.NewProductModule(db)

	// gRPC setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(cfg.GrpcPort)))
	if err != nil {
		log.Fatalf("❌ Lỗi lắng nghe port 50001: %v", err)
	}
	// gRPC service setup
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userModule.Service)
	auth.RegisterAuthServiceServer(grpcServer, authModule.AuthService)
	product.RegisterProductServiceServer(grpcServer, productModule.Service)

	reflection.Register(grpcServer)

	log.Println("🚀 User gRPC Service đang chạy tại port :50001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Lỗi khởi chạy gRPC Server: %v", err)
	}
}