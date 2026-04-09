package product_config

import util "go-ecommerce/common/utils"

type DatabaseConfig struct {
	MongoUri    string
	MongoDBName string
}

type RedisConfig struct {
	RedisUri string
}

type RabbitMQConfig struct {
	Uri string
}

type ProductServiceConfig struct {
	DatabaseConfig *DatabaseConfig
	GrpcPort int
	RedisConfig RedisConfig
	RabbitMQConfig *RabbitMQConfig
}

func NewProductServiceConfig() *ProductServiceConfig {
	// mongo
	mongoUri := util.GetEnv("MONGO_URI", "mongodb://localhost:27020")
	mongoDBName := util.GetEnv("DB_NAME_PRODUCT_SERVICE", "go_ecommerce_product_service")
	
	// grpc
	grpcPort := util.GetIntEnv("GRPC_PRODUCT_SERVICE_PORT", 50002)

	// redis
	redisUri := util.GetEnv("REDIS_URI", "localhost:6379")

	// rabbitmq
	rabbitUri := util.GetEnv("RABBIT_MQ_URI", "amqp://root:admin123@localhost:5672/")
	return &ProductServiceConfig{
		DatabaseConfig: &DatabaseConfig{
			MongoUri: mongoUri,
			MongoDBName: mongoDBName,
		},
		GrpcPort: grpcPort,
		RedisConfig: RedisConfig{
			RedisUri: redisUri,
		},
		RabbitMQConfig: &RabbitMQConfig{
			Uri: rabbitUri,
		},
	}
}