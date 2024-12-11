package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 资料编辑相关

// UpdateAvatar 更新用户头像
func UpdateAvatar(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Avatar    string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Update("avatar", req.Avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully"})
}

// UpdateNickname 修改用户昵称
func UpdateNickname(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Nickname  string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Update("nickname", req.Nickname).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nickname"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Nickname updated successfully"})
}

// EditProfile 编辑其他用户资料（性别、签名等）
func EditProfile(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Gender    string `json:"gender"`
		Signature string `json:"signature"`
		ID        string `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Updates(models.AccountInfo{
		Gender:    req.Gender,
		Signature: req.Signature,
		ID:        req.ID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// 账号设置相关

// UpdatePhone 修改绑定手机号
func UpdatePhone(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Phone     string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Update("phone", req.Phone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update phone"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Phone updated successfully"})
}

// UpdatePassword 修改登录密码
func UpdatePassword(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Password  string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Update("password", req.Password).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// UpdateFriendPermission 设置好友添加权限
func UpdateFriendPermission(c *gin.Context) {
	var req struct {
		AccountID        uint   `json:"account_id"`
		FriendPermission string `json:"friend_permission"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Update("friend_permission", req.FriendPermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friend permission"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend permission updated successfully"})
}

// 系统设置相关

// UpdateSystemSettings 更新系统设置
func UpdateSystemSettings(c *gin.Context) {
	var req struct {
		AccountID  uint   `json:"account_id"`
		Background string `json:"background"`
		Theme      string `json:"theme"`
		Sound      string `json:"sound"`
		FontSize   int    `json:"font_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.SystemSetting{}).Where("account_id = ?", req.AccountID).Updates(models.SystemSetting{
		Background: req.Background,
		Theme:      req.Theme,
		Sound:      req.Sound,
		FontSize:   req.FontSize,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update system settings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "System settings updated successfully"})
}
