package singleChat

import (
	"TalkHive/config"
	"TalkHive/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// 获取数据库连接
var db = config.GetDB()

// GetFriendInfo 查看好友个人信息
func GetFriendInfo(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint) // 假设账户ID通过JWT等方式传递

	// 查找好友信息
	var friend models.AccountInfo
	if err := db.Where("account_id = ?", friendID).First(&friend).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// 查找用户的备注
	var contact models.Contacts
	if err := db.Where("account_id = ? AND contact_id = ?", accountID, friendID).First(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve contact information"})
		return
	}

	// 返回好友信息和备注
	c.JSON(http.StatusOK, gin.H{
		"friend": friend,
		"remark": contact.Remark,
	})
}

// SetFriendRemark 设置好友备注
func SetFriendRemark(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)
	var request struct {
		Remark string `json:"remark"`
	}

	// 解析请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 更新备注
	if err := db.Model(&models.Contacts{}).
		Where("account_id = ? AND contact_id = ?", accountID, friendID).
		Update("remark", request.Remark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update remark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Remark updated successfully"})
}

// AddToBlacklist 加入黑名单
func AddToBlacklist(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)

	// 更新联系人状态为黑名单
	if err := db.Model(&models.Contacts{}).
		Where("account_id = ? AND contact_id = ?", accountID, friendID).
		Update("is_blacklist", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added to blacklist"})
}

// RemoveFromBlacklist 移除黑名单
func RemoveFromBlacklist(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)

	// 更新联系人状态为非黑名单
	if err := db.Model(&models.Contacts{}).
		Where("account_id = ? AND contact_id = ?", accountID, friendID).
		Update("is_blacklist", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from blacklist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed from blacklist"})
}

// DeleteContact 删除联系人
func DeleteContact(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)

	// 删除联系人记录
	if err := db.Where("account_id = ? AND contact_id = ?", accountID, friendID).Delete(&models.Contacts{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

// SetMute 设置免打扰
func SetMute(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)

	// 更新免打扰状态
	if err := db.Model(&models.Contacts{}).
		Where("account_id = ? AND contact_id = ?", accountID, friendID).
		Update("is_mute", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set mute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mute set successfully"})
}

// SetPin 设置置顶
func SetPin(c *gin.Context) {
	friendID := c.Param("id")
	accountID := c.MustGet("account_id").(uint)

	// 更新置顶状态
	if err := db.Model(&models.Contacts{}).
		Where("account_id = ? AND contact_id = ?", accountID, friendID).
		Update("is_pinned", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set pin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pinned successfully"})
}
