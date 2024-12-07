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

// Login 登录处理
func Login(c *gin.Context) {
	var input struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json输入格式错误"})
		return
	}

	var account models.AccountInfo
	if err := global.Db.Where("phone = ? AND password = ?", input.Phone, input.Password).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "手机号或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"userphone": account.Phone,
	})
}

// Register 注册处理
func Register(c *gin.Context) {
	var input struct {
		Avatar   string `json:"avatar"`
		ID       string `json:"id"`
		Nickname string `json:"nickname"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json输入格式有误"})
		return
	}

	// 校验字段
	if input.ID == "" || input.Nickname == "" || input.Phone == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "所有数据不能为空"})
		return
	}

	// 校验手机号格式是否正确
	if !utils.ValidatePhone(input.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "电话号码格式不对"})
		return
	}

	// 检查账号是否已存在
	var existingUser models.AccountInfo
	if err := global.Db.Where("id = ?", input.ID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "账号已被使用"})
		return
	}

	// 检查手机号是否已存在
	if err := global.Db.Where("phone = ?", input.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "手机号已被使用"})
		return
	}

	// 创建新用户
	newUser := models.AccountInfo{
		Avatar:   input.Avatar,
		ID:       input.ID,
		Nickname: input.Nickname,
		Phone:    input.Phone,
		Password: input.Password,
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

// SendSmsCode 发送验证码
func SendSmsCode(c *gin.Context) {
	var input struct {
		Phone string `json:"phone"`
	}

	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json输入格式有误"})
		return
	}

	// 检测手机号码格式是否正确
	if !utils.ValidatePhone(input.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "电话号码格式不对"})
		return
	}

	// 生成6位随机验证码
	code := utils.RandomCode(6)

	// 将验证码缓存到Redis以供验证使用
	cacheKey := global.SmsCodeKey + input.Phone
	if err := global.RedisDB.Set(cacheKey, code, 5*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save SMS code."})
		return
	}

	/*
		发送验证码到用户手机。。。
		展示还未能实现
	*/

	// 返回验证码到前端
	c.JSON(http.StatusOK, gin.H{"message": "成功发送验证码", "code": code})
}

// SmsLogin 短信登录
func SmsLogin(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从缓存中获取存储的验证码
	cacheKey := global.SmsCodeKey + input.Phone
	storedCode, err := global.RedisDB.Get(cacheKey).Result()
	if err != nil || storedCode != input.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "验证码失效或错误"})
		return
	}

	// 校验验证码成功后，查询数据库获取用户信息
	var account models.AccountInfo
	err = global.Db.Where("phone = ?", input.Phone).First(&account).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户未找到"})
		return
	}

	// 构造返回数据
	response := gin.H{
		"success":    true,
		"avatar":     account.Avatar,
		"nickname":   account.Nickname,
		"account_id": account.ID,
		"message":    "登录成功",
	}

	/*
		使用到JWT进行授权验证，后面再进行修改
	*/

	// 返回给前端
	c.JSON(http.StatusOK, response)
}

// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {
	var input struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	// 解析JSON请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	// 校验手机号格式是否正确
	if !utils.ValidatePhone(input.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "电话号码格式不对"})
		return
	}

	// 验证密码长度
	if len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "密码至少需要6个字符"})
		return
	}

	var account models.AccountInfo
	// 在数据库中查找手机号对应的账号
	if err := global.Db.Where("phone = ?", input.Phone).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该手机号未注册"})
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

// FriendRequests 获取当前用户的所有好友请求列表
func FriendRequests(c *gin.Context) {
	id := c.Param("id")

	// 获取所有当前用户收到的好友申请
	var receivedApplyInfos []models.ApplyInfo
	err := global.Db.Where("receiver_id = ? AND apply_type = ?", id, "friend").Find(&receivedApplyInfos).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "无法获取好友申请"})
		return
	}

	// 获取所有当前用户发出的好友申请
	var sentApplyInfos []models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND apply_type = ?", id, "friend").Find(&sentApplyInfos).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "无法获取发出的好友申请"})
		return
	}

	var friendRequests []map[string]interface{}

	// 处理接收到的申请
	for _, applyInfo := range receivedApplyInfos {
		var senderInfo models.AccountInfo
		err := global.Db.Where("account_id = ?", applyInfo.SenderID).First(&senderInfo).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "无法获取申请者信息"})
			return
		}

		friendRequest := map[string]interface{}{
			"avatar":     senderInfo.Avatar,
			"name":       senderInfo.Nickname,
			"account_id": senderInfo.AccountID,
			"reason":     applyInfo.Reason,
			"status":     applyInfo.Status,
			"time":       time.Now().Format("2006-01-02 15:04:05"),
			"type":       "received",
		}
		friendRequests = append(friendRequests, friendRequest)
	}

	// 处理发出的申请
	for _, applyInfo := range sentApplyInfos {
		var receiverInfo models.AccountInfo
		err := global.Db.Where("account_id = ?", applyInfo.ReceiverID).First(&receiverInfo).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "无法获取申请接收者信息"})
			return
		}

		friendRequest := map[string]interface{}{
			"avatar":     receiverInfo.Avatar,
			"name":       receiverInfo.Nickname,
			"account_id": receiverInfo.AccountID,
			"reason":     applyInfo.Reason,
			"status":     applyInfo.Status,
			"time":       time.Now().Format("2006-01-02 15:04:05"),
			"type":       "sent",
		}
		friendRequests = append(friendRequests, friendRequest)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "friend_requests": friendRequests})
}

// AcceptFriendRequest 同意好友申请
func AcceptFriendRequest(c *gin.Context) {
	var input struct {
		AccountID uint `json:"account_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误", "details": err.Error()})
		return
	}

	receiverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "接收者ID格式错误", "details": err.Error()})
		return
	}

	senderID := input.AccountID

	fmt.Printf("发送者id: %d 接收者id: %d\n", senderID, receiverID)

	var contactCount int64
	// 查询是否已经是好友
	if err := global.Db.Model(&models.Contacts{}).
		Where("(contact_id = ? OR contact_id = ?) AND is_group_chat = FALSE", receiverID, senderID).
		Count(&contactCount).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询好友关系时出错", "details": err.Error()})
		return
	}
	if contactCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "已是好友关系"})
		return
	}

	var applyInfo models.ApplyInfo
	if err := global.Db.Where("receiver_id = ? AND sender_id = ? AND status = ?", receiverID, senderID, "pending").First(&applyInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "好友申请未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库时出错", "details": err.Error()})
		}
		return
	}

	// 开启事务
	tx := global.Db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开启事务时出错"})
		return
	}

	// 更新申请状态
	applyInfo.Status = "accepted"
	if err := tx.Save(&applyInfo).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存请求时出错", "details": err.Error()})
		return
	}

	// 将双方添加为好友
	contacts := []models.Contacts{
		{ContactID: uint(receiverID), IsGroupChat: false},
		{ContactID: senderID, IsGroupChat: false},
	}

	for _, contact := range contacts {
		if err := tx.Model(&models.Contacts{}).Create(&contact).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加好友关系时出错", "details": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务时出错", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "好友申请已同意"})
}

// RejectFriendRequest 拒绝好友申请
func RejectFriendRequest(c *gin.Context) {
	var input struct {
		AccountID uint `json:"account_id"`
	}

	// 绑定并验证JSON输入
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误", "details": err.Error()})
		return
	}

	receiverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "接收者ID格式错误", "details": err.Error()})
		return
	}

	senderID := input.AccountID

	fmt.Printf("发送者id: %d 接收者id: %d\n", senderID, receiverID)

	var applyInfo models.ApplyInfo
	if err := global.Db.Where("receiver_id = ? AND sender_id = ? AND status = ?", receiverID, senderID, "pending").First(&applyInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "好友申请未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询数据库时出错", "details": err.Error()})
		}
		return
	}

	// 开启事务
	tx := global.Db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开启事务时出错"})
		return
	}

	// 更新申请状态
	applyInfo.Status = "rejected"
	if err := tx.Save(&applyInfo).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新申请状态时出错", "details": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务时出错", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "好友申请已拒绝"})
}

