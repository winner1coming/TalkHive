package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Db         *gorm.DB
	RedisDB    *redis.Client
	SmsCodeKey = "sms_code:" // 短信验证码的 Redis 键前缀
)
