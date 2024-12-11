package groupChat

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TransferOwnership - 群主转让
func TransferOwnership(c *gin.Context) {
	groupID := c.Param("id")
	var request struct {
		NewOwnerAccountID uint `json:"new_owner_account_id" binding:"required"`
	}

	// 解析请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 检查群主权限
	accountID := c.MustGet("account_id").(uint)
	var currentOwner models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ? AND is_banned = false", groupID, accountID).First(&currentOwner).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the group owner"})
		return
	}

	// 检查新群主是否为群成员
	var newOwner models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ? AND is_banned = false", groupID, request.NewOwnerAccountID).First(&newOwner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "New owner is not a member of the group"})
		return
	}

	// 更新群主
	if err := db.Model(&models.GroupChatInfo{}).Where("group_id = ?", groupID).Update("group_owner_id", request.NewOwnerAccountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer ownership"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group ownership transferred successfully"})
}

// AddAdministrator - 添加管理员
func AddAdministrator(c *gin.Context) {
	groupID := c.Param("id")
	var request struct {
		AccountID uint `json:"account_id" binding:"required"`
	}

	// 解析请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 检查群主权限
	accountID := c.MustGet("account_id").(uint)
	var currentOwner models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ? AND is_banned = false", groupID, accountID).First(&currentOwner).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the group owner"})
		return
	}

	// 检查目标成员是否为群成员
	var member models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ? AND is_banned = false", groupID, request.AccountID).First(&member).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found or is banned"})
		return
	}

	// 检查当前管理员数量
	var adminCount int64
	if err := db.Model(&models.GroupMemberInfo{}).Where("group_id = ? AND is_admin = true", groupID).Count(&adminCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check admin count"})
		return
	}

	// 检查管理员是否已满
	if adminCount >= 5 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot add more than 5 administrators"})
		return
	}

	// 更新目标成员为管理员
	if err := db.Model(&models.GroupMemberInfo{}).Where("group_id = ? AND account_id = ?", groupID, request.AccountID).Update("is_admin", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add administrator"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added as administrator successfully"})
}
