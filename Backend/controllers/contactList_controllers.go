package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// ---------------------------------------------------------------------------
/*搜索添加陌生人*/

// SearchStrangers 搜索陌生人
func SearchStrangers(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中缺少用户ID"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID无效"})
		return
	}
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户已经注销"})
		return
	}
	var input struct {
		Key string `json:"key"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少关键字参数或关键字为空"})
		return
	}

	var strangers []gin.H

	// 查找账号
	var accounts []models.AccountInfo
	err = global.Db.Where("id LIKE ? OR phone LIKE ? OR email LIKE ? OR account_id LIKE ?", "%"+input.Key+"%", "%"+input.Key+"%", "%"+input.Key+"%", "%"+input.Key+"%").Find(&accounts).Error
	if err == nil {
		for _, account := range accounts {
			strangers = append(strangers, gin.H{
				"tid":      account.AccountID,
				"id":       account.ID,
				"nickname": account.Nickname,
				"avatar":   account.Avatar,
				"type":     "friend",
			})
		}
	}

	// 查找群组
	var groupchats []models.GroupChatInfo
	err = global.Db.Where("group_id LIKE ?", "%"+input.Key+"%").Find(&groupchats).Error
	if err == nil {
		for _, groupchat := range groupchats {
			strangers = append(strangers, gin.H{
				"tid":      groupchat.GroupID,
				"id":       groupchat.GroupID,
				"nickname": groupchat.GroupName,
				"avatar":   groupchat.GroupAvatar,
				"type":     "group",
			})
		}
	}

	// 判断是否有结果返回
	if len(strangers) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "没有找到匹配的结果"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "查询成功", "data": strangers})
	}
}

// ChangeRemark 修改备注
func ChangeRemark(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中缺少用户ID"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID无效"})
		return
	}
	var me models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&me).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if me.Deactivate {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户已经注销"})
		return
	}
	var input struct {
		ID          string `json:"id"`
		IsGroupChat bool   `json:"is_group_chat"`
		Remark      string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json绑定失败"})
		return
	}

	if input.IsGroupChat {
		var group models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", input.ID).First(&group).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群聊不存在"})
			return
		}
		var contacts models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, input.ID).First(&contacts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "当前id不在群聊中"})
			return
		}
		contacts.Remark = input.Remark
		if err := global.Db.Save(&contacts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "修改失败"})
			return
		}
	} else {
		var other models.AccountInfo
		if err := global.Db.Where("account_id = ?", input.ID).First(&other).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户不存在"})
			return
		}
		if other.Deactivate {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户已经注销"})
			return
		}
		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, input.ID).First(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "好友不存在"})
			return
		}
		contact.Remark = input.Remark
		if err := global.Db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "修改失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "修改成功"})
}

//-----------------------------------------------------------------------------

/*好友请求*/

// GetFriendRequests 获取当前用户的所有好友请求列表
func GetFriendRequests(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中缺少用户ID"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID无效"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// 获取所有收到和发出的好友申请
	var receivedApplyInfos, sentApplyInfos []models.ApplyInfo
	err = global.Db.Where("receiver_id = ? AND apply_type = ?", accountID, "friend").Find(&receivedApplyInfos).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取收到的好友申请失败"})
		return
	}

	err = global.Db.Where("sender_id = ? AND apply_type = ?", accountID, "friend").Find(&sentApplyInfos).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取发出的好友申请失败"})
		return
	}

	// 如果没有任何请求
	if len(receivedApplyInfos) == 0 && len(sentApplyInfos) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "没有好友请求"})
		return
	}

	// 合并好友请求信息
	friendRequests := make([]map[string]interface{}, 0)

	// 处理收到的申请
	for _, applyInfo := range receivedApplyInfos {
		var senderInfo models.AccountInfo
		err := global.Db.Where("account_id = ?", applyInfo.SenderID).First(&senderInfo).Error
		if err != nil {
			continue // 跳过无法获取发送者信息的请求
		}

		friendRequest := map[string]interface{}{
			"apply_id":    applyInfo.ApplyID,
			"avatar":      senderInfo.Avatar,
			"name":        senderInfo.Nickname,
			"sender_id":   applyInfo.SenderID,
			"receiver_id": applyInfo.ReceiverID,
			"reason":      applyInfo.Reason,
			"status":      applyInfo.Status,
			"time":        time.Now().Format("2006-01-02 15:04:05"),
		}
		friendRequests = append(friendRequests, friendRequest)
	}

	// 处理发出的申请
	for _, applyInfo := range sentApplyInfos {
		var receiverInfo models.AccountInfo
		err := global.Db.Where("account_id = ?", applyInfo.ReceiverID).First(&receiverInfo).Error
		if err != nil {
			continue // 跳过无法获取接收者信息的请求
		}

		friendRequest := map[string]interface{}{
			"apply_id":    applyInfo.ApplyID,
			"avatar":      receiverInfo.Avatar,
			"name":        receiverInfo.Nickname,
			"sender_id":   applyInfo.SenderID,
			"receiver_id": applyInfo.ReceiverID,
			"reason":      applyInfo.Reason,
			"status":      applyInfo.Status,
			"time":        time.Now().Format("2006-01-02 15:04:05"),
		}
		friendRequests = append(friendRequests, friendRequest)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功！", "data": friendRequests})
}

// FriendRequestPend 处理好友请求(其他人申请当前id用户为好友)
func FriendRequestPend(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}
	var me models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "当前id用户不存在"})
		return
	}
	if me.Deactivate {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "当前id用户已经注销"})
		return
	}
	var input struct {
		AccountID uint `json:"account_id"`
		Accept    bool `json:"accept"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求体格式错误"})
		return
	}
	var other models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.AccountID).First(&other).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "申请好友的用户不存在"})
		return
	}
	if other.Deactivate {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "申请好友的用户已经注销"})
		return
	}

	// 判断是否有待处理的好友申请
	var apply models.ApplyInfo
	err = global.Db.Where("receiver_id = ? AND sender_id = ? AND apply_type = ? AND status = ?", accountID, input.AccountID, "friend", "pending").First(&apply).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该好友申请不存在或已被处理"})
		return
	}

	// 更新申请状态
	if input.Accept {
		apply.Status = "accepted"
	} else {
		apply.Status = "rejected"
	}
	if err := global.Db.Save(&apply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "无法更新申请状态"})
		return
	}

	// 如果接受请求，处理好友关系
	if input.Accept {

		// me对other的联系人关系
		meContact := models.Contacts{
			OwnerID:     uint(accountID),
			ContactID:   input.AccountID,
			IsBlocked:   false,
			IsPinned:    false,
			Divide:      "未分组",
			IsMute:      false,
			IsBlacklist: false,
			IsGroupChat: false,
			Remark:      "",
		}

		// 检查是否已经存在该联系人关系
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, input.AccountID).First(&meContact).Error; err != nil {
			// 插入联系人关系
			if err := global.Db.Create(&meContact).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加好友失败"})
				return
			}
		}

		// 检查是否已经有未分组
		var meDivide models.FriendDivide
		err := global.Db.Where("account_id = ? AND fd_name = ?", accountID, "未分组").First(&meDivide).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				meDivide = models.FriendDivide{
					AccountID: uint(accountID),
					FDName:    "未分组",
				}
				if err := global.Db.Create(&meDivide).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加到未分组失败"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询未分组失败"})
				return
			}
		}

		// other对me的联系人关系
		otherContact := models.Contacts{
			OwnerID:     input.AccountID,
			ContactID:   uint(accountID),
			IsBlocked:   false,
			IsPinned:    false,
			Divide:      "未分组",
			IsMute:      false,
			IsBlacklist: false,
			IsGroupChat: false,
			Remark:      "",
		}

		if err := global.Db.Where("owner_id = ? AND contact_id = ?", input.AccountID, accountID).First(&otherContact).Error; err != nil {
			if err := global.Db.Create(&otherContact).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加好友失败"})
				return
			}
		}

		// 检查是否已经有未分组
		var otherDivide models.FriendDivide
		err = global.Db.Where("account_id = ? AND fd_name = ?", input.AccountID, "未分组").First(&otherDivide).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				otherDivide = models.FriendDivide{
					AccountID: input.AccountID,
					FDName:    "未分组",
				}
				if err := global.Db.Create(&otherDivide).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加到未分组失败"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询未分组失败"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "添加好友成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "已拒绝好友请求"})
	}
}

