package workSpace

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 我的笔记！！！

// CreateNote - 新建笔记
func CreateNote(c *gin.Context) {
	var note models.Notes
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.IsShow = true
	if err := DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}
	c.JSON(http.StatusOK, note)
}

// GetAllNotes - 查看所有笔记
func GetAllNotes(c *gin.Context) {
	var notes []models.Notes
	if err := DB.Where("is_show = ?", true).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}
	c.JSON(http.StatusOK, notes)
}

// GetNote - 查看特定笔记
func GetNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Notes
	if err := DB.First(&note, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(http.StatusOK, note)
}

// EditNote - 编辑笔记
func EditNote(c *gin.Context) {
	id := c.Param("id")
	var updatedNote models.Notes
	if err := c.ShouldBindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := DB.Model(&models.Notes{}).Where("note_id = ?", id).Updates(updatedNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

// DeleteNote - 删除笔记
func DeleteNote(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Model(&models.Notes{}).Where("note_id = ?", id).Update("is_show", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

// GetNotesByCategory - 按分类查看笔记
func GetNotesByCategory(c *gin.Context) {
	category := c.Param("category")
	var notes []models.Notes
	if err := DB.Where("type = ? AND is_show = ?", category, true).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes by category"})
		return
	}
	c.JSON(http.StatusOK, notes)
}
