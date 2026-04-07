package pkg_product_redis

import (
	"context"
	"fmt"
	product_config "go-ecommerce/product-service/internal/config"
	"log"

	"github.com/redis/go-redis/v9"
)
var ctx = context.Background()
func ConnectRedis(cfg product_config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUri, // Địa chỉ Redis server
		Password: "",               // Mật khẩu (để trống nếu không có)
		DB:       0,                // ID database (mặc định là 0)
		PoolSize: 10,               // Số lượng connection tối đa trong pool
	})

	// 2. Kiểm tra kết nối bằng Ping
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}
	fmt.Println("✅ Kết nối Redis thành công!")

	return rdb
}