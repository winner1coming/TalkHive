package controllers

import (
	"TalkHive/config"
	"TalkHive/global"
	"TalkHive/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

// Register 用户注册
func Register(c *gin.Context) {
	var tempAccount struct {
		Avatar   string `json:"avatar"`
		ID       string `json:"id"`
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&tempAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json输入格式有误"})
		return
	}

	// 校验字段
	if tempAccount.ID == "" || tempAccount.Nickname == "" || tempAccount.Phone == "" || tempAccount.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "所有数据不能为空"})
		return
	}

	if global.Db == nil {
		fmt.Println("数据库连接失败：Register")
		return
	} else {
		fmt.Println("数据库连接成功：Register")
	}

	// 校验手机号格式是否正确（以 1 开头的 11 位数字）
	phoneRegex := regexp.MustCompile(`^1[0-9]{10}$`)
	if !phoneRegex.MatchString(tempAccount.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "手机号格式无效"})
		return
	}

	// 检查账号是否已存在
	var existingUser models.AccountInfo
	if err := global.Db.Where("id = ?", tempAccount.ID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "账号已被使用"})
		return
	}

	// 检查手机号是否已存在
	if err := global.Db.Where("phone = ?", tempAccount.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "手机号已被使用"})
		return
	}

	/*
		to be continued...
		哈希加密，需要使用到utils中的组件，后面再完善
	*/

	// 创建新用户
	newUser := models.AccountInfo{
		Avatar:   tempAccount.Avatar,
		ID:       tempAccount.ID,
		Nickname: tempAccount.Nickname,
		Phone:    tempAccount.Phone,
		Password: tempAccount.Password,
	}

	// 保存到数据库
	if err := global.Db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存用户信息失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "注册成功",
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var tempUser struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&tempUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json输入格式错误"})
		return
	}

	if global.Db == nil {
		fmt.Println("数据库连接失败：Login函数处失败")
		return
	} else {
		fmt.Println("数据库连接成功：Login函数处成功")
	}
	var account models.AccountInfo
	if err := global.Db.Where("phone = ? AND password = ?", tempUser.Phone, tempUser.Password).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "手机号或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"userphone": account.Phone,
	})
}

// 忘记密码（重置密码）
func ResetPassword(c *gin.Context) {
	var tempAccount struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	// 解析JSON请求体
	if err := c.ShouldBindJSON(&tempAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	fmt.Println("手机号:%s", tempAccount.Phone)
	// 校验手机号格式是否正确（以 1 开头的 11 位数字）
	phoneRegex := regexp.MustCompile(`^1[0-9]{10}$`)
	if !phoneRegex.MatchString(tempAccount.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "手机号格式无效"})
		return
	}

	// 验证密码长度
	if len(tempAccount.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "密码至少需要6个字符"})
		return
	}

	if global.Db == nil {
		fmt.Println("数据库连接失败：resetPassword函数处失败")
		return
	}
	var account models.AccountInfo
	// 在数据库中查找手机号对应的账号
	if err := global.Db.Where("phone = ?", tempAccount.Phone).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该手机号未注册"})
		return
	}

	// 更新数据库中的密码
	account.Password = string(tempAccount.Password)
	if err := global.Db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "密码更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "密码重置成功"})
}

// 发送短信验证码
func SendSmsCode(c *gin.Context) {

}

// 短信登录
func SmsLogin(c *gin.Context) {

}

// GetProfile 获取用户信息
func GetProfile(c *gin.Context) {
	id := c.Param("id")

	var account models.AccountInfo
	if err := config.DB.Where("id = ?", id).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_id": account.AccountID,
		"id":         account.ID,
		"nickname":   account.Nickname,
		"email":      account.Email,
		"avatar":     account.Avatar,
		"signature":  account.Signature,
		"gender":     account.Gender,
		"birthday":   account.Birthday,
	})
}

// UpdateProfile 更新用户信息
func UpdateProfile(c *gin.Context) {
	id := c.Param("id")

	var account models.AccountInfo
	if err := config.DB.Where("id = ?", id).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		Signature string `json:"signature"`
		Gender    string `json:"gender"`
		Birthday  string `json:"birthday"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新信息
	account.Nickname = updateData.Nickname
	account.Avatar = updateData.Avatar
	account.Signature = updateData.Signature
	account.Gender = updateData.Gender

	if err := config.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
