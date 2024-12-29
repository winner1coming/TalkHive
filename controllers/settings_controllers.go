package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"TalkHive/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

//--------------------------------------------------------------------------
// 个人主业

// ShowProfile 返回个人信息
func ShowProfile(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err = global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	responseData := gin.H{
		"id":        user.ID,
		"nickname":  user.Nickname,
		"gender":    user.Gender,
		"birthday":  user.Birthday,
		"signature": user.Signature,
		"email":     user.Email,
		"phone":     user.Phone,
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询成功",
		"data":    responseData,
	})
}

// SaveEdit 保存编辑后的个人资料
func SaveEdit(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		Avatar    string `json:"avatar"`
		Nickname  string `json:"nickname"`
		Gender    string `json:"gender"`
		Birthday  string `json:"birthday"`
		Signature string `json:"signature"`
		Phone     string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求数据格式错误"})
		return
	}
	user.Avatar = input.Avatar
	user.Nickname = input.Nickname
	user.Gender = input.Gender
	user.Birthday = input.Birthday
	user.Signature = input.Signature
	user.Phone = input.Phone

	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "保存成功"})
}

//---------------------------------------------------------------------------------
// 安全设置

// GetUserInfo 展示安全信息
func GetUserInfo(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功获取用户信息",
		"data": gin.H{
			"id":                        user.ID,
			"email":                     user.Email,
			"password":                  user.Password,
			"friend_permissionID":       user.FriendPermissionID,
			"friend_permissionNickname": user.FriendPermissionNickName,
		},
	})
}

// GetCode 向id用户的新邮箱中发送验证码
func GetCode(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		NewEmail string `json:"new_email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的请求体"})
		return
	}
	fmt.Println("新邮箱", input.NewEmail)
	if !utils.ValidateEmail(input.NewEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮箱格式不正确"})
		return
	}

	// 新旧邮箱一致
	if input.NewEmail == user.Email {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "新邮箱不能与旧邮箱相同"})
		return
	}

	var existingUser models.AccountInfo
	err = global.Db.Where("email = ?", input.NewEmail).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮箱已被其他用户使用"})
		return
	}

	// 生成随机的验证码，并将验证码存储到 Redis 中，设置过期时间为 5 分钟
	code := utils.RandomCode(6)
	cacheKey := global.SmsCodeKey + input.NewEmail
	if err := global.RedisDB.Set(cacheKey, code, 5*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存验证码失败"})
		return
	}

	err = utils.SendSms(input.NewEmail, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "发送短信失败", "code": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "验证码已发送", "code": code})
}

// SaveEmail 修改账号绑定的邮箱
func SaveEmail(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		NewEmail string `json:"new_email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的请求体"})
		return
	}

	if !utils.ValidateEmail(input.NewEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮箱格式不正确"})
		return
	}

	// 新旧邮箱一致
	if input.NewEmail == user.Email {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "新邮箱不能与旧邮箱相同"})
		return
	}

	var existingUser models.AccountInfo
	err = global.Db.Where("email = ?", input.NewEmail).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮箱已被其他用户使用"})
		return
	}

	user.Email = input.NewEmail
	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存新邮箱失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "保存邮箱成功"})
	}
}

// SavePassword 保存修改的密码
func SavePassword(c *gin.Context) {
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "HTTP header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		NewPassword string `json:"newpassword"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的请求体"})
	}

	user.Password = input.NewPassword
	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改密码失败"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改密码成功", "password": user.Password})
	}
}

// IsIDAdd 查询id用户是否能被添加
func IsIDAdd(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}

	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		FriendPermissionID bool `json:"friend_permission_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的请求体"})
		return
	}

	user.FriendPermissionID = input.FriendPermissionID
	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改权限失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改权限成功"})
}

// IsNickNameAdd 查询昵称用户是否能被添加
func IsNickNameAdd(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		FriendPermissionNickName bool `json:"friend_permissionNickname"`
	}
	if c.ShouldBindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的请求体"})
		return
	}
	user.FriendPermissionNickName = input.FriendPermissionNickName
	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改权限失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改权限成功"})
	}
}

// ConfirmDeactivation 确认当前id用户是否注销
func ConfirmDeactivation(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	user.Deactivate = true
	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改权限失败"})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户注销成功"})
}

// --------------------------------------------------------------------------
// 系统设置

// ChangeTheme 更改系统主题颜色
func ChangeTheme(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// light dark system
	var input struct {
		Theme string `json:"theme"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数解析失败"})
		return
	}
	if input.Theme != "light" && input.Theme != "dark" && input.Theme != "system" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的主题值"})
		return
	}

	var systemSetting models.SystemSetting
	err = global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录不存在则创建新的新的系统设置，并且绑定id
			systemSetting = models.SystemSetting{
				AccountID: uint(accountID),
				Theme:     input.Theme,
			}
			if err := global.Db.Create(&systemSetting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "插入系统设置失败"})
				return
			}
		} else { // 正常情况下不会出现
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询系统设置失败（正常情况下不会出现）"})
			return
		}
	}

	systemSetting.Theme = input.Theme
	if err := global.Db.Save(&systemSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改主题失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改主题成功"})
	}
}

