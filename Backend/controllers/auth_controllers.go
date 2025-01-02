package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"TalkHive/utils"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// Login 登录处理
func Login(c *gin.Context) {
	var input struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	// 解析 JSON 请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json输入格式错误"})
		return
	}

	// 查询数据库中账号信息
	var account models.AccountInfo
	if err := global.Db.Where("ID = ? AND password = ?", input.Account, input.Password).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "账号名称或密码错误"})
		return
	}

	// 生成 Token
	token, err := utils.GenerateJWT(account.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成 Token 失败",
		})
		return
	}

	//修改登录状态
	account.Status = "online"
	if err := global.Db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新用户状态失败"})
		return
	}

	var avatarBase64 string
	var mimeType string
	if account.Avatar != "" {
		avatarPath := account.Avatar
		file, err := os.Open(avatarPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "读取头像文件失败",
			})
			return
		}
		defer file.Close()

		fileContents, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "读取头像文件内容失败"})
			return
		}
		avatarBase64 = base64.StdEncoding.EncodeToString(fileContents)
		mimeType = mime.TypeByExtension(filepath.Ext(avatarPath))
		if mimeType == "" {
			mimeType = "image/jpg"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "登录成功",
		"nickname":   account.Nickname,
		"account_id": account.AccountID,
		"data": gin.H{
			"account": account.ID,
			"token":   token,
		},
		"email":    account.Email,
		"avatar":   avatarBase64,
		"mimeType": mimeType,
	})
}

// Register 注册处理
func Register(c *gin.Context) {
	var input struct {
		Avatar   *multipart.FileHeader `form:"avatar"`
		ID       string                `form:"id"`
		Nickname string                `form:"nickname"`
		Gender   string                `form:"gender"`
		Birthday string                `form:"birthday"`
		Email    string                `form:"email"`
		Phone    string                `form:"phone"`
		Password string                `form:"password"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	// 校验字段
	if input.ID == "" || input.Nickname == "" || input.Email == "" || input.Password == "" || input.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "数据不能为空"})
		return
	}

	// 校验手机号格式是否正确
	if !utils.ValidatePhone(input.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "电话号码格式不对"})
		return
	}

	// 检查账号ID是否已存在
	var existingUser models.AccountInfo
	if err := global.Db.Where("id = ?", input.ID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "当前账号ID已被使用"})
		return
	}

	// 检查邮箱是否已存在
	if err := global.Db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "邮箱已被使用"})
		return
	}

	// 保存头像文件到 images 文件夹
	var avatarPath string
	if input.Avatar != nil {
		// 构建目标文件夹路径//但是这个是在现有文件夹下的，感觉可以让用户选择一个文件夹，或者默认在C盘什么的
		baseDir := "D:/TalkHive" // 用户默认的 TalkHive 文件夹
		avatarDir := filepath.Join(baseDir, "User_Avatar")

		// 确保文件夹存在，如果不存在则创建
		if err := os.MkdirAll(avatarDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建文件夹失败"})
			return
		}

		// 构建文件路径
		avatarPath = filepath.Join(avatarDir, input.Avatar.Filename)

		// 保存文件
		if err := c.SaveUploadedFile(input.Avatar, avatarPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存头像失败"})
			return
		}
	}

	// 创建用户
	newUser := models.AccountInfo{
		Avatar:                   avatarPath,
		ID:                       input.ID,
		Nickname:                 input.Nickname,
		Gender:                   input.Gender,
		Birthday:                 input.Birthday,
		Email:                    input.Email,
		Phone:                    input.Phone,
		Password:                 input.Password,
		FriendPermissionID:       true,
		FriendPermissionNickName: true,
	}
	if err := global.Db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存用户信息失败"})
		return
	}

	// 创建系统设置
	newUserSetting := models.SystemSetting{
		AccountID:   newUser.AccountID,
		Background:  "",
		FontStyle:   "Microsoft YaHei",
		FontSize:    16,
		Theme:       "light",
		Sound:       "dingdo.mp3",
		Notice:      true,
		NoticeGroup: true,
	}
	if err := global.Db.Create(&newUserSetting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存用户设置失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "注册成功"})
}

// SendSmsCode 发送验证码
func SendSmsCode(c *gin.Context) {
	var input struct {
		Command string `json:"command"`
		Email   string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json输入格式有误"})
		return
	}

	if !utils.ValidateEmail(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮箱格式不正确"})
		return
	}

	switch input.Command {
	case "smsLogin":
		if !utils.CheckEmailRegistered(input.Email) { // 没有注册
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "该邮箱未注册"})
			return
		}
	case "register":
		if utils.CheckEmailRegistered(input.Email) { // 已经注册
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "该邮箱已注册，不可重复注册"})
			return
		}
	case "resetPassword":
		if !utils.CheckEmailRegistered(input.Email) { // 没有注册
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "该邮箱未注册"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的命令"})
		return
	}

	//生成验证码，并且往Redis中保存
	code := utils.RandomCode(6)
	cacheKey := global.SmsCodeKey + input.Email
	if err := global.RedisDB.Set(cacheKey, code, 5*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存验证码到Redis失败", "code": ""})
		return
	}

	//邮箱发送验证码
	err := utils.SendSms(input.Email, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "发送短信失败", "code": ""})
		return
	}

	// 返回结果
	var message string
	switch input.Command {
	case "smsLogin":
		message = "短信登录验证码发送成功"
	case "register":
		message = "短信验证码已发送，请查看您的邮箱"
	case "resetPassword":
		message = "重置密码的验证码已发送，请查看您的邮箱"
	}

	// 返回验证码到前端
	c.JSON(http.StatusOK, gin.H{"success": true, "message": message, "code": code})
}

// SmsLogin 短信登录
func SmsLogin(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无法解析Json"})
		return
	}

	var account models.AccountInfo
	if err := global.Db.Where("email = ?", input.Email).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户未找到"})
		return
	}

	if account.Deactivate == true {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "该账号已经注销"})
		return
	}

	var avatarBase64 string
	var mimeType string
	if account.Avatar != "" {
		avatarPath := account.Avatar
		file, err := os.Open(avatarPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "读取头像文件失败",
			})
			return
		}
		defer file.Close()

		fileContents, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "读取头像文件内容失败"})
			return
		}
		avatarBase64 = base64.StdEncoding.EncodeToString(fileContents)
		mimeType = mime.TypeByExtension(filepath.Ext(avatarPath))
		if mimeType == "" {
			mimeType = "image/jpg"
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "登录成功",
		"nickname":   account.Nickname,
		"account_id": account.AccountID,
		"avatar":     avatarBase64,
		"mimeType":   mimeType,
	})
}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 解析JSON请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	var account models.AccountInfo
	// 在数据库中查找手机号对应的账号
	if err := global.Db.Where("email = ?", input.Email).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该邮箱未注册"})
		return
	}

	// 更新数据库中的密码
	account.Password = string(input.Password)
	if err := global.Db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "密码更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "密码重置成功"})
}