// AddFriend 申请添加好友
func AddFriend(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var me models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户不存在"})
		return
	}
	if me.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户已注销，无法进行操作"})
		return
	}
	var input struct {
		AccountID uint   `json:"account_id"` // 被申请人ID
		Reason    string `json:"reason"`     // 申请理由
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}
	var other models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.AccountID).First(&other).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不存在"})
		return
	}
	if other.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户已注销，无法添加好友"})
		return
	}

	// 检查是否已经是好友
	var contacts models.Contacts
	err = global.Db.Where("(owner_id = ? AND contact_id = ?) OR (owner_id = ? AND contact_id = ?)", me.AccountID, other.AccountID, other.AccountID, me.AccountID).First(&contacts).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "你们已经是好友，请勿重复添加"})
		return
	}

	// 检查是否已发送过好友申请
	var existingApply models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND apply_type = ?", me.AccountID, other.AccountID, "friend").First(&existingApply).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "好友申请已发送，请勿重复申请"})
		return
	}

	newApply := models.ApplyInfo{
		ApplyType:  "friend",        // 申请类型
		SenderID:   me.AccountID,    // 发送者ID
		ReceiverID: other.AccountID, // 接收者ID
		Status:     "pending",       // 初始状态
		Reason:     input.Reason,    // 申请理由
	}

	if err := global.Db.Create(&newApply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "好友申请保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "好友申请已发送", "apply_id": newApply.ApplyID})
}

// ---------------------------------------------------------------------------
/*群聊请求*/

