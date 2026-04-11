package main

import (
	"fmt"
	"go-ecommerce/common/gen-proto/auth"
	user "go-ecommerce/common/gen-proto/users"
	"go-ecommerce/common/pkg/jwt"
	"go-ecommerce/common/pkg/rabbitmq"
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

	// rabbit mq
	rabbitMQService, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQConfig.Uri)
	if err != nil {
		log.Fatal(err)
	}

	// init jwt
	jwtService := jwt.NewJWTService(cfg.JwtConfig.JwtSecret, cfg.JwtConfig.JwtAccessExp, cfg.JwtConfig.JwtRefreshExp, cfg.JwtConfig.JwtIssuer)
	// load module
	userModule := module.NewUserModule(db, rabbitMQService)
	authModule := module.NewAuthModule(db, jwtService)

	// gRPC setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(cfg.GrpcPort)))
	if err != nil {
		log.Fatalf("❌ Lỗi lắng nghe port 50001: %v", err)
	}
	// gRPC service setup
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userModule.Service)
	auth.RegisterAuthServiceServer(grpcServer, authModule.AuthService)

	

	reflection.Register(grpcServer)

	log.Println("🚀 User gRPC Service đang chạy tại port :50001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Lỗi khởi chạy gRPC Server: %v", err)
	}
}