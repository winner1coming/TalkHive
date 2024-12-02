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