// DealGroupInviteRequest 处理群聊邀请请求
func DealGroupInviteRequest(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}
	var me models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if me.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"` // 邀请人的ID
		GroupID   uint `json:"group_id"`   // 群聊 ID
		Accept    bool `json:"accept"`     // 是否接受邀请
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "json解析失败"})
		return
	}

	var other models.AccountInfo
	if err = global.Db.Where("account_id = ?", input.AccountID).First(&other).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邀请人不存在"})
		return
	}
	if other.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邀请人已注销"})
		return
	}

	// 检查群聊是否存在
	var groupChat models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&groupChat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 检查ApplyInfo中是否存在这个邀请
	var applyInfo models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND group_id = ? AND apply_type = ?", input.AccountID, accountID, input.GroupID, "groupInvitation").First(&applyInfo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "数据库中没有这个群聊邀请申请记录"})
		return
	}

	if input.Accept { // 接受邀请，更新申请状态为 accepted
		applyInfo.Status = "accepted"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新申请表信息失败"})
			return
		}

		// 将用户添加到GroupMemberInfo中
		groupMember := models.GroupMemberInfo{
			AccountID:     uint(accountID),
			GroupID:       input.GroupID,
			GroupNickname: groupChat.GroupName,
			IsBanned:      false,
			GroupRole:     "group_ordinary",
		}
		if err := global.Db.Create(&groupMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "将用户添加到群聊中失败"})
			return
		}

		// 生成contacts关系表
		contact := models.Contacts{
			OwnerID:     uint(accountID),
			ContactID:   input.GroupID,
			IsGroupChat: true,
			IsPinned:    false,
			Divide:      "未分组",
			IsMute:      false,
			IsBlocked:   false,
			IsBlacklist: false,
			Remark:      "",
		}
		if err := global.Db.Create(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成contacts关系表失败"})
			return
		}

		// 检查GroupDivide表是否已经存在“未分组”分组
		var groupDivide models.GroupDivide
		err = global.Db.Where("account_id = ? AND gd_name = ?", input.AccountID, "未分组").First(&groupDivide).Error

		// 如果未找到"未分组"分组，才需要创建新的分组
		if errors.Is(err, gorm.ErrRecordNotFound) {
			groupdivide := models.GroupDivide{
				AccountID: uint(accountID),
				GDName:    "未分组",
			}
			if err := global.Db.Create(&groupdivide).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成群聊分组失败"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "已接受群聊邀请"})
	} else { // 拒绝邀请，更新申请状态为 rejected
		applyInfo.Status = "rejected"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新申请表信息失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "已拒绝群聊邀请"})
	}
}

