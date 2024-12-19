package settings

import (
	"TalkHive/config"
	"TalkHive/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

// IsNotice 更新全局消息通知设置
func IsNotice(c *gin.Context) {
	var req struct {
		ID     string `json:"id"`     // 用户不可修改的 ID
		Notice bool   `json:"notice"` // 是否开启全局消息通知
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
	if err := config.DB.Where("id = ?", req.ID).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "System settings not found",
		})
		return
	}

	// 更新全局消息通知设置
	setting.Notice = req.Notice

	// 保存更新
	if err := config.DB.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update global notification settings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Global notification settings updated successfully",
	})
}

// IsNoticeGroup 更新群消息声音通知设置
func IsNoticeGroup(c *gin.Context) {
	var req struct {
		ID          string `json:"id"`          // 用户不可修改的 ID
		NoticeGroup bool   `json:"noticeGroup"` // 是否接收群消息声音通知
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
	if err := config.DB.Where("id = ?", req.ID).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "System settings not found",
		})
		return
	}

	// 检查全局通知开关
	if !setting.Notice {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Global notification is disabled. Cannot enable group notifications.",
		})
		return
	}

	// 更新群消息声音通知设置
	setting.NoticeGroup = req.NoticeGroup

	// 保存更新
	if err := config.DB.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update group notification settings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Group notification settings updated successfully",
	})
}

// ChangeSound 更换提示音
func ChangeSound(c *gin.Context) {
	var req struct {
		ID    string `json:"id"`    // 用户不可修改的 ID
		Sound string `json:"sound"` // 新的提示音路径
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	var setting models.SystemSetting
	if err := config.DB.Where("id = ?", req.ID).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "System settings not found",
		})
		return
	}

	setting.Sound = req.Sound
	if err := config.DB.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update sound",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sound updated successfully",
	})
}

// SubmitSound 上传新的提示音
func SubmitSound(c *gin.Context) {
	id := c.PostForm("id")
	file, err := c.FormFile("newsound") // 接收音频文件
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No file uploaded",
		})
		return
	}

	// 限制音频格式为mp3等
	//if !strings.HasSuffix(file.Filename, ".mp3") && !strings.HasSuffix(file.Filename, ".wav") {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"success": false,
	//		"message": "Invalid file format",
	//	})
	//	return
	//}

	// 保存音频文件
	//savePath := filepath.Join("uploads/sounds", file.Filename)
	savePath := fmt.Sprintf("uploads/sounds/%d_%s", time.Now().Unix(), file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save file",
		})
		return
	}

	// 更新数据库路径
	var setting models.SystemSetting
	if err := config.DB.Where("id = ?", id).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "System settings not found",
		})
		return
	}

	setting.Sound = savePath
	if err := config.DB.Save(&setting).Error; err != nil {
		// 如果数据库更新失败，删除已保存的文件
		_ = os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update sound path",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sound uploaded successfully",
		"path":    savePath, // 返回保存路径
	})
}

// ChangeBackground 更改聊天背景
func ChangeBackground(c *gin.Context) {
	id := c.PostForm("id")
	file, err := c.FormFile("background") // 接收背景图片文件
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No file uploaded",
		})
		return
	}

	// 校验文件格式，仅允许 .jpg、.png
	// 可补充：限制文件大小
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid file format. Only .jpg or .png allowed",
		})
		return
	}

	// 保存图片文件
	saveDir := "uploads/backgrounds"
	_ = os.MkdirAll(saveDir, os.ModePerm) // 确保目录存在
	savePath := filepath.Join(saveDir, fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename))
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save file",
		})
		return
	}

	// 更新数据库字段
	var setting models.SystemSetting
	if err := config.DB.Where("id = ?", id).First(&setting).Error; err != nil {
		// 删除已保存的文件
		_ = os.Remove(savePath)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "System settings not found",
		})
		return
	}

	setting.Background = savePath
	if err := config.DB.Save(&setting).Error; err != nil {
		// 删除已保存的文件
		_ = os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update background",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Background updated successfully",
	})
}

// GetSystemSetting 获取系统设置
func GetSystemSetting(c *gin.Context) {
	// 获取请求参数中的用户 ID
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID is required",
		})
		return
	}

	// 查询数据库获取系统设置
	var setting models.SystemSetting
	if err := config.DB.Where("id = ?", id).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "System settings not found",
		})
		return
	}

	// 返回系统设置
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "System settings retrieved successfully",
		"data": map[string]interface{}{
			"theme":       setting.Theme,
			"fontSize":    setting.FontSize,
			"fontStyle":   setting.FontStyle,
			"sound":       setting.Sound,
			"background":  setting.Background,
			"notice":      setting.Notice,
			"noticeGroup": setting.NoticeGroup,
		},
	})
}

// Logout 用户退出登录
func Logout(c *gin.Context) {
	var req struct {
		ID string `json:"id"` // 用户不可修改的 ID
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 查找账号信息
	var account models.AccountInfo
	if err := config.DB.Where("id = ?", req.ID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Account not found",
		})
		return
	}

	// 更新退出时间和状态
	now := time.Now()
	account.LastLogout = &now
	account.Status = "offline" // 标记用户为离线状态

	if err := config.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update logout information",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}
