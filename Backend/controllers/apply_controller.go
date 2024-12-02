package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SubmitApply 提交申请
func SubmitApply(c *gin.Context) {
	var applyInfo models.ApplyInfo
	if err := c.ShouldBindJSON(&applyInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&applyInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
}

// ReviewApply 审核申请
func ReviewApply(c *gin.Context) {
	var review struct {
		ApplyID uint   `json:"apply_id"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.ApplyInfo{}).Where("apply_id = ?", review.ApplyID).Update("status", review.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application reviewed successfully"})
}
