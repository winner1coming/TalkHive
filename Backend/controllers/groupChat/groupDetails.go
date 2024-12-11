package groupChat

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetGroupDetails 显示群聊详情
func GetGroupDetails(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var group models.GroupChatInfo
	if err := db.Where("id = ?", groupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// 查询成员列表
	var members []models.GroupMemberInfo
	if err := db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group":   group,
		"members": members,
	})
}

// SearchMember 搜索群成员
func SearchMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		Query string `json:"query"` // 搜索关键词
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var members []models.GroupMemberInfo
	if err := db.Where("group_id = ? AND group_nickname LIKE ?", groupID, "%"+request.Query+"%").
		Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

// InviteMember 邀请好友加入
func InviteMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		AccountIDs []uint `json:"account_ids"` // 被邀请的好友ID列表
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 添加成员
	for _, accountID := range request.AccountIDs {
		newMember := models.GroupMemberInfo{
			GroupID:       uint(groupID),
			AccountID:     accountID,
			GroupNickname: "", // 可根据需求设置默认昵称
		}
		if err := db.Create(&newMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite member"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Members invited successfully"})
}

// SetGroupNickname 设置群昵称
func SetGroupNickname(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	accountID := c.MustGet("account_id").(uint)
	var request struct {
		Nickname string `json:"nickname"` // 新昵称
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Model(&models.GroupMemberInfo{}).
		Where("group_id = ? AND account_id = ?", groupID, accountID).
		Update("group_nickname", request.Nickname).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nickname"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nickname updated successfully"})
}

// SetGroupNote 设置群聊备注
func SetGroupNote(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	accountID := c.MustGet("account_id").(uint)
	var request struct {
		Note string `json:"note"` // 新备注
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Model(&models.GroupMemberInfo{}).
		Where("group_id = ? AND account_id = ?", groupID, accountID).
		Update("note", request.Note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

// UpdateGroupSettings 更新群聊设置
func UpdateGroupSettings(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	accountID := c.MustGet("account_id").(uint)
	var request struct {
		IsMuted  bool `json:"is_muted"`  // 消息免打扰
		IsPinned bool `json:"is_pinned"` // 是否置顶
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Model(&models.GroupMemberInfo{}).
		Where("group_id = ? AND account_id = ?", groupID, accountID).
		Updates(map[string]interface{}{
			"is_muted":  request.IsMuted,
			"is_pinned": request.IsPinned,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// ExitGroup 退出群聊
func ExitGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	accountID := c.MustGet("account_id").(uint)
	if err := db.Where("group_id = ? AND account_id = ?", groupID, accountID).
		Delete(&models.GroupMemberInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exit group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exited group successfully"})
}
