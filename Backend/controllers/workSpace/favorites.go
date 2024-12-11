package workSpace

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 我的收藏！！！

// GetFavorites - 查看收藏列表
func GetFavorites(c *gin.Context) {
	var favorites []models.Favorites
	if err := DB.Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}
	c.JSON(http.StatusOK, favorites)
}

// ViewFavorite - 查看收藏内容
func ViewFavorite(c *gin.Context) {
	id := c.Param("id")
	var favorite models.Favorites
	if err := DB.First(&favorite, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}
	c.JSON(http.StatusOK, favorite)
}
