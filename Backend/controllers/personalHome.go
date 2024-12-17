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
