package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 获取数据库连接
var DB = config.GetDB()

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

// ShowProfile 获取用户个人主页资料
func ShowProfile(c *gin.Context) {
	// 获取前端传递的 id 参数
	id := c.Query("id")

	// 检查是否传递了 id 参数
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Missing required parameter: id",
		})
		return
	}

	// 查询数据库
	var account models.AccountInfo
	if err := config.DB.Where("account_id = ?", id).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// 返回符合前端格式的 JSON 数据
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User profile retrieved successfully",
		"data": gin.H{
			"id":        account.AccountID,
			"nickname":  account.Nickname,
			"gender":    account.Gender,
			"birthday":  account.Birthday,
			"signature": account.Signature,
		},
	})
}

// SaveEdit 保存编辑后的用户信息
func SaveEdit(c *gin.Context) {
	// 接收前端发送的 JSON 数据
	var updateData struct {
		ID        string `json:"id" binding:"required"`
		Avatar    string `json:"avatar"`
		Nickname  string `json:"nickname"`
		Gender    string `json:"gender"`
		Birthday  string `json:"birthday"` // 字符串形式
		Signature string `json:"signature"`
	}

	// 绑定 JSON 数据并检查必填字段
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data: " + err.Error(),
		})
		return
	}

	// 查询数据库中是否存在对应的用户
	var account models.AccountInfo
	if err := config.DB.Where("account_id = ?", updateData.ID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// 更新用户信息
	account.Nickname = updateData.Nickname
	account.Avatar = updateData.Avatar
	account.Signature = updateData.Signature
	account.Gender = updateData.Gender

	// 转换生日为 time.Time 类型
	if updateData.Birthday != "" {
		parsedBirthday, err := time.Parse("2006-01-02", updateData.Birthday) // 假设前端格式为 "YYYY-MM-DD"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid date format for birthday. Expected format: YYYY-MM-DD",
			})
			return
		}
		account.Birthday = parsedBirthday
	}

	// 保存更新后的数据到数据库
	if err := config.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
		})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
	})
}
