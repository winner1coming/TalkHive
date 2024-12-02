package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	var message struct {
		SenderID uint `json:"sender_id"`
		//ReceiverID uint   `json:"receiver_id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := models.MessageInfo{
		SendAccountID: message.SenderID,
		// ReceiverID:    message.ReceiverID,
		Content: message.Content,
	}

	if err := config.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

// GetMessages 获取消息记录
func GetMessages(c *gin.Context) {
	receiverID := c.Param("id")

	var messages []models.MessageInfo
	if err := config.DB.Where("receiver_id = ?", receiverID).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
