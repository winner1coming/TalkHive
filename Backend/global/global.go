package global

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Db         *gorm.DB
	RedisDB    *redis.Client
	SmsCodeKey = "sms_code:"                    // 短信验证码的 Redis 键前缀
	Clients    = make(map[uint]*websocket.Conn) // 存储WebSocket连接，每个用户一个连接
	Broadcast  = make(chan Message)             // Broadcast 广播消息的通道
	Mu         sync.Mutex                       // 互斥锁
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

// Message 一条消息的类型
type Message struct {
	MessageID  uint   `json:"message_id"`
	AccountID  uint   `json:"sender_id"`
	TargetID   uint   `json:"target_id"`
	Avatar     string `json:"avatar"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	IsGroup    bool   `json:"is_group"`
	CreateTime string `json:"create_time"`
}
