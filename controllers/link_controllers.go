package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// AddLinks 控制器函数
func AddLinks(c *gin.Context) {
	// 获取请求头中的用户ID
	userID := c.GetHeader("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 定义请求体结构
	var seq struct {
		URL     string `json:"url" binding:"required"`
		URLName string `json:"url_name" binding:"required"`
		Icon    string `json:"icon" binding:"required"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 检查数据库中是否已有相同 URL 的数据项
	var existingLink models.Links
	if err := global.Db.Where("account_id = ? AND url = ?", userID, seq.URL).First(&existingLink).Error; err == nil {
		// 如果数据库中已经存在该 URL，返回失败信息
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "URL already exists",
		})
		return
	}

	// 如果不存在，创建新的链接记录
	newLink := models.Links{
		AccountID: userID,
		URL:       seq.URL,
		URLName:   seq.URLName,
		Icon:      seq.Icon,
	}

	// 插入数据到数据库
	if err := global.Db.Create(&newLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to add link",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Link added successfully",
	})
}

func DelLinks(c *gin.Context) {
	// 获取请求头中的用户ID
	userID := c.GetHeader("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 定义请求体结构
	var seq struct {
		URL string `json:"url" binding:"required"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&seq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 检查数据库中是否存在该 URL
	var existingLink models.Links
	if err := global.Db.Where("account_id = ? AND url = ?", userID, seq.URL).First(&existingLink).Error; err != nil {
		// 如果数据库中没有找到该 URL，返回失败信息
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Link not found",
			})
		} else {
			// 其他错误处理
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to check link existence",
			})
		}
		return
	}

	// 删除链接
	if err := global.Db.Where("account_id = ? AND url = ?", userID, seq.URL).Delete(&models.Links{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete link",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Link deleted successfully",
	})
}
