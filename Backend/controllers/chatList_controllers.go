package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"TalkHive/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// ---------------------------------------------------
// 聊天列表

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
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "当前用户无聊天记录"})
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

		// 查询contacts表
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, chat.TargetID).First(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询contacts表失败"})
			return
		}

		// 如果没有消息
		var lastMessage models.MessageInfo
		if err := global.Db.Where("sender_chat_id = ? OR receiver_chat_id = ?", chat.ChatID, chat.ChatID).Order("create_time desc").First(&lastMessage).Error; err != nil {
			var tags []string
			if contact.IsGroupChat {
				tags = append(tags, "group")
			} else {
				tags = append(tags, "friend")
			}

			if contact.IsPinned {
				tags = append(tags, "Pinned")
			}
			if contact.IsBlocked {
				tags = append(tags, "blocked")
			}

			if !chat.IsGroup {
				avatarBase64, mimeType, err := utils.GetFileContentAndType(friend.Avatar)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
					return
				}
				avatarBase64 = "data:" + mimeType + ";base64," + avatarBase64
				chatResponse := gin.H{
					"id":              chat.TargetID,
					"avatar":          avatarBase64,
					"mimeType":        mimeType,
					"name":            friend.Nickname,
					"remark":          contact.Remark,
					"lastMessage":     nil,
					"lastMessageTime": nil,
					"unreadCount":     nil,
					"tags":            tags,
				}
				response = append(response, chatResponse)
				continue
			} else {
				avatarBase64, mimeType, err := utils.GetFileContentAndType(friend.Avatar)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
					return
				}
				avatarBase64 = "data:" + mimeType + ";base64," + avatarBase64

				chatResponse := gin.H{
					"id":              chat.TargetID,
					"avatar":          avatarBase64,
					"mimeType":        mimeType,
					"name":            groupChat.GroupName,
					"remark":          contact.Remark,
					"lastMessage":     nil,
					"lastMessageTime": nil,
					"unreadCount":     nil,
					"tags":            tags,
				}
				response = append(response, chatResponse)
				continue
			}
		}

		// 获取未读消息数/或者查询Contacts表中unread_message_num
		var unreadCount int64
		if err := global.Db.Model(&models.MessageInfo{}).Where("chat_id = ? AND target_id = ? AND is_read = ?", chat.ChatID, accountID, false).Count(&unreadCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询未读消息失败"})
			return
		}

		var tags []string
		if chat.IsGroup != true {
			tags := append(tags, "friend")
			if contact.IsPinned {
				tags = append(tags, "Pinned")
			}
			if contact.IsBlocked {
				tags = append(tags, "blocked")
			}

			avatarBase64, mimeType, err := utils.GetFileContentAndType(friend.Avatar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
				return
			}
			avatarBase64 = "data:" + mimeType + ";base64," + avatarBase64

			chatResponse := gin.H{
				"id":              chat.TargetID,
				"avatar":          avatarBase64,
				"name":            friend.Nickname,
				"remark":          contact.Remark,
				"lastMessage":     lastMessage.Content,
				"lastMessageTime": lastMessage.CreateTime,
				"unreadCount":     unreadCount,
				"tags":            tags,
			}
			response = append(response, chatResponse)
		} else {
			tags := append(tags, "group")
			if contact.IsPinned {
				tags = append(tags, "Pinned")
			}
			if contact.IsBlocked {
				tags = append(tags, "blocked")
			}
			avatarBase64, mimeType, err := utils.GetFileContentAndType(friend.Avatar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
				return
			}
			avatarBase64 = "data:" + mimeType + ";base64," + avatarBase64
			chatResponse := gin.H{
				"id":              chat.TargetID,
				"avatar":          avatarBase64,
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

// CreateChat 创建单个群聊
func CreateChat(c *gin.Context) {
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
	if me.Deactivate {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
		return
	}
	var input struct {
		Tid     int  `json:"tid"`
		IsGroup bool `json:"is_group"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json绑定失败"})
		return
	}

	var response []gin.H
	var tags []string
	if input.IsGroup {
		var group models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", input.Tid).First(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊失败"})
			return
		}
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ? AND is_group_chat = ?", accountID, input.Tid, true).First(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询contacts表失败"})
			return
		}
		contact.UnreadMessageNum = 0
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新contacts表失败"})
			return
		}

		// 创建me对群聊的聊天记录
		chatMe := models.ChatInfo{
			AccountID:  uint(accountID),
			TargetID:   uint(input.Tid),
			IsGroup:    input.IsGroup,
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := global.Db.Create(&chatMe).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建聊天记录失败"})
			return
		}

		tags := append(tags, "group")
		if contact.IsPinned {
			tags = append(tags, "pinned")
		}
		if contact.IsBlocked {
			tags = append(tags, "blocked")
		}

		avatarBase64, mimeType, err := utils.GetFileContentAndType(group.GroupAvatar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		avatarBase64 = "data:" + mimeType + ";base64," + avatarBase64
		chatResponse := gin.H{
			"id":              chatMe.TargetID,
			"avatar":          avatarBase64,
			"name":            group.GroupName,
			"remark":          contact.Remark,
			"lastMessage":     nil,
			"lastMessageTime": nil,
			"unreadCount":     nil,
			"tags":            tags,
		}
		response = append(response, chatResponse)

	} else {
		var friend models.AccountInfo
		if err := global.Db.Where("account_id = ?", input.Tid).First(&friend).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询对方用户失败"})
			return
		}
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, input.Tid).First(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询联系人失败"})
			return
		}

		contact.UnreadMessageNum = 0
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新contacts表失败"})
			return
		}

		// 生成me对other以及other对me的聊天消息记录
		chat := models.ChatInfo{
			AccountID:  uint(accountID),
			TargetID:   uint(input.Tid),
			IsGroup:    input.IsGroup,
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := global.Db.Create(&chat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建聊天记录失败"})
			return
		}

		chatOther := models.ChatInfo{
			AccountID:  uint(input.Tid),
			TargetID:   uint(accountID),
			IsGroup:    input.IsGroup,
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := global.Db.Create(&chatOther).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建聊天记录失败"})
			return
		}

		tags := append(tags, "friend")
		if contact.IsPinned {
			tags = append(tags, "Pinned")
		}
		if contact.IsBlocked {
			tags = append(tags, "blocked")
		}

		avatarBase64, mimeType, err := utils.GetFileContentAndType(friend.Avatar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		avatarBase64 = "data:" + mimeType + ";base64," + avatarBase64

		chatResponse := gin.H{
			"id":              chat.TargetID,
			"avatar":          avatarBase64,
			"name":            friend.Nickname,
			"remark":          contact.Remark,
			"lastMessage":     nil,
			"lastMessageTime": nil,
			"unreadCount":     nil,
			"tags":            tags,
		}
		response = append(response, chatResponse)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建聊天记录成功", "data": response})
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "搜索关键字为空"})
		return
	}

	// 查询AccountInfo表，查询昵称nickname
	var accounts []models.AccountInfo
	if err := global.Db.Where("nickname LIKE ? OR account_id = ?", "%"+keyword+"%", keyword).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}

	// 查询Contacts表，查询备注remark
	var contacts []models.Contacts
	if err := global.Db.Where("owner_id = ? AND remark LIKE ?", accountID, "%"+keyword+"%").Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询联系人失败"})
		return
	}

	var targetIDs []uint
	for _, contact := range contacts {
		targetIDs = append(targetIDs, contact.ContactID)
	}
	for _, account := range accounts {
		targetIDs = append(targetIDs, account.AccountID)
	}

	var chatList []models.ChatInfo
	if len(targetIDs) > 0 {
		if err := global.Db.Where("target_id IN ?", targetIDs).Find(&chatList).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询聊天记录失败"})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []gin.H{}})
		return
	}

	var response []gin.H
	for _, chat := range chatList {
		var friend models.AccountInfo
		var groupChat models.GroupChatInfo

		// 判断好友/群聊是否存在
		if chat.IsGroup {
			if err := global.Db.Where("group_id = ?", chat.TargetID).First(&groupChat).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群聊信息不存在"})
				return
			}
		} else {
			if err := global.Db.Where("account_id = ?", chat.TargetID).First(&friend).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "好友信息不存在"})
			}
		}

		// 查询Contacts表
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, chat.TargetID).First(&contact).Error; err != nil {
			continue // 查询联系人信息失败则跳过
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

		var chatResponse gin.H
		if chat.IsGroup {
			avatarBase64, _, err := utils.GetFileContentAndType(groupChat.GroupAvatar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
				return
			}

			chatResponse = gin.H{
				"id":              chat.TargetID,
				"avatar":          avatarBase64,
				"name":            groupChat.GroupName,
				"remark":          contact.Remark,
				"lastMessage":     lastMessage.Content,
				"lastMessageTime": lastMessage.CreateTime,
				"unreadCount":     unreadCount,
			}
		} else {
			avatarBase64, _, err := utils.GetFileContentAndType(friend.Avatar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
				return
			}

			chatResponse = gin.H{
				"id":              chat.TargetID,
				"avatar":          avatarBase64,
				"name":            friend.Nickname,
				"remark":          contact.Remark,
				"lastMessage":     lastMessage.Content,
				"lastMessageTime": lastMessage.CreateTime,
				"unreadCount":     unreadCount,
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
		IsPinned bool   `json:"is_pinned"`
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

	// 查询Contacts表
	var contact models.Contacts
	if err := global.Db.Where("account_id = ? AND contact_id = ?", me.AccountID, other.AccountID).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Contacts表中无此条记录"})
		return
	}
	contact.IsPinned = input.IsPinned
	if err := global.Db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "置顶状态更新成功"})
}

// ReadMessages 标记为已读或未读消息
func ReadMessages(c *gin.Context) {
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
		Tid     string `json:"tid"`
		IsRead  bool   `json:"is_read"`
		IsGroup bool   `json:"is_group"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json绑定失败"})
		return
	}

	if input.IsGroup {
		var friend models.AccountInfo
		if err := global.Db.Where("account_id = ?", input.Tid).First(&friend).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
			return
		}
		if friend.Deactivate == true {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
			return
		}

		// 查询Contact表
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", me.AccountID, friend.AccountID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Contacts表中无此条记录"})
			return
		}

		if input.IsRead {
			contact.UnreadMessageNum = 0
		} else {
			contact.UnreadMessageNum = contact.UnreadMessageNum + 1
		}
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
	} else {
		var group models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", input.Tid).First(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊失败"})
			return
		}

		// 查询Contact表
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", me.AccountID, group.GroupID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Contacts表中无此条记录"})
			return
		}
		if input.IsRead {
			contact.UnreadMessageNum = 0
		} else {
			contact.UnreadMessageNum = contact.UnreadMessageNum + 1
		}
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "消息状态更新成功"})
}

