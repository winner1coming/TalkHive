package workSpace

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// DDL记录！！！

// CreateDDL - 新建DDL事件
func CreateDDL(c *gin.Context) {
	var ddl models.DDLS
	if err := c.ShouldBindJSON(&ddl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ddl.IsCompleted = false
	if err := DB.Create(&ddl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DDL"})
		return
	}
	c.JSON(http.StatusOK, ddl)
}

// GetAllDDL - 查看所有DDL事件
func GetAllDDL(c *gin.Context) {
	var ddls []models.DDLS
	if err := DB.Find(&ddls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DDLs"})
		return
	}
	c.JSON(http.StatusOK, ddls)
}

// GetDDL - 查看单个DDL事件
func GetDDL(c *gin.Context) {
	id := c.Param("id")
	var ddl models.DDLS
	if err := DB.First(&ddl, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DDL not found"})
		return
	}
	c.JSON(http.StatusOK, ddl)
}

// EditDDL - 编辑DDL事件
func EditDDL(c *gin.Context) {
	id := c.Param("id")
	var updatedDDL models.DDLS
	if err := c.ShouldBindJSON(&updatedDDL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Model(&models.DDLS{}).Where("ddl_id = ?", id).Updates(updatedDDL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DDL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DDL updated successfully"})
}

// MarkDDLComplete - 勾选完成事件
func MarkDDLComplete(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Model(&models.DDLS{}).Where("ddl_id = ?", id).Update("is_complete", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark DDL as complete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DDL marked as complete"})
}

// DeleteDDL - 删除DDL事件
func DeleteDDL(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Delete(&models.DDLS{}, "ddl_id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DDL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DDL deleted successfully"})
}

// GetUpcomingDDL - 查看主页面DDL提醒
func GetUpcomingDDL(c *gin.Context) {
	var upcomingDDL []models.DDLS
	if err := DB.Where("is_complete = ? AND deadline >= ?", false, time.Now()).Order("deadline asc").Limit(5).Find(&upcomingDDL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming DDLs"})
		return
	}
	c.JSON(http.StatusOK, upcomingDDL)
}
