package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddToFavorites 添加到收藏
func AddToFavorites(c *gin.Context) {
	var favorite models.Favorites
	if err := c.ShouldBindJSON(&favorite); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites"})
}