// DeleteChat 删除聊天信息
func DeleteChat(c *gin.Context) {
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
		Tid     uint `json:"tid"`
		IsGroup bool `json:"is_group"`
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
	if other.Deactivate {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
		return
	}

	// 查询聊天记录是否存在
	var chat models.ChatInfo
	if err = global.Db.Where("account_id = ? AND target_id = ? AND is_group = ?", accountID, input.Tid, input.IsGroup).First(&chat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "聊天记录未找到"})
		return
	}

	// 删除聊天记录
	if err = global.Db.Where("account_id = ? AND target_id = ? AND is_group = ?", accountID, input.Tid, input.IsGroup).Delete(&models.ChatInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除聊天记录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "聊天记录删除成功"})
}

// SetMute 设置静音
func SetMute(c *gin.Context) {
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
		Tid     string `json:"tid"`
		IsMute  bool   `json:"is_mute"`
		IsGroup bool   `json:"is_group"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json绑定失败"})
		return
	}

	if input.IsGroup {
		var group models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", input.Tid).First(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群聊不存在"})
			return
		}

		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", me.AccountID, group.GroupID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Contacts表中无此条记录"})
		}
		contact.IsMute = input.IsMute
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
	} else {
		var other models.AccountInfo
		if err := global.Db.Where("account_id = ?", input.Tid).First(&other).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
			return
		}
		if other.Deactivate == true {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
			return
		}
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", me.AccountID, other.AccountID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Contacts表中无此条记录"})
		}
		contact.IsMute = input.IsMute
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "消息状态更新成功"})
}

// BlockChat 屏蔽聊天
func BlockChat(c *gin.Context) {
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
		Tid       string `json:"tid"`
		IsBlocked bool   `json:"is_blocked"`
		IsGroup   bool   `json:"is_group"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json绑定失败"})
		return
	}

	if input.IsGroup {
		var group models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", input.Tid).First(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群聊不存在"})
			return
		}

		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", me.AccountID, group.GroupID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Contacts表中无此条记录"})
		}
		contact.IsBlocked = input.IsBlocked
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
	} else {
		var other models.AccountInfo
		if err := global.Db.Where("account_id = ?", input.Tid).First(&other).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
			return
		}
		if other.Deactivate == true {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户注销"})
			return
		}
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", me.AccountID, other.AccountID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Contacts表中无此条记录"})
		}
		contact.IsBlocked = input.IsBlocked
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "消息状态更新成功"})
}

