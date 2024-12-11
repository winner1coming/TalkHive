package controllers

import (
	"TalkHive/config"
	"TalkHive/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

// 资料编辑相关

// UpdateAvatar 更新用户头像
func UpdateAvatar(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Avatar    string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Update("avatar", req.Avatar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully"})
}

// UpdateNickname 修改用户昵称
func UpdateNickname(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Nickname  string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Update("nickname", req.Nickname).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nickname"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Nickname updated successfully"})
}

// EditProfile 编辑其他用户资料（性别、签名等）
func EditProfile(c *gin.Context) {
	var req struct {
		AccountID uint   `json:"account_id"`
		Gender    string `json:"gender"`
		Signature string `json:"signature"`
		ID        string `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.AccountInfo{}).Where("account_id = ?", req.AccountID).Updates(models.AccountInfo{
		Gender:    req.Gender,
		Signature: req.Signature,
		ID:        req.ID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// 账号设置相关

// SavePhone 保存新的手机号
func SavePhone(c *gin.Context) {
	var req struct {
		ID       string `json:"id"`       // 用户不可修改的id
		NewPhone string `json:"newphone"` // 新手机号
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 1. 检查新手机号是否已被注册
	var existingUser models.AccountInfo
	if err := config.DB.Where("phone = ?", req.NewPhone).First(&existingUser).Error; err == nil {
		// 如果手机号已注册，返回错误
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号已注册",
		})
		return
	}

	// 2. 更新手机号
	if err := config.DB.Model(&models.AccountInfo{}).Where("id = ?", req.ID).Update("phone", req.NewPhone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update phone",
		})
		return
	}

	// 3. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Phone updated successfully",
	})
}

// SavePassword 保存密码的更改
func SavePassword(c *gin.Context) {
	var req struct {
		ID          string `json:"id"`          // 用户不可修改的id
		NewPassword string `json:"newpassword"` // 新密码
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 1. 获取当前用户的账号信息
	var account models.AccountInfo
	if err := config.DB.Where("id = ?", req.ID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// 2. 更新密码
	account.Password = req.NewPassword // 新密码更新

	if err := config.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update password",
		})
		return
	}

	// 3. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password updated successfully",
	})
}

// checkPassword 验证密码是否正确
func checkPassword(storedPassword, providedPassword string) bool {
	// 假设密码是加密存储的，可以用哈希比较的方法
	// 比如 bcrypt.CompareHashAndPassword 等函数（具体实现取决于密码加密方式）
	// 如果密码没有加密存储，可以直接比较字符串
	return storedPassword == providedPassword
}

// IsIDAdd 设置是否允许通过ID查找
func IsIDAdd(c *gin.Context) {
	var req struct {
		IDStatus string `json:"idStatus"` // 是否允许通过ID查找，值为 "1" 或 "0"
		ID       string `json:"id"`       // 用户ID（不可修改的）
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 获取用户信息
	var account models.AccountInfo
	if err := config.DB.Where("id = ?", req.ID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// 更新friend_permission字段
	// 如果idStatus为"1"，表示允许通过ID查找；否则，表示不允许
	var newFriendPermission string
	if req.IDStatus == "1" {
		// 如果允许通过ID查找，存储id作为friend_permission
		newFriendPermission = req.ID
	} else {
		// 如果不允许通过ID查找，设置为""或其他标识
		newFriendPermission = ""
	}

	// 更新数据库中的friend_permission字段
	if err := config.DB.Model(&account).Update("friend_permission", newFriendPermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update friend permission",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Friend permission updated successfully",
	})
}

// IsPhoneAdd 设置是否允许通过手机号查找
func IsPhoneAdd(c *gin.Context) {
	var req struct {
		PhoneStatus string `json:"phoneStatus"` // 是否允许通过手机号查找，值为 "on" 或 "off"
		ID          string `json:"id"`          // 用户ID（不可修改的）
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 获取用户信息
	var account models.AccountInfo
	if err := config.DB.Where("id = ?", req.ID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// 更新friend_permission字段
	// 如果phoneStatus为"on"，表示允许通过手机号查找；否则，表示不允许
	var newFriendPermission string
	if req.PhoneStatus == "on" {
		// 如果允许通过手机号查找，存储手机号作为friend_permission
		newFriendPermission = account.Phone
	} else {
		// 如果不允许通过手机号查找，设置为""或其他标识
		newFriendPermission = ""
	}

	// 更新数据库中的friend_permission字段
	if err := config.DB.Model(&account).Update("friend_permission", newFriendPermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update friend permission",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Friend permission updated successfully",
	})
}

// ConfirmDeactivation 确认注销账号
func ConfirmDeactivation(c *gin.Context) {
	var req struct {
		ID string `json:"id"` // 用户ID，不能修改
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 查找用户
	var account models.AccountInfo
	if err := config.DB.Where("id = ?", req.ID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// 检查用户是否已注销
	if account.Deactivation {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User is already deactivated",
		})
		return
	}

	// 更新用户的注销状态
	account.Deactivation = true
	currentTime := time.Now()
	account.LastLogout = &currentTime // 记录注销时间

	if err := config.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to deactivate account",
		})
		return
	}

	// 处理其他相关逻辑，例如清除用户的消息等
	// 可以根据项目需求进一步清除或更新其他相关数据

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account deactivated successfully",
	})
}

// GetPhone 获取用户的手机号、密码和好友权限
func GetPhone(c *gin.Context) {
	// 获取请求参数 id
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID is required",
		})
		return
	}

	// 查询数据库，获取用户信息
	var account models.AccountInfo
	if err := config.DB.Where("account_id = ?", id).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// 返回手机号、密码和好友权限信息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data fetched successfully",
		"data": gin.H{
			"phone":             account.Phone,
			"password":          account.Password, // 注意：密码通常不直接返回，最好通过加密/哈希处理或通过特定机制传递
			"friend_permission": account.FriendPermission,
		},
	})
}

// GetCode 获取验证码
func GetCode(c *gin.Context) {
	// 获取请求参数
	var requestData struct {
		ID       string `json:"id"`
		NewPhone string `json:"newphone"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 1. 检查新手机号是否已被注册
	var existingUser models.AccountInfo
	if err := config.DB.Where("phone = ?", requestData.NewPhone).First(&existingUser).Error; err == nil {
		// 如果手机号已注册，返回错误
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号已注册",
		})
		return
	}

	// 2. 生成验证码逻辑（这里可以自定义验证码生成方法）
	code := generateVerificationCode() // 需要实现生成验证码的逻辑

	// 3. 发送验证码（此处只返回验证码示例，实际生产中需要通过短信平台发送）
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "验证码获取成功",
		"code":    code, // 返回验证码，实际应用中会发送短信给用户
	})
}

// generateVerificationCode 用于生成验证码的函数
func generateVerificationCode() string {
	// 这里简单生成一个6位数的验证码，实际可以根据需求修改
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// 系统设置相关

// UpdateSystemSettings 更新系统设置
func UpdateSystemSettings(c *gin.Context) {
	var req struct {
		AccountID  uint   `json:"account_id"`
		Background string `json:"background"`
		Theme      string `json:"theme"`
		Sound      string `json:"sound"`
		FontSize   int    `json:"font_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.SystemSetting{}).Where("account_id = ?", req.AccountID).Updates(models.SystemSetting{
		Background: req.Background,
		Theme:      req.Theme,
		Sound:      req.Sound,
		FontSize:   req.FontSize,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update system settings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "System settings updated successfully"})
}
