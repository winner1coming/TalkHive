package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

// RedisClient 全局 Redis 客户端
var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() {
	// 创建 Redis 客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		DB:       AppConfig.Redis.DB,
		Password: AppConfig.Redis.Password,
	})

	// 测试 Redis 连接
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}

	log.Println("连接Redis成功")
}

// GetRedis 返回全局数据库实例
func GetRedis() *redis.Client {
	if RedisClient == nil {
		log.Fatal("Redis实例为空")
	}
	return RedisClient
}