// ---------------------------------------------------------------------------
// 聊天消息

// GetMessages 获取聊天消息
func GetMessages(c *gin.Context) {
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
		Tid     uint `json:"tid"`
		IsGroup bool `json:"is_group"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json绑定失败"})
		return
	}

	if input.IsGroup == true {
		// 查询GroupChatInfo表
		var group models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", input.Tid).First(&group).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
			return
		}

		// 查询GroupMemberInfo表
		var groupMember models.GroupMemberInfo
		if err := global.Db.Where("account_id = ? AND group_id = ?", me.AccountID, group.GroupID).First(&groupMember).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "GroupMemberInfo表查询失败"})
			return
		}

		// 查询聊天记录
		var chat models.ChatInfo
		if err := global.Db.Where("account_id = ? AND target_id = ? AND is_group", me.AccountID, group.GroupID, true).First(&chat).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "聊天记录不存在"})
			return
		}

		// 查询当前聊天记录下的message
		var messages []models.MessageInfo
		if err := global.Db.Where("chat_id = ?", chat.ChatID).Order("create_time DESC").Find(&messages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询messageInfo表失败"})
			return
		}

		// 标记消息已读
		if err := global.Db.Model(&messages).Updates(map[string]interface{}{"is_read": true}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新消息状态失败"})
			return
		}

		// 返回消息
		var result []gin.H
		for _, message := range messages {
			// 查询GroupMemberInfo表
			var senderGroupMember models.GroupMemberInfo
			if err := global.Db.Where("account_id = ?", message.SendAccountID).First(&senderGroupMember).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询发送者失败"})
				return
			}

			// 查询AccountInfo表
			var sender models.AccountInfo
			if err := global.Db.Where("account_id = ?", message.SendAccountID).First(&sender).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询发送者失败"})
				return
			}

			avatarBase64, _, err := utils.GetFileContentAndType(sender.Avatar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
				return
			}

			result = append(result, gin.H{
				"message_id":      message.MessageID,
				"send_account_id": message.SendAccountID,
				"content":         message.Content,
				"sender":          senderGroupMember.GroupNickname,
				"avatar":          avatarBase64,
				"create_time":     message.CreateTime,
				"type":            message.Type,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功",
			"data": gin.H{
				"is_all_banned": group.IsAllBanned,
				"is_banned":     groupMember.IsBanned,
				"group_role":    groupMember.GroupRole,
				"messages":      result,
			},
		})
	} else {
		var friend models.AccountInfo
		if err := global.Db.Where("account_id = ?", input.Tid).First(&friend).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不存在"})
			return
		}
		if friend.Deactivate {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
			return
		}

		// 查询聊天记录
		var chat models.ChatInfo
		if err := global.Db.Where("account_id = ? AND target_id = ? AND is_group = ?", me.AccountID, friend.AccountID, false).First(&chat).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "聊天记录不存在"})
			return
		}

		// 查询当前聊天记录下的message
		var messages []models.MessageInfo
		if err := global.Db.Where("chat_id = ?", chat.ChatID).Order("create_time DESC").Find(&messages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询messageInfo表失败"})
			return
		}

		if len(messages) == 0 {
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "当前聊天记录为空"})
			return
		}

		// 标记消息已读
		if err := global.Db.Model(&messages).Updates(map[string]interface{}{"is_read": true}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新消息状态失败"})
			return
		}

		// 返回消息
		var result []gin.H
		for _, message := range messages {
			// 查询Contacts表
			// 接收者对发送者的contacts表
			var receiver_sender models.Contacts
			if err := global.Db.Where("owner_id = ? AND contact_id = ?", message.TargetID, message.SendAccountID).First(&receiver_sender).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询contacts表失败"})
				return
			}

			// 查询AccountInfo表
			var sender models.AccountInfo
			if err := global.Db.Where("account_id = ?", message.SendAccountID).First(&sender).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询发送者失败"})
				return
			}

			avatarBase64, mimeType, err := utils.GetFileContentAndType(friend.Avatar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
				return
			}
			avatarBase64 = "data:" + mimeType + ";base64," + avatarBase64

			result = append(result, gin.H{
				"message_id":      message.MessageID,
				"send_account_id": message.SendAccountID,
				"content":         message.Content,
				"sender":          receiver_sender.Remark,
				"avatar":          avatarBase64,
				"create_time":     message.CreateTime,
				"type":            message.Type,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功",
			"data": gin.H{
				"is_all_banned": nil,
				"is_banned":     nil,
				"group_role":    nil,
				"messages":      result,
			},
		})
	}
}

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
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
	if me.Deactivate {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		Tid     uint   `json:"tid"`
		Content string `json:"content"`
		Type    string `json:"type"`
		IsGroup bool   `json:"is_group"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	if input.IsGroup {
		// 查询群聊是否存在
		var group models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", input.Tid).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
			return
		}
		// 查询是否被禁言
		if group.IsAllBanned {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "所有成员都被禁言"})
			return
		}

		// 查询是否为群成员
		var groupMember models.GroupMemberInfo
		if err := global.Db.Where("account_id = ? AND group_id = ?", me.AccountID, group.GroupID).First(&groupMember).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不在该群聊中"})
			return
		}

		// 单个成员是否被禁言
		if groupMember.IsBanned {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "id用户被禁言"})
			return
		}

		// 查询当前用户与该群聊的聊天记录，如果没有则创建聊天记录
		var chat models.ChatInfo
		if err := global.Db.Where("account_id = ? AND target_id = ?", me.AccountID, group.GroupID).First(&chat).Error; err != nil {
			chat = models.ChatInfo{
				AccountID:  me.AccountID,
				TargetID:   group.GroupID,
				IsGroup:    true,
				CreateTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			if err := global.Db.Create(&chat).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建聊天记录失败"})
				return
			}
		}

		// 创建消息
		message := models.MessageInfo{
			SendAccountID: me.AccountID,
			TargetID:      group.GroupID,
			//ChatID:        chat.ChatID,
			Content:    input.Content,
			Type:       input.Type,
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := global.Db.Create(&message).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存消息失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "消息发送成功",
			"data": gin.H{
				"message_id":      message.MessageID,
				"create_time":     message.CreateTime,
				"send_account_id": message.SendAccountID,
				"target_id":       message.TargetID,
				"content":         message.Content,
				"type":            message.Type,
				//"chat_id":         message.ChatID,
				"is_read": false,
			},
		})

	} else {
		// 查询目标用户是否存在和注销
		var friend models.AccountInfo
		if err := global.Db.Where("account_id = ?", input.Tid).First(&friend).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不存在"})
			return
		}
		if friend.Deactivate {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "目标用户已注销"})
			return
		}

		// 如果没有则创建两个聊天记录
		var chat_sender, chat_receiver models.ChatInfo
		if err := global.Db.Where("account_id = ? AND target_id = ?", me.AccountID, friend.AccountID).First(&chat_sender).Error; err != nil {
			chat_sender = models.ChatInfo{
				AccountID:  me.AccountID,
				TargetID:   friend.AccountID,
				IsGroup:    false,
				CreateTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			if err := global.Db.Create(&chat_sender).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建发送者聊天记录失败"})
				return
			}
		}
		if err := global.Db.Where("account_id = ? AND target_id = ?", friend.AccountID, me.AccountID).First(&chat_receiver).Error; err != nil {
			chat_receiver = models.ChatInfo{
				AccountID:  friend.AccountID,
				TargetID:   me.AccountID,
				IsGroup:    false,
				CreateTime: time.Now().Format("2006-01-02 15:04:05"),
			}
			if err := global.Db.Create(&chat_receiver).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建接收者聊天记录失败"})
				return
			}
		}

		// 创建发送人的消息
		message := models.MessageInfo{
			SendAccountID:  me.AccountID,
			TargetID:       friend.AccountID,
			SenderChatID:   chat_sender.ChatID,
			ReceiverChatID: chat_receiver.ChatID,
			Content:        input.Content,
			Type:           input.Type,
			IsRead:         false,
			CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
		}
		if err := global.Db.Create(&message).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存消息失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "消息发送成功",
			"data": gin.H{
				"message_id":       message.MessageID,
				"create_time":      message.CreateTime,
				"send_account_id":  message.SendAccountID,
				"target_id":        message.TargetID,
				"content":          message.Content,
				"type":             message.Type,
				"sender_chat_id":   message.SenderChatID,
				"receiver_chat_id": message.ReceiverChatID,
				"is_read":          false,
			},
		})
	}
}

