package config

import util "go-ecommerce/common/utils"

type DatabaseConfig struct {
	MongoUri    string
	MongoDBName string
}

type JwtConfig struct {
	JwtSecret     string
	JwtAccessExp  int
	JwtRefreshExp int
	JwtIssuer     string
}

type UserServiceConfig struct {
	DatabaseConfig *DatabaseConfig
	JwtConfig      *JwtConfig
	GrpcPort int
}

func NewUserServiceConfig() *UserServiceConfig {
	// mongo
	mongoUri := util.GetEnv("MONGO_URI", "mongodb://localhost:27020")
	mongoDBName := util.GetEnv("DB_NAME_USER_SERVICE", "go_ecommerce_user_service")
	// jwt
	jwtSecret := util.GetEnv("JWT_SECRET", "01234567890123456789012345678912")
	jwtIssuer := util.GetEnv("JWT_ISSUER", "go_ecommerce")
	jwtAccessExp := util.GetIntEnv("JWT_ACCESS_EXP", 5)
	jwtRefreshExp := util.GetIntEnv("JWT_REFRESH_EXP", 120)
	// grpc
	grpcPort := util.GetIntEnv("GRPC_USER_SERVICE_PORT", 50001)
	return &UserServiceConfig{
		DatabaseConfig: &DatabaseConfig{
			MongoUri: mongoUri,
			MongoDBName: mongoDBName,
		},
		JwtConfig: &JwtConfig{
			JwtSecret: jwtSecret,
			JwtIssuer: jwtIssuer,
			JwtAccessExp: jwtAccessExp,
			JwtRefreshExp: jwtRefreshExp,
		},
		GrpcPort: grpcPort,
	}
}