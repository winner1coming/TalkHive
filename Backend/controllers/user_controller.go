package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 用户注册
func Register(c *gin.Context) {
	var account models.AccountInfo
	// 绑定 JSON 数据到 account 对象
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名或邮箱是否已被注册
	var existingAccount models.AccountInfo
	if err := config.DB.Where("id = ? OR email = ?", account.ID, account.Email).First(&existingAccount).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ID or Email already exists"})
		return
	}

	// 创建新用户
	if err := config.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account registered successfully"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var loginRequest struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户并验证密码
	var account models.AccountInfo
	if err := config.DB.Where("id = ? AND password = ?", loginRequest.ID, loginRequest.Password).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"nickname": account.Nickname,
		"email":    account.Email,
	})
}

// GetProfile 获取用户信息
func GetProfile(c *gin.Context) {
	id := c.Param("id")

	var account models.AccountInfo
	if err := config.DB.Where("id = ?", id).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": account.AccountID,
		"id":         account.ID,
		"nickname":   account.Nickname,
		"email":      account.Email,
		"avatar":     account.Avatar,
		"signature":  account.Signature,
		"gender":     account.Gender,
		"birthday":   account.Birthday,
	})
}

// UpdateProfile 更新用户信息
func UpdateProfile(c *gin.Context) {
	id := c.Param("id")

	var account models.AccountInfo
	if err := config.DB.Where("id = ?", id).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Signature string `json:"signature"`
		Gender    string `json:"gender"`
		Birthday  string `json:"birthday"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新信息
	account.Nickname = updateData.Nickname
	account.Avatar = updateData.Avatar
	account.Signature = updateData.Signature
	account.Gender = updateData.Gender

	if err := config.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
