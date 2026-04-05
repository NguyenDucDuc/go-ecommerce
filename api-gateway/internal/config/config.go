package config

import util "go-ecommerce/common/utils"

type Config struct {
	Port int
	GrpcUserServiceUri string
	GrpcProductServiceUri string
	GrpcOrderServiceUri string
}

func NewConfig() *Config {
	return &Config{
		Port: util.GetIntEnv("PORT", 4000),
		GrpcUserServiceUri: util.GetEnv("GRPC_USER_SERVICE_URI", "localhost:50001"),
		GrpcProductServiceUri: util.GetEnv("GRPC_PRODUCT_SERVICE_URI", "localhost:50002"),
		GrpcOrderServiceUri: util.GetEnv("GRPC_ORDER_SERVICE_URI", "localhost:50003"),
	}
}