// DealGroupApplyRequest 处理群聊申请请求
func DealGroupApplyRequest(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}
	AccountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID转换失败"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
		GroupID   uint `json:"group_id"`
		Accept    bool `json:"accept"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json解析失败"})
		return
	}

	// 检查申请记录是否存在
	var applyInfo models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND group_id = ? AND apply_type = ?", input.AccountID, AccountID, input.GroupID, "groupApply").First(&applyInfo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "数据库中没有这个群申请！"})
		return
	}

	if input.Accept { // 接受邀请
		applyInfo.Status = "accepted"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新申请表信息失败"})
			return
		}

		// 查询该群聊
		var groupChat models.GroupChatInfo
		err := global.Db.Where("group_id = ?", input.GroupID).First(&groupChat).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
			return
		}

		// 修改GroupMemberInfo表
		groupMember := models.GroupMemberInfo{
			AccountID:     uint(AccountID),
			GroupID:       input.GroupID,
			GroupNickname: groupChat.GroupName,
			IsBanned:      false,
			GroupRole:     "group_ordinary",
		}
		if err := global.Db.Create(&groupMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "将用户添加到群聊中失败"})
			return
		}

		// 修改Contacts表
		contacts := models.Contacts{
			OwnerID:     input.AccountID,
			ContactID:   input.GroupID,
			IsBlacklist: false,
			IsPinned:    false,
			Divide:      "未分组",
			IsMute:      false,
			IsBlocked:   false,
			IsGroupChat: true,
			Remark:      "",
		}
		if err := global.Db.Create(&contacts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "修改Contacts表失败"})
			return
		}

		// 检查GroupDivide表是否已经存在“未分组”分组
		var groupDivide models.GroupDivide
		err = global.Db.Where("account_id = ? AND gd_name = ?", input.AccountID, "未分组").First(&groupDivide).Error

		// 如果未找到"未分组"分组，才需要创建新的分组
		if errors.Is(err, gorm.ErrRecordNotFound) {
			divide := models.GroupDivide{
				AccountID: input.AccountID,
				GDName:    "未分组",
			}
			if err := global.Db.Create(&divide).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "修改GroupDivide表失败"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "同意群聊加入申请"})
	} else { // 拒绝邀请
		applyInfo.Status = "rejected"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "更新申请表信息失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "拒绝群聊加入申请"})
	}
}

// GetGroupRequests 获取群组申请列表
func GetGroupRequests(c *gin.Context) {
	// 从请求头获取用户ID
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	// 转换用户ID为整数
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}

	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id用户已被注销"})
		return
	}

	// 获取当前用户的所有群聊申请（包括作为群主的和作为申请者的，以及申请类型的和邀请类型的）
	var applyInfos []models.ApplyInfo
	err = global.Db.Where("(receiver_id = ? OR sender_id = ? ) AND (apply_type = ? OR apply_type = ?)", accountID, accountID, "groupInvitation", "groupApply").Find(&applyInfos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "数据库查询失败"})
		return
	}

	var groupRequests []map[string]interface{}
	for _, applyInfo := range applyInfos {
		var groupChat models.GroupChatInfo
		// 判断此群是否存在
		err := global.Db.Where("group_id = ?", applyInfo.GroupID).First(&groupChat).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "这个群没找到"})
			return
		}

		var groupMember models.GroupMemberInfo
		// 申请人SenderID是否已经在群聊中
		err = global.Db.Where("account_id = ? AND group_id = ?", applyInfo.SenderID, applyInfo.GroupID).First(&groupMember).Error
		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "申请人已经在群聊中"})
			return
		}

		// 获取申请人的信息
		var accountInfo models.AccountInfo
		err = global.Db.Where("account_id = ?", applyInfo.SenderID).First(&accountInfo).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取申请人信息失败"})
			return
		}
		if accountInfo.Deactivate == true {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "申请人已被注销"})
			return
		}

		groupRequest := map[string]interface{}{
			"apply_id":     applyInfo.ApplyID,                        // 申请ID
			"avatar":       groupChat.GroupAvatar,                    // 群聊头像
			"group_name":   groupChat.GroupName,                      // 群名称
			"account_name": accountInfo.Nickname,                     // 申请人或群主的昵称
			"sender_id":    applyInfo.SenderID,                       // 申请人或群主的ID
			"receiver_id":  applyInfo.ReceiverID,                     // 接收者ID
			"group_id":     applyInfo.GroupID,                        // 群聊ID
			"reason":       applyInfo.Reason,                         // 申请理由
			"apply_type":   applyInfo.ApplyType,                      // 申请类型
			"status":       applyInfo.Status,                         // 状态：pending、accepted、rejected等
			"time":         time.Now().Format("2006-01-02 15:04:05"), // 申请时间
		}

		groupRequests = append(groupRequests, groupRequest)
	}

	c.JSON(http.StatusOK, gin.H{"success": false, "message": "成功!", "群聊申请请求": groupRequests})
}

// AddGroup 申请加入群聊
func AddGroup(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户不存在"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户已注销，无法进行操作"})
		return
	}
	var input struct {
		GroupID uint   `json:"group_id"` // 群组ID
		Reason  string `json:"reason"`   // 申请理由
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查用户是否已经是群成员
	var membership models.GroupMemberInfo
	err = global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, user.AccountID).First(&membership).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "你已经是该群的成员，请勿重复申请"})
		return
	}

	// 检查是否已发送过加入申请
	var existingApply models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND group_id = ? AND apply_type = ?", user.AccountID, input.GroupID, "groupInvitation").First(&existingApply).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "加入群组申请已发送，请勿重复申请"})
		return
	}

	// 构造新的群组申请记录
	newApply := models.ApplyInfo{
		ApplyType:  "groupInvitation", // 申请类型
		SenderID:   user.AccountID,    // 申请者ID
		ReceiverID: group.GroupOwner,  // 群主作为接收方
		Status:     "pending",         // 初始状态
		Reason:     input.Reason,      // 申请理由
	}

	if err := global.Db.Create(&newApply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群组申请保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "加入群组申请已发送", "apply_id": newApply.ApplyID})
}

// ---------------------------------------------------------------------------
/*黑名单*/

// GetBlackList 获取当前用户的黑名单列表
func GetBlackList(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	// 转换用户ID为整数
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}

	// 查询Contacts表，查找当前用户的黑名单联系人
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ? AND is_blacklist = ? AND is_group_chat = ?", accountID, true, false).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询黑名单失败"})
		return
	}

	// 获取黑名单用户的详细信息
	var blackList []gin.H
	for _, contact := range contacts {
		var account models.AccountInfo
		err := global.Db.Where("account_id = ?", contact.ContactID).First(&account).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "查询用户信息失败"})
			return
		}

		blackList = append(blackList, gin.H{
			"avatar":     account.Avatar,
			"name":       account.Nickname,
			"account_id": account.AccountID,
			"signature":  account.Signature,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"blacklist": blackList,
	})
}

// RemoveFromBlacklist 从黑名单中移除用户
func RemoveFromBlacklist(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	var contact models.Contacts
	result := global.Db.Where("owner_id = ? AND contact_id = ? AND is_group_chat = ?", accountID, input.AccountID, false).First(&contact)

	// 如果没有找到该记录
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "没有找到这个关系"})
		return
	}

	// contact_id不在黑名单中
	if !contact.IsBlacklist {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "黑名单已经不在黑名单中"})
		return
	}

	contact.IsBlacklist = false
	if err := global.Db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "黑名单移除错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "黑名单移除成功"})
}

// AddToBlacklist 将用户添加到黑名单
func AddToBlacklist(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	AccountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "转换失败"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	var contact models.Contacts
	result := global.Db.Where("owner_id = ? AND contact_id = ? AND is_group_chat = ?", AccountID, input.AccountID, false).First(&contact)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "没有找到这个关系"})
		return
	}
	if contact.IsBlacklist {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "黑名单已经存在"})
		return
	}

	contact.IsBlacklist = true
	if err := global.Db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加黑名单错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "添加黑名单成功"})
}

// ---------------------------------------------------------------------------
/*好友列表*/

// GetFriends 获取好友列表
func GetFriends(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var me models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户不存在"})
		return
	}
	if me.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户已注销，无法进行操作"})
		return
	}

	// 定义好友返回结构体
	type FriendResponse struct {
		Avatar    string `json:"avatar"`
		AccountID uint   `json:"account_id"`
		Remark    string `json:"remark"`
		Status    string `json:"status"`
		Signature string `json:"signature"`
		Tag       string `json:"tag"`
	}

	// 查询联系人
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ?", accountID).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取当前ID的好友列表失败"})
		return
	}

	var friends []FriendResponse
	for _, contact := range contacts {
		if contact.IsBlacklist { // 黑名单不显示
			continue
		}
		var accountInfo models.AccountInfo
		err := global.Db.First(&accountInfo, contact.ContactID).Error
		if err != nil { // 如果某个联系人信息获取失败，继续下一个联系人
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取当前好友信息失败"})
			continue
		}

		friend := FriendResponse{
			Avatar:    accountInfo.Avatar,
			AccountID: accountInfo.AccountID,
			Remark:    contact.Remark,
			Status:    accountInfo.Status,
			Signature: accountInfo.Signature,
			Tag:       contact.Divide,
		}

		friends = append(friends, friend)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取好友列表成功", "data": friends})
}

// ---------------------------------------------------------------------------
/*分组部分*/

// GetDivides 获取分组名称
func GetDivides(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	var groups []string // 只存储分组名称的数组
	if groupType == "groups" {
		if err := global.Db.Model(&models.GroupDivide{}).Where("account_id = ?", accountID).Pluck("gd_name", &groups).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取群聊分组失败"})
			return
		}
	} else {
		if err := global.Db.Table("friend_divides").Where("account_id = ?", accountID).Pluck("fd_name", &groups).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取好友分组失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "divides": groups})
}

// CreateDivide 创建分组
func CreateDivide(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// 获取URL中的type参数
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 获取分组名称
	var input struct {
		Divide string `json:"divide" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "json绑定失败"})
		return
	}

	// 根据groupType选择不同的分组类型
	var existingDivide interface{}
	var divideName = input.Divide

	if groupType == "groups" {
		// 检查群聊分组是否存在
		existingDivide = &models.GroupDivide{}
		err = global.Db.Where("account_id = ? AND divide = ?", accountID, divideName).First(existingDivide).Error
	} else {
		// 检查好友分组是否存在
		existingDivide = &models.FriendDivide{}
		err = global.Db.Where("account_id = ? AND divide = ?", accountID, divideName).First(existingDivide).Error
	}

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": groupType + "分组名称已存在"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询分组名称失败"})
		return
	}

	// 创建新的分组
	if groupType == "groups" {
		// 创建群聊分组
		newDivide := models.GroupDivide{
			AccountID: uint(accountID),
			GDName:    divideName,
		}
		if err := global.Db.Create(&newDivide).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建群聊分组失败"})
			return
		}
	} else {
		// 创建好友分组
		newDivide := models.FriendDivide{
			AccountID: uint(accountID),
			FDName:    divideName,
		}
		if err := global.Db.Create(&newDivide).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建好友分组失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建分组成功"})
}

