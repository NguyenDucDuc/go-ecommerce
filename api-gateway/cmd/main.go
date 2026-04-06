package main

import (
	"fmt"
	"go-ecommerce/api-gateway/internal/config"
	"go-ecommerce/api-gateway/internal/module"
	"go-ecommerce/common/gen-proto/auth"
	product "go-ecommerce/common/gen-proto/products"
	user "go-ecommerce/common/gen-proto/users"
	util "go-ecommerce/common/utils"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	util.LoadEnv()
	cfg := config.NewConfig()

	// setup grpc service
	userConn, err := grpc.NewClient(cfg.GrpcUserServiceUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Không thể kết nối tới User Service: %v", err)
	}
	defer userConn.Close()

	productConn, err := grpc.NewClient(cfg.GrpcProductServiceUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Không thể kết nối tới Product Service: %v", err)
	}
	defer productConn.Close()

	r := gin.Default()
	v1 := r.Group("/api/v1")

	// load module
	userClient := user.NewUserServiceClient(userConn)
	authClient := auth.NewAuthServiceClient(userConn)
	productClient := product.NewProductServiceClient(productConn)

	userModule := module.NewUserModule(userClient)
	userModule.Routes.RegisterRoutes(v1)

	authModule := module.NewAuthModule(authClient)
	authModule.Routes.RegisterRoutes(v1)

	productModule := module.NewProductModule(productClient)
	productModule.Routes.RegisterRoutes(v1)

	r.Run(fmt.Sprintf(":%s", strconv.Itoa(cfg.Port)))
}

