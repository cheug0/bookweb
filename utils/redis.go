package utils

import (
	"bookweb/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	redisCtx    = context.Background()
)

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.RedisConfig) error {
	if cfg == nil || !cfg.Enabled {
		RedisClient = nil
		return nil
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	_, err := RedisClient.Ping(redisCtx).Result()
	if err != nil {
		RedisClient = nil
		return err
	}
	return nil
}

// CacheGet 从缓存获取数据
func CacheGet(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("redis not enabled")
	}
	return RedisClient.Get(redisCtx, key).Result()
}

// CacheSet 设置缓存数据
func CacheSet(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("redis not enabled")
	}
	return RedisClient.Set(redisCtx, key, value, expiration).Err()
}

// CacheDel 删除缓存
func CacheDel(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("redis not enabled")
	}
	return RedisClient.Del(redisCtx, key).Err()
}

// IsRedisEnabled 检查 Redis 是否已启用
func IsRedisEnabled() bool {
	return RedisClient != nil
}