// DeleteDivide 删除分组
func DeleteDivide(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}
	var input struct {
		Divide string `json:"divide"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Divide == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "分组名称不能为空"})
		return
	}

	// 通用删除和成员移入未分组操作
	deleteAndMoveToUncategorized := func(groupName, groupType string) error {
		var group interface{}
		var err error
		if groupType == "groups" { // 根据groupType选择对应的分组模型
			group = &models.GroupDivide{}
			err = global.Db.Where("account_id = ? AND gd_name = ?", accountID, groupName).First(group).Error
		} else {
			group = &models.FriendDivide{}
			err = global.Db.Where("account_id = ? AND fd_name = ?", accountID, groupName).First(group).Error
		}

		if err != nil { // 分组不存在
			return fmt.Errorf("分组不存在")
		}
		if err := global.Db.Delete(group).Error; err != nil { // 删除分组
			return fmt.Errorf("删除分组失败")
		}

		// 更新分组中的成员，移入"未分类"分组
		if err := global.Db.Model(&models.Contacts{}).Where("owner_id = ? AND divide = ?", accountID, groupName).Update("divide", "未分组").Error; err != nil {
			return fmt.Errorf("成员分组更新失败")
		}

		return nil
	}

	// 执行删除和成员移入未分组的操作
	if err := deleteAndMoveToUncategorized(input.Divide, groupType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("%s分组删除成功，成员已移入未分类分组", groupType)})
}

// RenameDivide 重命名分组
func RenameDivide(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}

	// 获取URL中的type参数（friends 或 groups）
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 解析请求体中的参数
	var requestBody struct {
		OldFdName string `json:"old_fd_name"` // 旧分组名称
		NewFdName string `json:"new_fd_name"` // 新分组名称
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "JSON绑定失败"})
		return
	}
	if requestBody.OldFdName == "" || requestBody.NewFdName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "old_fd_name和new_fd_name参数不能为空"})
		return
	}

	// 确定是否是群聊
	var isGroupChat bool
	if groupType == "groups" {
		isGroupChat = true
	} else {
		isGroupChat = false
	}

	// 检查旧分组是否存在
	var oldGroup models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, requestBody.OldFdName, isGroupChat).First(&oldGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "旧分组不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询旧分组失败"})
		}
		return
	}

	// 检查新分组名称是否已存在
	var newGroup models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, requestBody.NewFdName, isGroupChat).First(&newGroup).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "新分组名称已存在"})
		return
	}

	// 更新旧分组的名称为新分组名称
	err = global.Db.Model(&oldGroup).Update("divide", requestBody.NewFdName).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新分组名称失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "分组名称更新成功"})
}

// MoveInDivide 将传入的成员的分组改为相应分组
func MoveInDivide(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户注销"})
		return
	}

	// 获取URL中的type参数（friends 或 groups）
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 解析请求体中的参数
	var requestBody struct {
		TID    string `json:"tid"`    // 成员ID或者群号
		Divide string `json:"divide"` // 目标分组名称
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "JSON绑定失败"})
		return
	}
	if requestBody.TID == "" || requestBody.Divide == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tid和divide参数不能为空"})
		return
	}

	// 确定是否是群聊
	var isGroupChat bool
	if groupType == "groups" {
		isGroupChat = true
	} else {
		isGroupChat = false
	}

	//检测进行分组的这个人是否存在，是否注销
	if !isGroupChat {
		var user1 models.AccountInfo
		if err := global.Db.Where("account_id = ?", requestBody.TID).First(&user1).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "需要进行分组的用户不存在"})
			return
		}
		if user1.Deactivate == true {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "需要进行分组的用户已经注销"})
			return
		}
	}

	// 检查目标分组是否存在
	var targetGroup models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, requestBody.Divide, isGroupChat).First(&targetGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "目标分组不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询目标分组失败"})
		}
		return
	}

	// 检查当前id用户的这个好友/群聊是否存在
	var member models.Contacts
	err = global.Db.Where("owner_id = ? AND tid = ? AND is_group_chat = ?", accountID, requestBody.TID, isGroupChat).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "成员不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询成员失败"})
		}
		return
	}

	err = global.Db.Model(&member).Update("divide", requestBody.Divide).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新成员分组失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成员已成功移入目标分组"})
}

// ---------------------------------------------------------------------------
/*群聊部分*/

// GetGroups 获取群聊列表
func GetGroups(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// 查询用户所属的群聊信息
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ? AND is_group_chat = ? AND is_blocked = ?", accountID, true, false).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊失败"})
		return
	}

	var groupList []map[string]interface{}

	for _, contact := range contacts {
		var group models.GroupChatInfo
		err := global.Db.Where("group_id = ?", contact.ContactID).First(&group).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊失败"})
			continue
		}

		// 将群聊信息添加到结果数组
		groupList = append(groupList, map[string]interface{}{
			"avatar":     group.GroupAvatar,       // 群头像
			"account_id": group.GroupID,           // 群ID
			"signature":  group.GroupIntroduction, // 群介绍
			"remark":     contact.Remark,          // 群名称或备注
			"tag":        contact.Divide,          // 分组名称
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "data": groupList})
}

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		GroupName        string `json:"group_name"`
		GroupAvatar      string `json:"group_avatar"`
		GroupDescription string `json:"group_description"`
		AllowInvite      bool   `json:"allow_invite"`
		AllowIDSearch    bool   `json:"allow_id_search"`
		AllowNameSearch  bool   `json:"allow_name_search"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数解析失败"})
		return
	}
	if input.GroupAvatar == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊头像不能为空"})
		return
	}
	if input.GroupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊名称不能为空"})
		return
	}

	// 写入群聊总表
	groupchat := models.GroupChatInfo{
		GroupAvatar:       input.GroupAvatar,
		GroupName:         input.GroupName,
		IsAllBanned:       false,
		GroupOwner:        uint(accountID),
		GroupIntroduction: input.GroupDescription,
		ALlowNameSearch:   input.AllowNameSearch,
		AllowIDSearch:     input.AllowIDSearch,
		AllowInvite:       input.AllowInvite,
	}
	if err := global.Db.Create(&groupchat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建群聊失败"})
		return
	}

	// 群聊表/通讯录表
	contact := models.Contacts{
		OwnerID:     uint(accountID),
		ContactID:   groupchat.GroupID,
		IsBlacklist: false,
		IsPinned:    false,
		IsGroupChat: true,
		Divide:      "未分组",
		IsMute:      false,
		IsBlocked:   false,
		Remark:      "",
	}
	if err := global.Db.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建群聊失败"})
		return
	}

	// 群成员表
	groupmember := models.GroupMemberInfo{
		AccountID:     uint(accountID),
		GroupID:       groupchat.GroupID,
		GroupNickname: input.GroupName,
		IsBanned:      false,
		GroupRole:     "group_owner",
	}
	if err := global.Db.Create(&groupmember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建群聊失败"})
		return
	}

	// 检查群分组表中是否已经有“未分类”的分组
	var existingGroupDivide models.GroupDivide
	err = global.Db.Where("gd_name = ?", "未分类").First(&existingGroupDivide).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊分组失败"})
		return
	}

	// 如果没有找到“未分类”分组，则添加新的分组
	if err == gorm.ErrRecordNotFound {
		groupdivide := models.GroupDivide{
			AccountID: uint(accountID),
			GDName:    "未分类",
		}
		if err := global.Db.Create(&groupdivide).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建群聊失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建群聊成功"})
}