// CollectMessage 收藏消息
func CollectMessage(c *gin.Context) {
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
	if me.Deactivate {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		MessageID int `json:"message_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数错误"})
		return
	}
}

// ReplyMessage 回复消息
func ReplyMessage(c *gin.Context) {

}

// ForwardMessage 转发聊天记录
func ForwardMessage(c *gin.Context) {

}

// DeleteMessage 删除消息
func DeleteMessage(c *gin.Context) {
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
	if me.Deactivate {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		MessageID uint `json:"message_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	// 查询MessageInfo表，确认该消息是用户的
	var message models.MessageInfo
	if err := global.Db.Where("message_id = ? AND (send_account_id = ? OR target_id = ?)", input.MessageID, accountID, accountID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "消息不存在当前用户的聊天中"})
		return
	}

	// 查询DeleteInfo表，判断该消息是否已移入删除列表
	var deleteInfo models.DeleteInfo
	if err := global.Db.Where("message_id = ?", input.MessageID).First(&deleteInfo).Error; err == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该消息已经移入删除记录"})
		return
	}

	// 将消息ID移入DeleteInfo表
	deleteRecord := models.DeleteInfo{
		AccountID: uint(accountID),
		MessageID: input.MessageID,
	}
	if err := global.Db.Create(&deleteRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "移入删除记录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "消息删除成功"})
}
