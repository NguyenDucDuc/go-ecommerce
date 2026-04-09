package config

import util "go-ecommerce/common/utils"


type JwtConfig struct {
	JwtSecret     string
	JwtAccessExp  int
	JwtRefreshExp int
	JwtIssuer     string
}
type Config struct {
	Port int
	GrpcUserServiceUri string
	GrpcProductServiceUri string
	GrpcOrderServiceUri string
	JwtConfig *JwtConfig
}



func NewConfig() *Config {
	jwtSecret := util.GetEnv("JWT_SECRET", "01234567890123456789012345678912")
	jwtIssuer := util.GetEnv("JWT_ISSUER", "go_ecommerce")
	jwtAccessExp := util.GetIntEnv("JWT_ACCESS_EXP", 5)
	jwtRefreshExp := util.GetIntEnv("JWT_REFRESH_EXP", 120)
	return &Config{
		Port: util.GetIntEnv("PORT", 4000),
		GrpcUserServiceUri: util.GetEnv("GRPC_USER_SERVICE_URI", "localhost:50001"),
		GrpcProductServiceUri: util.GetEnv("GRPC_PRODUCT_SERVICE_URI", "localhost:50002"),
		GrpcOrderServiceUri: util.GetEnv("GRPC_ORDER_SERVICE_URI", "localhost:50003"),
		JwtConfig: &JwtConfig{
			JwtSecret: jwtSecret,
			JwtIssuer: jwtIssuer,
			JwtAccessExp: jwtAccessExp,
			JwtRefreshExp: jwtRefreshExp,
		},
	}
}