// DisMissGroup 解散群聊
func DisMissGroup(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP Header 中用户 ID 为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID 解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		GroupID uint `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "JSON绑定失败"})
		return
	}

	// 判断群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 验证是否是群主才能解散群聊
	if uint(accountID) != group.GroupOwner {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主才能解散群聊"})
		return
	}

	// 删除 GroupChatInfo 表中的群聊记录
	if err := global.Db.Where("group_id = ?", input.GroupID).Delete(&models.GroupChatInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "解散群聊失败"})
		return
	}

	// 删除 Contacts 表中所有与该群聊相关的联系人关系
	if err := global.Db.Where("owner_id = ? AND contact_id = ? AND is_group_chat = ?", accountID, input.GroupID, true).Delete(&models.Contacts{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除contacts表中的群聊信息失败"})
		return
	}

	// 删除 GroupMemberInfo 表中所有该群聊的成员记录
	if err := global.Db.Where("group_id = ?", input.GroupID).Delete(&models.GroupMemberInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除GroupMemberInfo中群成员信息表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "解散群聊成功"})
}

// Invite 邀请其他人加入群聊
func Invite(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
		GroupID   uint `json:"group_id"`
	}

	//检测被邀请的用户是否存在
	var other models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.AccountID).First(&other).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "想邀请的用户不存在"})
		return
	}
	if other.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "想邀请的用户已注销"})
		return
	}

	// 群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 是否已经在群聊中
	var member models.GroupMemberInfo
	if err := global.Db.Where("account_id = ? AND group_id = ?", input.AccountID, input.GroupID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "该成员已经在群聊中"})
		return
	}

	newMember := models.GroupMemberInfo{
		GroupID:   input.GroupID,
		AccountID: input.AccountID,
		GroupRole: "group_ordinary",
	}
	if err := global.Db.Create(&newMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "邀请失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "邀请成功"})
}

// Quit 退出群聊
func Quit(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		GroupID string `json:"group_id"`
	}

	if input.GroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少必要的参数"})
		return
	}

	// 检测该群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 检测该用户是否为群聊成员
	var groupMember models.GroupMemberInfo
	if err := global.Db.Where("AccountID = ? AND group_id = ?", accountID, input.GroupID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该用户不是群聊成员，无法退出"})
		return
	}

	// 退出群聊
	if err := global.Db.Where("AccountID = ? AND group_id = ?", accountID, input.GroupID).Delete(&groupMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "退出群聊失败，请稍后重试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功退出群聊"})
}

// GetGroupInfo 获取群聊的详细信息
func GetGroupInfo(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	groupID := c.Param("group_id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少群聊ID参数"})
		return
	}

	// 检测群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", groupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 获取Group_id群成员信息
	var members []models.GroupMemberInfo
	if err := global.Db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取群聊成员失败"})
		return
	}

	memberList := make([]gin.H, len(members))
	for i, member := range members {
		memberList[i] = gin.H{
			"user_id":   member.AccountID,
			"user_name": member.GroupNickname,
		}
	}

	var groupChat models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", groupID).First(&groupChat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取群聊信息失败"})
		return
	}

	// 前后端可能返回不一致，需注意
	c.JSON(http.StatusOK, gin.H{
		"success":            true,
		"message":            "成功",
		"group_id":           groupID,
		"group_name":         group.GroupName,
		"group_introduction": groupChat.GroupIntroduction,
		"members":            memberList,
	})
}

// ChangeNickname 更改用户在群聊内的昵称
func ChangeNickname(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		GroupID       uint   `json:"group_id"`       // 群组ID
		GroupNickname string `json:"group_nickname"` // 新的群昵称
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 检查群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 检查用户是否是群成员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, accountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不是该群的成员"})
		return
	}

	// 更新用户的群昵称
	membership.GroupNickname = input.GroupNickname
	if err := global.Db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新昵称失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "昵称更新成功"})
}

// SetAllowInvite 设置是否允许群成员邀请他人
func SetAllowInvite(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		GroupID     uint `json:"group_id"`
		AllowInvite bool `json:"allow_invite"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数格式不正确"})
		return
	}

	// 查询群聊是否存在
	var groupchat models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&groupchat).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 查询GroupMember表
	var groupMember models.GroupMemberInfo
	if err := global.Db.Where("account_id = ? AND group_id = ?", accountID, input.GroupID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不在群聊中"})
	}

	//判断是否为管理员或群主
	if groupMember.GroupRole != "group_owner" && groupMember.GroupRole != "group_admin" {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不是群主或管理员，无权限修改设置"})
		return
	}

	// 更新群聊设置
	groupchat.AllowInvite = input.AllowInvite
	if err := global.Db.Save(&groupchat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "设置成功"})
}

