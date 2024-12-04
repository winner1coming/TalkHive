package singleChat

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// SendFriendCard 发送好友名片
func SendFriendCard(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)

	// 获取好友信息
	var friend models.AccountInfo
	if err := db.Where("account_id = ?", friendID).First(&friend).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	// 发送名片消息
	message := models.MessageInfo{
		CreateTime:    time.Now(),
		SendAccountID: accountID,
		Content:       friend.Nickname + " - " + friend.ID, // 名片内容可以根据实际需求定制
		Type:          "card",
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card sent successfully"})
}

// DeleteChatRecords 删除聊天记录
func DeleteChatRecords(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)

	// 删除指定好友的聊天记录
	if err := db.Where("send_account_id = ? AND receive_account_id = ?", accountID, friendID).
		Delete(&models.MessageInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat records deleted"})
}
