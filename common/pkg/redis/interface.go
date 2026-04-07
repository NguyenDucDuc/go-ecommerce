package pkg_redis

import (
	"context"
	"time"
)

type IRedisService interface {
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
}