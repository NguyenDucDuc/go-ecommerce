package config

import util "go-ecommerce/common/utils"

type DatabaseConfig struct {
	MongoUri    string
	MongoDBName string
}

type ProductServiceConfig struct {
	DatabaseConfig *DatabaseConfig
	GrpcPort int
}

func NewProductServiceConfig() *ProductServiceConfig {
	// mongo
	mongoUri := util.GetEnv("MONGO_URI", "mongodb://localhost:27020")
	mongoDBName := util.GetEnv("DB_NAME_PRODUCT_SERVICE", "go_ecommerce_product_service")
	
	// grpc
	grpcPort := util.GetIntEnv("GRPC_PRODUCT_SERVICE_PORT", 50002)
	return &ProductServiceConfig{
		DatabaseConfig: &DatabaseConfig{
			MongoUri: mongoUri,
			MongoDBName: mongoDBName,
		},
		GrpcPort: grpcPort,
	}
}