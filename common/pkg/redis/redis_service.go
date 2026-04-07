package pkg_redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) IRedisService {
	return &RedisService{
		client: client,
	}
}

// SetJSON: Chuyển bất kỳ struct nào thành JSON string để lưu vào Redis
func (r *RedisService) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

// GetJSON: Lấy JSON string và nạp ngược lại vào struct (dest phải là con trỏ)
func (r *RedisService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err // Có thể là redis.Nil nếu không tìm thấy
	}
	return json.Unmarshal(data, dest)
}

func (r *RedisService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}