package config

import (
	"TalkHive/global"
	"github.com/go-redis/redis"
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
	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis, got error: %v", err)
	}

	global.RedisDB = RedisClient
	log.Println("连接Redis成功")
}

// GetRedis 返回全局数据库实例
func GetRedis() *redis.Client {
	if RedisClient == nil {
		log.Fatal("Redis实例为空")
	}
	return RedisClient
}