// SetAllowIdSearch 设置群聊中能否通过ID搜索群中其他人
func SetAllowIdSearch(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		GroupID       uint `json:"group_id"`
		AllowIdSearch bool `json:"allow_id_search"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "json参数格式不正确"})
		return
	}

	// 查询群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 查询GroupMember表
	var groupMember models.GroupMemberInfo
	if err := global.Db.Where("account_id = ? AND group_id = ?", accountID, input.GroupID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不在群聊中"})
	}

	//判断是否为管理员或群主
	if groupMember.GroupRole != "group_owner" && groupMember.GroupRole != "group_admin" {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不是群主或管理员，无权限修改设置"})
		return
	}

	// 更改是否能通过ID搜索
	group.AllowIDSearch = input.AllowIdSearch
	if err := global.Db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "设置成功"})
}

// SetAllowNameSearch 设置群聊中能否通过名字搜索其他人
func SetAllowNameSearch(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		GroupID         uint `json:"group_id"`
		AllowNameSearch bool `json:"allow_name_search"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数格式不正确"})
		return
	}

	// 查询群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 查询GroupMember表
	var groupMember models.GroupMemberInfo
	if err := global.Db.Where("account_id = ? AND group_id = ?", accountID, input.GroupID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不在群聊中"})
	}

	//判断是否为管理员或群主
	if groupMember.GroupRole != "group_owner" && groupMember.GroupRole != "group_admin" {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不是群主或管理员，无权限修改设置"})
		return
	}

	// 修改是否能通过名称进行搜索
	group.ALlowNameSearch = input.AllowNameSearch
	if err := global.Db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "设置成功"})
}

