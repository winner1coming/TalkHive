package controllers

import (
	"net/http"
	"strconv"

	"TalkHive/models"
	"github.com/gin-gonic/gin"
)

// SendApplication 发送好友/群聊申请
func SendApplication(c *gin.Context) {
	var application models.ApplyInfo

	// 绑定请求体参数
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存到数据库
	if err := DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application sent successfully"})
}

// AcceptApplication 接受好友申请/被通过
func AcceptApplication(c *gin.Context) {
	accountID := c.PostForm("account_id")
	applyType := c.PostForm("apply_type")
	groupID := c.DefaultPostForm("group_id", "")

	var application models.ApplyInfo

	// 查找申请记录
	if err := DB.Where("account_id = ? AND apply_type = ?", accountID, applyType).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// 更新状态为已接受
	application.Status = "accepted"
	if err := DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept application"})
		return
	}

	// 如果是好友申请，更新联系人表
	if applyType == "friend" {
		contact := models.Contacts{
			ContactID:   application.AccountID,
			IsGroupChat: false,
		}
		if err := DB.Create(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friend list"})
			return
		}
	}

	// 如果是群聊申请，加入群聊
	if applyType == "group" {
		groupIDUint, _ := strconv.Atoi(groupID)
		groupMember := models.GroupMemberInfo{
			AccountID: uint(application.AccountID),
			GroupID:   uint(groupIDUint),
		}
		if err := DB.Create(&groupMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add group member"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application accepted"})
}

// RemoveFriend 删除好友/被删除
func RemoveFriend(c *gin.Context) {
	accountID := c.Query("account_id")

	// 删除好友记录
	if err := DB.Delete(&models.Contacts{}, "contact_id = ?", accountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed successfully"})
}

// GetApplicationStatus 获取好友申请状态
func GetApplicationStatus(c *gin.Context) {
	accountID := c.Query("account_id")
	var applications []models.ApplyInfo

	// 查询申请记录
	if err := DB.Where("account_id = ?", accountID).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch application status"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// AcceptGroupApplication 入群申请被接受
func AcceptGroupApplication(c *gin.Context) {
	applyID := c.PostForm("apply_id")
	var application models.ApplyInfo

	// 查找申请记录
	if err := DB.First(&application, "apply_id = ?", applyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// 更新状态
	application.Status = "accepted"
	if err := DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept application"})
		return
	}

	// 更新群聊成员列表
	groupMember := models.GroupMemberInfo{
		AccountID: application.AccountID,
		GroupID:   application.GroupID,
	}
	if err := DB.Create(&groupMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add group member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group application accepted"})
}

// LeaveGroup 踢出群聊/被踢出
func LeaveGroup(c *gin.Context) {
	accountID := c.PostForm("account_id")
	groupID := c.PostForm("group_id")

	// 删除群聊成员记录
	if err := DB.Delete(&models.GroupMemberInfo{}, "account_id = ? AND group_id = ?", accountID, groupID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left the group"})
}

// GetGroupApplicationStatus 获取自己发送的申请群聊的状态
func GetGroupApplicationStatus(c *gin.Context) {
	accountID := c.Query("account_id")
	var applications []models.ApplyInfo

	// 查询记录
	if err := DB.Where("account_id = ? AND apply_type = 'group'", accountID).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group application status"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// AddGroup 添加分组
func AddGroup(c *gin.Context) {
	groupName := c.PostForm("group_name")
	accountIDStr := c.PostForm("account_id") // account_id 作为字符串获取

	// 将 accountID 字符串转换为 uint
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32) // 转换为 uint 类型
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// 插入新分组
	groupDivide := models.GroupDivide{
		GDName:    groupName,
		AccountID: uint(accountID), // 将 uint 转换为模型需要的 uint 类型
	}

	if err := DB.Create(&groupDivide).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group added successfully"})
}

// DeleteGroup 删除分组
func DeleteGroup(c *gin.Context) {
	groupName := c.Query("group_name")
	accountID := c.Query("account_id")

	// 删除分组记录
	if err := DB.Delete(&models.GroupDivide{}, "gd_name = ? AND account_id = ?", groupName, accountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// RenameGroup 重命名分组
func RenameGroup(c *gin.Context) {
	oldName := c.PostForm("old_group_name")
	newName := c.PostForm("new_group_name")
	accountID := c.PostForm("account_id")

	// 更新分组名称
	if err := DB.Model(&models.GroupDivide{}).Where("gd_name = ? AND account_id = ?", oldName, accountID).Update("gd_name", newName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group renamed successfully"})
}

// GetGroupStatistics 获取分组人数统计
func GetGroupStatistics(c *gin.Context) {
	accountID := c.MustGet("account_id").(uint)

	// 查询分组统计数据
	var statistics []struct {
		Divide string `json:"divide"`
		Count  int    `json:"count"`
	}

	if err := DB.Model(&models.Contacts{}).
		Select("divide, COUNT(*) as count").
		Where("account_id = ?", accountID).
		Group("divide").
		Find(&statistics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve group statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statistics": statistics})
}

// GetSummary 获取总好友数和群聊数
func GetSummary(c *gin.Context) {
	accountID := c.MustGet("account_id").(uint)

	var friendCount int64
	var groupCount int64

	// 查询好友总数
	if err := DB.Model(&models.Contacts{}).
		Where("account_id = ? AND is_groupchat = ?", accountID, false).
		Count(&friendCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve friend count"})
		return
	}

	// 查询群聊总数
	if err := DB.Model(&models.Contacts{}).
		Where("account_id = ? AND is_groupchat = ?", accountID, true).
		Count(&groupCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve group count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"friend_count": friendCount,
		"group_count":  groupCount,
	})
}

// SearchContact 搜索好友或群聊
func SearchContact(c *gin.Context) {
	accountID := c.MustGet("account_id").(uint)
	keyword := c.Query("keyword")

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword is required"})
		return
	}

	var friends []models.AccountInfo
	var groups []models.GroupChatInfo

	// 搜索好友
	if err := DB.Table("contacts").
		Select("account_infos.*").
		Joins("JOIN account_infos ON contacts.contact_id = account_infos.account_id").
		Where("contacts.account_id = ? AND contacts.is_groupchat = ? AND (account_infos.nickname LIKE ? OR account_infos.phone LIKE ?)", accountID, false, "%"+keyword+"%", "%"+keyword+"%").
		Scan(&friends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search friends"})
		return
	}

	// 搜索群聊
	if err := DB.Table("contacts").
		Select("group_chat_infos.*").
		Joins("JOIN group_chat_infos ON contacts.contact_id = group_chat_infos.group_id").
		Where("contacts.account_id = ? AND contacts.is_groupchat = ? AND group_chat_infos.group_name LIKE ?", accountID, true, "%"+keyword+"%").
		Scan(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"friends": friends,
		"groups":  groups,
	})
}