// ChangeFontsize 更改字体大小
func ChangeFontsize(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		FontSize uint `json:"font_size"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数解析失败"})
		return
	}

	var systemSetting models.SystemSetting
	if err = global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录不存在则创建新的新的系统设置，并且绑定id
			systemSetting = models.SystemSetting{
				AccountID: uint(accountID),
				FontSize:  input.FontSize,
			}
			if err := global.Db.Create(&systemSetting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "插入系统设置失败"})
				return
			}
		} else { // 正常情况下不会出现
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询系统设置失败（正常情况下不会出现）"})
			return
		}
	}
	systemSetting.FontSize = input.FontSize
	if err := global.Db.Save(&systemSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改字体大小失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改字体大小成功"})
	}
}

// ChangeFontstyle 更换字体风格
func ChangeFontstyle(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}
	var input struct {
		FontStyle string `json:"font_style"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数解析失败"})
		return
	}

	var systemSetting models.SystemSetting
	if err = global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录不存在则创建新的新的系统设置，并且绑定id
			systemSetting = models.SystemSetting{
				AccountID: uint(accountID),
				FontStyle: input.FontStyle,
			}
			if err := global.Db.Create(&systemSetting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "插入系统设置失败"})
				return
			}
		} else { // 正常情况下不会出现
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询系统设置失败（正常情况下不会出现）"})
			return
		}
	}
	systemSetting.FontStyle = input.FontStyle
	if err := global.Db.Save(&systemSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改字体风格失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改字体风格成功"})
	}
}

// IsNotice 是否有个人消息通知的声音
func IsNotice(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		Notice bool `json:"notice"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数解析失败"})
		return
	}
	var systemSetting models.SystemSetting
	if err = global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录不存在则创建新的新的系统设置，并且绑定id
			systemSetting = models.SystemSetting{
				AccountID: uint(accountID),
				Notice:    input.Notice,
			}
			if err := global.Db.Create(&systemSetting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "插入系统设置失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询系统设置失败（正常情况下不会出现）"})
			return
		}
	}
	systemSetting.Notice = input.Notice
	if err := global.Db.Save(&systemSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改通知设置失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改通知设置成功"})
	}
}

// IsNoticeGroup 是否有群消息通知的声音
func IsNoticeGroup(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		NoticeGroup bool `json:"notice_group"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数解析失败"})
		return
	}
	var systemSetting models.SystemSetting
	if err = global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录不存在则创建新的新的系统设置，并且绑定id
			systemSetting = models.SystemSetting{
				AccountID:   uint(accountID),
				NoticeGroup: input.NoticeGroup,
			}
			if err := global.Db.Create(&systemSetting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "插入系统设置失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询系统设置失败（正常情况下不会出现）"})
			return
		}
	}
	systemSetting.NoticeGroup = input.NoticeGroup
	if err := global.Db.Save(&systemSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改群消息通知声音设置失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改群消息通知声音设置成功"})
	}
}

// ChangeSound 更换提示音
func ChangeSound(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		Sound string `json:"sound"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数解析失败"})
		return
	}
	var systemSetting models.SystemSetting
	if err = global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录不存在则创建新的新的系统设置，并且绑定id
			systemSetting = models.SystemSetting{
				AccountID: uint(accountID),
				Sound:     input.Sound,
			}
			if err := global.Db.Create(&systemSetting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "插入新的系统设置失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询系统设置失败（正常情况下不会出现）"})
			return
		}
	}
	systemSetting.Sound = input.Sound
	if err := global.Db.Save(&systemSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改提示声音成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改提示声音成功"})
	}
}

// SubmitSound 上传新的提示音
func SubmitSound(c *gin.Context) {

}

// ChangeBackground 更换背景
func ChangeBackground(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var input struct {
		Background string `json:"back_ground"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "参数解析失败"})
		return
	}
	var systemSetting models.SystemSetting
	if err = global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果记录不存在则创建新的新的系统设置，并且绑定id
			systemSetting = models.SystemSetting{
				AccountID:  uint(accountID),
				Background: input.Background,
			}
			if err := global.Db.Create(&systemSetting).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "插入新的系统设置失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询系统设置失败（正常情况下不会出现）"})
			return
		}
	}
	systemSetting.Background = input.Background
	if err := global.Db.Save(&systemSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更改背景成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "更改背景成功"})
	}
}

// GetSystemSetting 返回系统设置
func GetSystemSetting(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	var systemSetting models.SystemSetting
	if err := global.Db.Where("account_id = ?", accountID).First(&systemSetting).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "系统设置未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "系统设置获取成功",
		"data": gin.H{
			"theme":       systemSetting.Theme,
			"fontSize":    systemSetting.FontSize,
			"fontStyle":   systemSetting.FontStyle,
			"sound":       systemSetting.Sound,
			"background":  systemSetting.Background,
			"notice":      systemSetting.Notice,
			"noticeGroup": systemSetting.NoticeGroup,
		},
	})
}

// Logout 退出登录
func Logout(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Header中的User-ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户未找到"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// 记录当前时间
	currentTime := time.Now().Format("2006-01-02 15:04:05") // 获取当前时间
	if err := global.Db.Model(&user).Update("last_logout", currentTime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新退出时间失败"})
		return
	}

	/* 继续完善JWT*/
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "退出登录成功"})
}
