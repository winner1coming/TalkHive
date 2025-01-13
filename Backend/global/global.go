package global

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Db         *gorm.DB
	RedisDB    *redis.Client
	SmsCodeKey = "sms_code:" // 短信验证码的 Redis 键前缀
)

// ParseUint - 将字符串转换为 uint 类型
func ParseUint(input string) uint {
	value, _ := strconv.ParseUint(input, 10, 32) // 忽略错误，出错返回 0
	return uint(value)
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
