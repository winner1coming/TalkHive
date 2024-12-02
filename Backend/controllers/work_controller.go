package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetNotes 获取笔记
func GetNotes(c *gin.Context) {
	userID := c.Param("user_id")
	var notes []models.Notes
	if err := config.DB.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

// AddCode 添加代码片段
func AddCode(c *gin.Context) {
	var code models.Codes
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&code).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add code snippet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Code snippet added successfully"})
}

// ManageTodos 管理代办事项
func ManageTodos(c *gin.Context) {
	var todo models.DDLS
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to manage todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo managed successfully"})
}
