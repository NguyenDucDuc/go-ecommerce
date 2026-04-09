package order_config

import util "go-ecommerce/common/utils"

type DatabaseConfig struct {
	MongoUri    string
	MongoDBName string
}

type RedisConfig struct {
	RedisUri string
}

type GrpcConfig struct {
	GrpcPort int
	GrpcProductUri string
}

type ProductServiceConfig struct {
	DatabaseConfig *DatabaseConfig
	RedisConfig RedisConfig
	GrpcConfig GrpcConfig
}


func NewProductServiceConfig() *ProductServiceConfig {
	// mongo
	mongoUri := util.GetEnv("MONGO_URI", "mongodb://localhost:27020")
	mongoDBName := util.GetEnv("DB_NAME_ORDER_SERVICE", "go_ecommerce_order_service")
	
	// grpc
	grpcPort := util.GetIntEnv("GRPC_ORDER_SERVICE_PORT", 50003)
	grpcProductUri := util.GetEnv("GRPC_PRODUCT_SERVICE_URI", "localhost:50002")

	// redis
	redisUri := util.GetEnv("REDIS_URI", "localhost:6379")
	return &ProductServiceConfig{
		DatabaseConfig: &DatabaseConfig{
			MongoUri: mongoUri,
			MongoDBName: mongoDBName,
		},
		GrpcConfig: GrpcConfig{
			GrpcPort: grpcPort,
			GrpcProductUri: grpcProductUri,
		},
		RedisConfig: RedisConfig{
			RedisUri: redisUri,
		},
	}
}