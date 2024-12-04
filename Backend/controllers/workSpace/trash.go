package workSpace

import (
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 搜索栏！！！

// SearchByKeyword - 根据关键字匹配内容
func SearchByKeyword(c *gin.Context) {
	keyword := c.Query("keyword")
	var notes []models.Notes
	var codes []models.Codes

	// 搜索笔记标题和内容
	noteQuery := DB.Where("is_show = ? AND (title LIKE ? OR content LIKE ?)", true, "%"+keyword+"%", "%"+keyword+"%")
	if err := noteQuery.Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search notes"})
		return
	}

	// 搜索代码文件标题和内容
	codeQuery := DB.Where("is_show = ? AND (title LIKE ? OR content LIKE ?)", true, "%"+keyword+"%", "%"+keyword+"%")
	if err := codeQuery.Find(&codes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search codes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
		"codes": codes,
	})
}

// 回收站！！！

// GetTrashItems - 查看回收站内容
func GetTrashItems(c *gin.Context) {
	var trashedNotes []models.Notes
	var trashedCodes []models.Codes

	// 获取回收站中的笔记
	if err := DB.Where("is_show = ?", false).Find(&trashedNotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trashed notes"})
		return
	}

	// 获取回收站中的代码文件
	if err := DB.Where("is_show = ?", false).Find(&trashedCodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trashed codes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trashedNotes": trashedNotes,
		"trashedCodes": trashedCodes,
	})
}

// RestoreNote - 恢复回收站笔记
func RestoreNote(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Model(&models.Notes{}).Where("note_id = ?", id).Update("is_show", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore note"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note restored successfully"})
}

// DeletePermanently - 永久删除
func DeletePermanently(c *gin.Context) {
	id := c.Param("id")

	// 删除笔记或代码
	if err := DB.Delete(&models.Notes{}, "note_id = ?", id).Error; err != nil {
		if err := DB.Delete(&models.Codes{}, "code_id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted permanently"})
}
