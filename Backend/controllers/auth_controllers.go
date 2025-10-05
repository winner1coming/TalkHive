package controllers

import (
	"TalkHive/global"
	"TalkHive/models"
	"TalkHive/utils"
	"net/http"
	"strconv"
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
	// if account.Status == "online" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "该账号已登录，请勿重复登录"})
	// 	return
	// }

	// 生成 Token
	token, err := utils.GenerateJWT(account.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成 Token 失败"})
		return
	}

	//修改登录状态
	account.Status = "online"
	if err := global.Db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新用户状态失败"})
		return
	}

	avatarBase64, mimeType, err := utils.GetFileContentAndType(account.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	var links []models.Links
	if err := global.Db.Where("account_id = ?", account.AccountID).Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
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
		"links":    links,
	})
}

// Register 注册处理
func Register(c *gin.Context) {
	var input struct {
		Avatar   string `JSON:"avatar"`
		ID       string `JSON:"id"`
		Nickname string `JSON:"nickname"`
		Gender   string `JSON:"gender"`
		Birthday string `JSON:"birthday"`
		Email    string `JSON:"email"`
		Phone    string `JSON:"phone"`
		Password string `JSON:"password"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	// 校验字段
	if input.ID == "" || input.Nickname == "" || input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "数据不能为空"})
		return
	}

	// // 校验手机号格式是否正确
	// if !utils.ValidatePhone(input.Phone) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "电话号码格式不对"})
	// 	return
	// }

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

	// 创建用户
	newUser := models.AccountInfo{
		ID:                       input.ID,
		Nickname:                 input.Nickname,
		Gender:                   input.Gender,
		Birthday:                 input.Birthday,
		Email:                    input.Email,
		Phone:                    input.Phone,
		Password:                 input.Password,
		FriendPermissionID:       true,
		FriendPermissionNickName: true,
		Signature:                "这个人还没有自定义签名",
	}

	if err := global.Db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存用户信息失败"})
		return
	}

	if input.Avatar != "" {
		// 调用 GetAvatarPath 方法生成文件路径并保存图片
		avatarPath, err := utils.GetAvatarPath(input.Avatar, strconv.Itoa(int(newUser.AccountID)), "User_Avatar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		// 更新用户的 Avatar 字段
		newUser.Avatar = avatarPath
		if err := global.Db.Save(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新用户头像失败"})
			return
		}
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
	var err error
	if account.Avatar != "" {
		// 调用 GetFileContentAndType 方法获取文件内容和类型
		avatarBase64, mimeType, err = utils.GetFileContentAndType(account.Avatar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
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