// SetAllBanned 设置是否全体禁言
func SetAllBanned(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		GroupID     uint `json:"group_id"`
		IsAllBanned bool `json:"is_all_banned"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数格式不正确"})
		return
	}

	// 查询群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 查询GroupMember表
	var groupMember models.GroupMemberInfo
	if err := global.Db.Where("account_id = ? AND group_id = ?", accountID, input.GroupID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不在群聊中"})
	}

	//判断是否为管理员或群主
	if groupMember.GroupRole != "group_owner" && groupMember.GroupRole != "group_admin" {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不是群主或管理员，无权限修改设置"})
		return
	}

	// 更新群聊禁言状态
	group.IsAllBanned = input.IsAllBanned
	if err := global.Db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "设置成功"})
}

// FetchFriendsNotInGroup 搜索不在群聊内的好友
func FetchFriendsNotInGroup(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		GroupID uint `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 检查群聊是否存在
	var group models.GroupChatInfo
	if err = global.Db.Where("group_id=?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 获取用户的好友列表
	var friends []models.Contacts
	if err := global.Db.Where("owner_id = ?", accountID).Find(&friends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取好友列表失败"})
		return
	}

	// 获取已在群聊内的成员列表
	var groupMembers []models.GroupMemberInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).Find(&groupMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取群成员列表失败"})
		return
	}

	// 将已在群聊内的成员ID存入map
	memberIDs := make(map[uint]bool)
	for _, member := range groupMembers {
		memberIDs[member.AccountID] = true
	}

	// 筛选不在群聊内的好友
	var friendsNotInGroup []models.Contacts
	for _, friend := range friends {
		if !memberIDs[friend.ContactID] {
			friendsNotInGroup = append(friendsNotInGroup, friend)
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": friendsNotInGroup})
}

// RemoveMember 移除某个群成员
func RemoveMember(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 查找当前群组信息
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主或管理员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, accountID).First(&membership).Error; err != nil || (membership.GroupRole != "group_admin" && group.GroupOwner != uint(accountID)) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主或管理员可以移除成员"})
		return
	}

	// 检查目标用户是否是群成员
	if err := global.Db.Where("group_id = ? AND member_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员，没法移除"})
		return
	}

	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).Delete(&models.GroupMemberInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "移除成员失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功移除群成员"})
}

// SetBanned 设置禁言或取消禁言
func SetBanned(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
		IsBanned  bool `json:"is_banned"`  // 是否禁言
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 查找当前群组信息
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主或管理员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, accountID).First(&membership).Error; err != nil || (membership.GroupRole != "group_ordinary" && group.GroupOwner != uint(accountID)) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主或管理员可以设置禁言"})
		return
	}

	// 检查目标用户是否是群成员
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员"})
		return
	}

	// 更新禁言状态
	membership.IsBanned = input.IsBanned
	if err := global.Db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新禁言状态失败"})
		return
	}

	// 返回成功响应
	action := "设置禁言"
	if !input.IsBanned {
		action = "取消禁言"
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("成功%s", action)})
}

// SetAdmin 设置管理员
func SetAdmin(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
		IsAdmin   bool `json:"is_admin"`   // 是否设置为管理员
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主
	if group.GroupOwner != uint(accountID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主可以设置管理员"})
		return
	}

	// 检查目标用户是否是群成员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员"})
		return
	}

	if input.IsAdmin {
		membership.GroupRole = "group_admin"
	} else {
		membership.GroupRole = "group_ordinary"
	}
	if err := global.Db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新管理员状态失败"})
		return
	}

	action := "设置为管理员"
	if !input.IsAdmin {
		action = "取消管理员"
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("成功%s", action)})
}

// TransferOwner 更换群主
func TransferOwner(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 查找当前群组信息
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主
	if group.GroupOwner != uint(accountID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主可以转让群主权限"})
		return
	}

	// 检查目标用户是否是群成员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员"})
		return
	}

	// 转让群主权限
	if err := global.Db.Model(&group).Update("owner_id", input.AccountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群主转让失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "群主权限已成功转让"})
}

// ChangeGroupAvatar 更改群头像
func ChangeGroupAvatar(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		GroupID     uint   `json:"group_id"`
		GroupAvatar string `json:"group_avatar"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数格式不正确"})
		return
	}

	// 查询群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 查询GroupMember表
	var groupMember models.GroupMemberInfo
	if err := global.Db.Where("account_id = ? AND group_id = ?", accountID, input.GroupID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不在群聊中"})
	}

	//判断是否为管理员或群主
	if groupMember.GroupRole != "group_owner" && groupMember.GroupRole != "group_admin" {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不是群主或管理员，无权限修改设置"})
		return
	}

	// 更新群头像
	group.GroupAvatar = input.GroupAvatar
	if err := global.Db.Save(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "头像更新成功"})
}

// -------------------------------------------------------------------------
/* 获取各类资料卡片*/

// GetProfileCard 获取指定用户的资料卡片
func GetProfileCard(c *gin.Context) {
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

	// 用户ID或者群号
	other := c.Param("tid")
	groupID := c.DefaultPostForm("group_id", "") // 获取POST请求中的group_id参数
	if groupID != "" {
		var groupChat models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", groupID).First(&groupChat).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
			return
		}

		var contact models.Contacts
		if err := global.Db.Where("contact_id = ? AND owner_id = ?", groupID, accountID).First(&contact).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "不是群聊成员"})
			return
		}

		responseData := gin.H{
			"tid":           groupChat.GroupID,
			"avatar":        groupChat.GroupAvatar,
			"remark":        contact.Remark,
			"groupNickname": groupChat.GroupName,
			"tag":           contact.Divide,
			"signature":     groupChat.GroupIntroduction,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功",
			"data":    responseData,
		})
	} else {
		var another models.AccountInfo
		if err := global.Db.Where("account_id = ?", other).First(&another).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		if another.Deactivate {
			c.JSON(http.StatusNotFound, gin.H{"error": "目标用户已注销"})
			return
		}

		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, another.AccountID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "不是好友关系"})
			return
		}

		responseData := gin.H{
			"tid":       another.AccountID,
			"avatar":    another.Avatar,
			"remark":    contact.Remark,
			"nickname":  another.Nickname,
			"tag":       contact.Divide,
			"signature": another.Signature,
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功",
			"data":    responseData,
		})
	}
}
