package workSpace

import (
	"TalkHive/global"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// DB 获取数据库连接
var DB = global.Db

// 我的代码！！！

// CreateCode - 新建代码文件
func CreateCode(c *gin.Context) {
	var code models.Codes
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code.IsShow = true
	code.SaveTime = time.Now()
	if err := DB.Create(&code).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create code"})
		return
	}
	c.JSON(http.StatusOK, code)
}

// GetCode - 查看代码文件
func GetCode(c *gin.Context) {
	id := c.Param("id")
	var code models.Codes
	if err := DB.First(&code, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code not found"})
		return
	}
	c.JSON(http.StatusOK, code)
}

// EditCode - 编辑代码文件
func EditCode(c *gin.Context) {
	id := c.Param("id")
	var updatedCode models.Codes
	if err := c.ShouldBindJSON(&updatedCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Model(&models.Codes{}).Where("code_id = ?", id).Updates(updatedCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Code updated successfully"})
}

// ShareCode - 分享代码文件
func ShareCode(c *gin.Context) {
	id := c.Param("id")
	// 分享逻辑需要根据具体需求实现，例如生成分享链接
	c.JSON(http.StatusOK, gin.H{"message": "Code shared successfully", "id": id})
}
