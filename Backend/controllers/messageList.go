package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddFriend 添加好友
func AddFriend(c *gin.Context) {
	var request struct {
		UserID   uint   `json:"user_id"`
		FriendID uint   `json:"friend_id"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact := models.Contacts{
		ContactID:   request.FriendID,
		Remark:      request.Remark,
		IsGroupChat: false,
	}

	if err := config.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	var group models.GroupChatInfo
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group created successfully"})
}

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

// SubmitApply 提交申请
func SubmitApply(c *gin.Context) {
	var applyInfo models.ApplyInfo
	if err := c.ShouldBindJSON(&applyInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&applyInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
}

// ReviewApply 审核申请
func ReviewApply(c *gin.Context) {
	var review struct {
		ApplyID uint   `json:"apply_id"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.ApplyInfo{}).Where("apply_id = ?", review.ApplyID).Update("status", review.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application reviewed successfully"})
}
