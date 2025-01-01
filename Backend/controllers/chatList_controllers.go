package controllers

//// GetChatList 聊天列表
//func GetChatList(c *gin.Context) {
//	userID := c.GetHeader("User-ID")
//	if userID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中缺少用户ID"})
//		return
//	}
//	accountID, err := strconv.Atoi(userID)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID无效"})
//		return
//	}
//	var user models.AccountInfo
//	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户不存在"})
//		return
//	}
//
//	if user.Deactivate {
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户已注销"})
//		return
//	}
//
//	// 查询聊天列表
//	var chats []models.Chat
//	err = global.Db.Where("user_id = ?", accountID).Order("last_message_time desc").Find(&chats).Error
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取聊天列表失败"})
//		return
//	}
//
//	// 构造聊天列表
//	chatList := make([]map[string]interface{}, 0)
//	for _, chat := range chats {
//		// 查询好友或群组信息
//		var chatPartner models.AccountInfo
//		err := global.Db.Where("account_id = ?", chat.Tid).First(&chatPartner).Error
//		if err != nil {
//			continue // 如果没有找到聊天对象，跳过
//		}
//
//		// 获取该聊天的最后一条消息
//		var lastMessage models.MessageInfo
//		err = global.Db.Where("chat_id = ?", chat.ChatID).Order("create_time desc").First(&lastMessage).Error
//		if err != nil {
//			continue // 如果获取不到最后一条消息，跳过
//		}
//
//		// 获取未读消息数
//		var unreadCount int64
//		err = global.Db.Model(&models.MessageInfo{}).Where("chat_id = ? AND send_account_id != ? AND status = ?", chat.ChatID, accountID, "unread").Count(&unreadCount).Error
//		if err != nil {
//			continue // 如果未读消息数查询失败，跳过
//		}
//
//		// 获取标签信息
//		var tags []string
//		if chat.IsPinned {
//			tags = append(tags, "pinned")
//		}
//		if unreadCount > 0 {
//			tags = append(tags, "unread")
//		}
//
//		// 组织每个聊天对象的信息
//		chatInfo := map[string]interface{}{
//			"id":              chat.Tid,
//			"avatar":          chatPartner.Avatar,
//			"name":            chatPartner.Nickname, // 或者使用 chat.Remark 根据需求
//			"lastMessage":     lastMessage.Content,
//			"lastMessageTime": lastMessage.CreateTime.Format("15:04"), // 格式化为 10:00 形式
//			"unreadCount":     unreadCount,
//			"tags":            tags,
//		}
//
//		chatList = append(chatList, chatInfo)
//	}
//	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功！", "chats": chatList})
//}
