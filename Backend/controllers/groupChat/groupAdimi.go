package groupChat

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"TalkHive/models"
	"github.com/gin-gonic/gin"
)

// UpdateGroupAvatar 修改群头像
func UpdateGroupAvatar(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		AvatarURL string `json:"avatar_url"` // 新头像URL
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Model(&models.GroupChatInfo{}).
		Where("id = ?", groupID).
		Update("avatar_url", request.AvatarURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully"})
}

// UpdateGroupName 修改群名称
func UpdateGroupName(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		Name string `json:"name"` // 新群名称
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Model(&models.GroupChatInfo{}).
		Where("id = ?", groupID).
		Update("name", request.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update name"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group name updated successfully"})
}

// UpdateJoinPermissions 设置入群权限
func UpdateJoinPermissions(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		Permissions string `json:"permissions"` // 入群权限: open, invite_only, approval_required
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Model(&models.GroupChatInfo{}).
		Where("id = ?", groupID).
		Update("join_permissions", request.Permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update join permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Join permissions updated successfully"})
}

// UpdateGroupIntro 编辑群简介
func UpdateGroupIntro(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		Intro string `json:"intro"` // 新简介内容
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Model(&models.GroupChatInfo{}).
		Where("id = ?", groupID).
		Update("intro", request.Intro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update intro"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group intro updated successfully"})
}

// SetMuteStatus 设置禁言（个体或全员）
func SetMuteStatus(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		AccountID uint  `json:"account_id"` // 被禁言的用户ID，0表示全员禁言
		IsMuted   bool  `json:"is_muted"`   // 禁言状态
		Duration  int64 `json:"duration"`   // 禁言持续时间（秒），0表示永久
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 如果是全员禁言
	if request.AccountID == 0 {
		if err := db.Model(&models.GroupChatInfo{}).
			Where("id = ?", groupID).
			Update("is_muted", request.IsMuted).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mute status"})
			return
		}
	} else {
		// 单用户禁言
		muteEndTime := time.Now().Add(time.Duration(request.Duration) * time.Second)
		if err := db.Model(&models.GroupMemberInfo{}).
			Where("group_id = ? AND account_id = ?", groupID, request.AccountID).
			Updates(map[string]interface{}{
				"is_muted":   request.IsMuted,
				"mute_until": muteEndTime,
			}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mute member"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mute status updated successfully"})
}

// RemoveMember 移除群成员
func RemoveMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var request struct {
		AccountID uint `json:"account_id"` // 被移除的用户ID
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := db.Where("group_id = ? AND account_id = ?", groupID, request.AccountID).
		Delete(&models.GroupMemberInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// DisbandGroup 解散群聊
func DisbandGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// 删除群聊及相关信息
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", groupID).Delete(&models.GroupChatInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", groupID).Delete(&models.GroupMemberInfo{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disband group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group disbanded successfully"})
}

// HandleJoinApplications 管理入群申请
func HandleJoinApplications(c *gin.Context) {
	groupID := c.Param("id")
	var input struct {
		AccountID uint   `json:"account_id" binding:"required"`
		Status    string `json:"status" binding:"required"` // "approved" 或 "rejected"
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找入群申请
	var application models.ApplyInfo
	if err := db.First(&application, "group_id = ? AND account_id = ?", groupID, input.AccountID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// 更新申请状态
	application.Status = input.Status
	if err := db.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle application"})
		return
	}

	// 返回处理结果
	c.JSON(http.StatusOK, gin.H{"message": "Application handled successfully"})
}
