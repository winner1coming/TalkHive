package settings

import (
	"TalkHive/config"
	"TalkHive/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ChangeTheme 更改主题颜色并存储到数据库
func ChangeTheme(c *gin.Context) {
	var req struct {
		ID    string `json:"id"`    // 用户不可修改的ID
		Theme string `json:"theme"` // 主题颜色: light / dark / system
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 检查传入的主题是否有效
	validThemes := []string{"light", "dark", "system"}
	isValid := false
	for _, t := range validThemes {
		if req.Theme == t {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid theme value",
		})
		return
	}

	// 更新数据库中的主题设置
	var setting models.SystemSetting
	if err := config.DB.First(&setting).Error; err != nil {
		// 如果没有找到系统设置记录，则创建一个新的设置记录
		setting = models.SystemSetting{
			Theme: req.Theme,
		}
		if err := config.DB.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to save theme to database",
			})
			return
		}
	} else {
		// 更新现有的系统设置记录
		setting.Theme = req.Theme
		if err := config.DB.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update theme in database",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Theme updated successfully",
	})
}

// ChangeFontsize 更改字体大小并存储到数据库
func ChangeFontsize(c *gin.Context) {
	var req struct {
		ID       string `json:"id"`       // 用户不可修改的ID
		FontSize int    `json:"fontSize"` // 字体大小
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 检查传入的字体大小是否有效（假设字体大小在10到100之间）
	if req.FontSize < 10 || req.FontSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Font size must be between 10 and 100",
		})
		return
	}

	// 获取当前系统设置
	var setting models.SystemSetting
	if err := config.DB.First(&setting).Error; err != nil {
		// 如果没有找到系统设置记录，则创建一个新的设置记录
		setting = models.SystemSetting{
			FontSize: req.FontSize,
		}
		if err := config.DB.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to save font size to database",
			})
			return
		}
	} else {
		// 更新现有的系统设置记录
		setting.FontSize = req.FontSize
		if err := config.DB.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update font size in database",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Font size updated successfully",
	})
}

// ChangeFontstyle 更改字体风格并存储到数据库
func ChangeFontstyle(c *gin.Context) {
	var req struct {
		ID        string `json:"id"`        // 用户不可修改的ID
		FontStyle string `json:"fontStyle"` // 字体风格
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 获取当前系统设置
	var setting models.SystemSetting
	if err := config.DB.First(&setting).Error; err != nil {
		// 如果没有找到系统设置记录，则创建一个新的设置记录
		setting = models.SystemSetting{
			FontStyle: req.FontStyle,
		}
		if err := config.DB.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to save font style to database",
			})
			return
		}
	} else {
		// 更新现有的系统设置记录
		setting.FontStyle = req.FontStyle
		if err := config.DB.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update font style in database",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Font style updated successfully",
	})
}

// UpdateNotification 更新消息通知设置
func UpdateNotification(c *gin.Context) {
	var req struct {
		ID     string `json:"id"`     // 用户不可修改的ID
		Notice bool   `json:"notice"` // 是否开启消息通知
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 获取当前系统设置
	var setting models.SystemSetting
	if err := config.DB.First(&setting).Error; err != nil {
		// 如果没有找到系统设置记录，则创建一个新的设置记录
		setting = models.SystemSetting{
			Sound: "on", // 默认声音开启
		}
		// 创建新的设置
		if err := config.DB.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to save notification settings to database",
			})
			return
		}
	} else {
		// 更新现有的系统设置记录
		if !req.Notice {
			// 如果关闭了消息通知，则同时关闭通知声音
			setting.Sound = "off"
		}
		// 更新通知设置
		// 注意：此处只更新消息通知的状态，声音设置由关闭通知时自动处理
		if err := config.DB.Model(&setting).Update("Sound", setting.Sound).Update("Sound", setting.Sound).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update notification settings in database",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification settings updated successfully",
	})
}
