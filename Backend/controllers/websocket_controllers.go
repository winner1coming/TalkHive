package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// HandleConnections 处理WebSocket连接
func HandleConnections(c *gin.Context) {
	userID := c.Param("account_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var me models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if me.Deactivate {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// 获取 WebSocket 连接
	conn, err := global.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// 将连接存储在全局map中
	global.Clients[uint(accountID)] = conn
	defer delete(global.Clients, uint(accountID))

	// 监听消息
	for {
		var msg global.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}
		// 处理接收到的消息
		HandleMessages(msg)
	}
}

// HandleMessages 处理消息逻辑
func HandleMessages(msg global.Message) {
	if msg.IsGroup {
		var groupMembers []models.GroupMemberInfo
		if err := global.Db.Where("group_id = ?", msg.AccountID).First(&groupMembers).Error; err != nil {
			log.Println("查询群成员失败:", err)
		}
		for _, member := range groupMembers {
			if targetConn, ok := global.Clients[member.AccountID]; ok {
				if err := targetConn.WriteJSON(msg); err != nil {
					log.Println("群聊消息发送失败:", err)
				}
			} else {
				log.Println("目标用户未连接")
			}
		}
	} else {
		if targetConn, ok := global.Clients[msg.TargetID]; ok {
			var sender models.AccountInfo
			if err := global.Db.Where("account_id = ?", msg.AccountID).First(&sender).Error; err != nil {
				log.Println("查询发送者信息失败:", err)
			}
			var targetUser models.AccountInfo
			if err := global.Db.Where("account_id = ?", msg.TargetID).First(&targetUser).Error; err != nil {
				log.Println("查询目标用户信息失败:", err)
			}
			var messageRecord models.MessageInfo
			if err := global.Db.Where("send_account_id = ? AND target_id = ?", msg.AccountID, msg.TargetID).First(&messageRecord).Error; err != nil {
				log.Println("查询消息记录失败:", err)
			}

			if err := targetConn.WriteJSON(msg); err != nil {
				log.Println("单聊消息发送失败:", err)
			}

		} else {
			log.Println("目标用户未连接")
		}
	}
}
