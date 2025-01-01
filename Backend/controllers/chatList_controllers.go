package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetChatList 获取聊天列表
func GetChatList(c *gin.Context) {
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
	if user.Deactivate == true {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
		return
	}

	// 获取用户的聊天记录
	var chatList []models.ChatInfo
	if err := global.Db.Where("account_id = ?", accountID).Find(&chatList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询聊天记录失败"})
		return
	}

	var response []gin.H
	for _, chat := range chatList {
		var friend models.AccountInfo
		var groupChat models.GroupChatInfo

		// 判断好友或群聊是否存在
		if chat.IsGroup {
			if err := global.Db.Where("group_id = ?", chat.TargetID).First(&groupChat).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
				return
			}
		} else {
			if err := global.Db.Where("account_id = ?", chat.TargetID).First(&friend).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "好友不存在"})
				return
			}
		}

		// 获取最后一条消息
		var lastMessage models.MessageInfo
		if err := global.Db.Where("chat_id = ?", chat.ChatID).Order("timestamp desc").First(&lastMessage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询最后一条消息失败"})
			return
		}

		// 获取未读消息数
		var unreadCount int64
		if err := global.Db.Model(&models.MessageInfo{}).Where("chat_id = ? AND receive_account_id = ? AND is_read = ?", chat.ChatID, accountID, false).Count(&unreadCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询未读消息失败"})
			return
		}

		// 查询contacts表
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, chat.TargetID).First(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询contacts表失败"})
			return
		}
		if chat.IsGroup != true {
			chatResponse := gin.H{
				"id":              chat.TargetID,
				"avatar":          friend.Avatar,
				"name":            friend.Nickname,
				"remark":          contact.Remark,
				"lastMessage":     lastMessage.Content,
				"lastMessageTime": lastMessage.CreateTime,
				"unreadCount":     unreadCount,
			}
			response = append(response, chatResponse)
		} else {
			chatResponse := gin.H{
				"id":              chat.TargetID,
				"avatar":          groupChat.GroupAvatar,
				"name":            groupChat.GroupName,
				"remark":          contact.Remark,
				"lastMessage":     lastMessage.Content,
				"lastMessageTime": lastMessage.CreateTime,
				"unreadCount":     unreadCount,
			}
			response = append(response, chatResponse)
		}

	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// GetChat 新建单个聊天
func GetChat(c *gin.Context) {

}

// SearchChats 搜索聊天
func SearchChats(c *gin.Context) {
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
	if user.Deactivate == true {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
		return
	}
	keyword := c.Param("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "搜索关键字不能为空"})
		return
	}

	// 查询ChatInfo表
	var chatList []models.ChatInfo
	if err := global.Db.Where("target_id LIKE ? OR chat_id LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&chatList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询聊天记录失败"})
		return
	}

	var response []gin.H
	for _, chat := range chatList {
		var friend models.AccountInfo
		var groupChat models.GroupChatInfo

		// 判断是好友还是群聊
		if chat.IsGroup {
			if err := global.Db.Where("group_id = ?", chat.TargetID).First(&groupChat).Error; err != nil {
				continue // 群聊信息不存在则跳过
			}
		} else {
			if err := global.Db.Where("account_id = ?", chat.TargetID).First(&friend).Error; err != nil {
				continue // 好友信息不存在则跳过
			}
		}

		// 获取最后一条消息
		var lastMessage models.MessageInfo
		if err := global.Db.Where("chat_id = ?", chat.ChatID).Order("timestamp desc").First(&lastMessage).Error; err != nil {
			continue // 查询最后一条消息失败则跳过
		}

		// 获取未读消息数
		var unreadCount int64
		if err := global.Db.Model(&models.MessageInfo{}).Where("chat_id = ? AND receive_account_id = ? AND is_read = ?", chat.ChatID, c.GetHeader("User-ID"), false).Count(&unreadCount).Error; err != nil {
			continue // 查询未读消息数失败则跳过
		}

		// 查询联系人信息
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", c.GetHeader("User-ID"), chat.TargetID).First(&contact).Error; err != nil {
			continue // 查询联系人信息失败则跳过
		}

		// 构造聊天记录响应
		var chatResponse gin.H
		if chat.IsGroup {
			chatResponse = gin.H{
				"id":              chat.TargetID,
				"avatar":          groupChat.GroupAvatar,
				"name":            groupChat.GroupName,
				"remark":          contact.Remark,
				"lastMessage":     lastMessage.Content,
				"lastMessageTime": lastMessage.CreateTime,
				"unreadCount":     unreadCount,
				"tags":            []string{"unread"}, // 示例标签，可以根据实际情况动态调整
			}
		} else {
			chatResponse = gin.H{
				"id":              chat.TargetID,
				"avatar":          friend.Avatar,
				"name":            friend.Nickname,
				"remark":          contact.Remark,
				"lastMessage":     lastMessage.Content,
				"lastMessageTime": lastMessage.CreateTime,
				"unreadCount":     unreadCount,
				"tags":            []string{"unread"}, // 示例标签，可以根据实际情况动态调整
			}
		}

		response = append(response, chatResponse)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

// PinChat 置顶或取消置顶聊天
func PinChat(c *gin.Context) {
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
	var me models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if me.Deactivate == true {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
		return
	}

	var input struct {
		Tid      string `json:"tid"`
		IsPinned string `json:"is_pinned"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json绑定失败"})
		return
	}
	var other models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.Tid).First(&other).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if other.Deactivate == true {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
		return
	}

	// 查找聊天记录
	var chat models.ChatInfo
	if err := global.Db.Where("chat_id = ?", chatID).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "聊天记录不存在"})
		return
	}

	// 更新聊天的置顶状态
	var newPinned bool
	if pinned == "true" {
		newPinned = true
	} else {
		newPinned = false
	}

	if err := global.Db.Model(&chat).Update("is_pinned", newPinned).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新置顶状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "置顶状态更新成功"})
}

// ReadMessages 标记为已读或未读消息
func ReadMessages(c *gin.Context) {

}
