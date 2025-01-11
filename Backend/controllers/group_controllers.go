package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetAllowSearch(c *gin.Context) {
	var req struct {
		GroupID       uint `json:"group_id" binding:"required"`
		AllowIDSearch bool `json:"allow_id_search" binding:"required"`
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data", "error": err.Error()})
		return
	}

	// 修改 models.GroupChatInfo 表，更新对应 GroupID 条目的 AllowInvite 属性
	var groupChatInfo models.GroupChatInfo
	if err := global.Db.Model(&groupChatInfo).
		Where("group_id = ?", req.GroupID).
		Update("allow_id_search", req.AllowIDSearch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update allow_invite", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Group invite settings updated successfully"})
}
