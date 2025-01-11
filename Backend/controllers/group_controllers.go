package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetAllowInvite(c *gin.Context) {
	var req struct {
		GroupID     uint `json:"group_id" binding:"required"`
		AllowInvite bool `json:"allow_invite" binding:"required"`
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
		Update("allow_invite", req.AllowInvite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update allow_invite", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Group invite settings updated successfully"})
}

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

func SetAllowNameSearch(c *gin.Context) {
	var req struct {
		GroupID         uint `json:"group_id" binding:"required"`
		AllowNameSearch bool `json:"allow_name_search" binding:"required"`
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
		Update("allow_name_search", req.AllowNameSearch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update allow_invite", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Group invite settings updated successfully"})
}

func SetAllBanned(c *gin.Context) {
	// group_id, is_all_banned
	var req struct {
		GroupID     uint `json:"group_id" binding:"required"`
		IsAllBanned bool `json:"is_all_banned" binding:"required"`
	}

	// 绑定请求体数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data", "error": err.Error()})
		return
	}

	// 修改models.GroupChatInfo表，更新对应GroupID条目的IsAllBanned属性
	var groupChatInfo models.GroupChatInfo
	if err := global.Db.Model(&groupChatInfo).
		Where("group_id = ?", req.GroupID).
		Update("is_all_banned", req.IsAllBanned).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update group invite settings", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Group invite settings updated successfully"})
}

func ChangeGroupAvatar(c *gin.Context) {
	// group_id, group_avatar
	var req struct {
		GroupID     uint   `json:"group_id" binding:"required"`
		GroupAvatar string `json:"group_avatar" binding:"required"`
	}

	// 绑定请求体数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form data", "error": err.Error()})
		return
	}

	// 修改 models.GroupChatInfo 表，更新对应 GroupID 条目的 GroupAvatar 属性
	var groupChatInfo models.GroupChatInfo
	if err := global.Db.Model(&groupChatInfo).
		Where("group_id = ?", req.GroupID).
		Update("group_avatar", req.GroupAvatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update group avatar", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Group avatar updated successfully"})
}
