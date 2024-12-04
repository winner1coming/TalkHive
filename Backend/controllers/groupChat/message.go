package groupChat

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CopyMessage 复制消息
func CopyMessage(c *gin.Context) {
	var request struct {
		MessageID uint `json:"message_id"`
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 查找需要复制的消息
	var message models.MessageInfo
	if err := db.Where("message_id = ?", request.MessageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// 创建一个消息，将原消息内容复制过来
	newMessage := models.MessageInfo{
		SendAccountID: message.SendAccountID,
		Content:       message.Content,
		Type:          message.Type,
		CreateTime:    time.Now(),
	}

	// 保存新消息
	if err := db.Create(&newMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message copied"})
}

// FavoriteMessage 收藏消息
func FavoriteMessage(c *gin.Context) {
	var request struct {
		MessageID uint `json:"message_id"`
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 查找需要收藏的消息
	var message models.MessageInfo
	if err := db.Where("message_id = ?", request.MessageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// 检查该消息是否已经被收藏
	var favorite models.Favorites
	if err := db.Where("table_name = ? AND id = ?", "message_info", request.MessageID).First(&favorite).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message already favorited"})
		return
	}

	// 将消息添加到收藏
	favorite = models.Favorites{
		TableName: "message_info",
		ID:        request.MessageID,
		AccountID: message.SendAccountID, // 假设使用消息发送者的 AccountID
	}

	if err := db.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to favorite message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message favorited"})
}

// ReplyMessage 回复消息
func ReplyMessage(c *gin.Context) {
	var request struct {
		MessageID uint   `json:"message_id"`
		Content   string `json:"content"`
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 查找原始消息
	var message models.MessageInfo
	if err := db.Where("message_id = ?", request.MessageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// 创建新消息作为回复
	replyMessage := models.MessageInfo{
		SendAccountID: message.SendAccountID, // 回复的账户ID
		Content:       request.Content,       // 回复的内容
		Type:          "reply",
		CreateTime:    time.Now(),
	}

	// 保存回复消息
	if err := db.Create(&replyMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reply message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message replied"})
}

// ForwardMessage 转发消息
func ForwardMessage(c *gin.Context) {
	var request struct {
		MessageID uint `json:"message_id"`
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 查找需要转发的消息
	var message models.MessageInfo
	if err := db.Where("message_id = ?", request.MessageID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// 创建转发的新消息
	forwardedMessage := models.MessageInfo{
		SendAccountID: message.SendAccountID, // 转发者的账户ID
		Content:       message.Content,
		Type:          "forward",
		CreateTime:    time.Now(),
	}

	// 保存转发消息
	if err := db.Create(&forwardedMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message forwarded"})
}

// DeleteMessage 删除消息
func DeleteMessage(c *gin.Context) {
	var request struct {
		MessageID uint `json:"message_id"`
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 查找并删除消息
	if err := db.Where("message_id = ?", request.MessageID).Delete(&models.MessageInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}

// BatchMessageOperation 多选消息操作（收藏、转发、删除）
func BatchMessageOperation(c *gin.Context) {
	var request struct {
		MessageIDs []uint `json:"message_ids"`
		Operation  string `json:"operation"` // "favorite", "forward", or "delete"
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	switch request.Operation {
	case "favorite":
		for _, messageID := range request.MessageIDs {
			var message models.MessageInfo
			if err := db.Where("message_id = ?", messageID).First(&message).Error; err != nil {
				continue // Skip if message not found
			}

			// 执行收藏操作
			favorite := models.Favorites{
				TableName: "message_info",
				ID:        messageID,
				AccountID: message.SendAccountID, // 假设使用发送者的 AccountID
			}

			if err := db.Create(&favorite).Error; err != nil {
				continue
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Messages favorited"})
	case "forward":
		for _, messageID := range request.MessageIDs {
			var message models.MessageInfo
			if err := db.Where("message_id = ?", messageID).First(&message).Error; err != nil {
				continue // Skip if message not found
			}

			// 执行转发操作
			forwardedMessage := models.MessageInfo{
				SendAccountID: message.SendAccountID, // 转发者的账户ID
				Content:       message.Content,
				Type:          "forward",
				CreateTime:    time.Now(),
			}

			if err := db.Create(&forwardedMessage).Error; err != nil {
				continue
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Messages forwarded"})
	case "delete":
		for _, messageID := range request.MessageIDs {
			if err := db.Where("message_id = ?", messageID).Delete(&models.MessageInfo{}).Error; err != nil {
				continue // Skip if deletion fails
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Messages deleted"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation"})
	}
}

// AtAllMembers 管理员 @所有人
func AtAllMembers(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// 验证当前用户是否是管理员
	accountID := c.MustGet("account_id").(uint) // 假设中间件传递了 account_id
	var groupMember models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ? AND is_banned = false", groupID, accountID).
		First(&groupMember).Error; err != nil || groupMember.GroupNickname != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only administrators can @all members"})
		return
	}

	// 创建 @所有人 的消息
	content := "@所有人：请注意重要信息！"
	message := models.MessageInfo{
		CreateTime:    time.Now(),
		SendAccountID: accountID,
		Content:       content,
		Type:          "text", // 假设消息类型为文本
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to all members"})
}

// AtSingleMember 普通成员 @单人
func AtSingleMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		TargetAccountID uint   `json:"target_account_id"` // 被@的用户ID
		Content         string `json:"content"`           // 消息内容
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 验证当前用户是否在群聊中
	accountID := c.MustGet("account_id").(uint)
	var groupMember models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ? AND is_banned = false", groupID, accountID).
		First(&groupMember).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not in this group"})
		return
	}

	// 验证目标用户是否在群聊中
	var targetMember models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ?", groupID, request.TargetAccountID).
		First(&targetMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found in this group"})
		return
	}

	// 创建消息
	content := "@" + targetMember.GroupNickname + ": " + request.Content
	message := models.MessageInfo{
		CreateTime:    time.Now(),
		SendAccountID: accountID,
		Content:       content,
		Type:          "text", // 假设消息类型为文本
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent to target member"})
}
