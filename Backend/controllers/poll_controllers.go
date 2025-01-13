package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Poll(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		LastAccessTime string `json:"lastAccessTime" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	var has_new_message bool
	var has_new_friendrequest bool
	var has_new_grouprequest bool

	// 查询ChatInfo表
	var chatInfos []models.ChatInfo
	if err := global.Db.Where("account_id = ?", accountID).Find(&chatInfos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询聊天信息失败"})
		return
	}
	for _, chatInfo := range chatInfos {
		var messages []models.MessageInfo
		if err := global.Db.Where("(sender_chat_id = ? OR receiver_chat_id = ?)AND create_time > ?", chatInfo.ChatID, chatInfo.ChatID, input.LastAccessTime).Find(&messages).Error; err == nil {
			has_new_message = true
			break
		}
	}

	// 查询FriendRequest
	var friendRequests []models.ApplyInfo
	if err := global.Db.Where("(sender_id = ? OR receiver_id = ? ) AND apply_type = ?", accountID, accountID, "friend").Find(&friendRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群组申请失败"})
		return
	}
	for _, friendRequest := range friendRequests {
		if friendRequest.SenderID == user.AccountID { // 判断对方同意时间
			if friendRequest.DealTime > input.LastAccessTime {
				has_new_friendrequest = true
			}
		} else if friendRequest.ReceiverID == user.AccountID { // 判断对方发送时间
			if friendRequest.SendTime > input.LastAccessTime {
				has_new_friendrequest = true
			}
		}
	}

	// 查询GroupRequest
	var groupRequests []models.ApplyInfo
	if err := global.Db.Where("(sender_id = ? OR receiver_id = ? ) AND (apply_type = ? OR apply_type = ?)", accountID, accountID, "groupInvitation", "groupApply").Find(&groupRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群组申请失败"})
		return
	}
	for _, groupRequest := range groupRequests {
		if groupRequest.ApplyType == "groupInvitation" { // 邀请加入群聊
			if groupRequest.SenderID == user.AccountID { // 判断对方同意时间
				if groupRequest.DealTime > input.LastAccessTime {
					has_new_grouprequest = true
				}
			} else if groupRequest.ReceiverID == user.AccountID { // 判断对方发送时间
				if groupRequest.SendTime > input.LastAccessTime {
					has_new_grouprequest = true
				}
			}
		} else { // 申请加入群聊
			if groupRequest.SenderID == user.AccountID { // 判断对方同意时间
				if groupRequest.DealTime > input.LastAccessTime {
					has_new_grouprequest = true
				}
			} else if groupRequest.ReceiverID == user.AccountID { // 判断对方发送时间
				if groupRequest.SendTime > input.LastAccessTime {
					has_new_grouprequest = true
				}
			}
		}
	}

	response := gin.H{
		"success": true,
		"message": "获取消息成功",
		"data": gin.H{
			"has_new_message":       has_new_message,
			"has_new_friendrequest": has_new_friendrequest,
			"has_new_grouprequest":  has_new_grouprequest,
		},
	}

	c.JSON(http.StatusOK, response)
}
