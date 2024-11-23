package controllers

import (
	"chatroom/config"
	"chatroom/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 用户注册
func Register(c *gin.Context) {
	var account models.AccountInfo
	// 绑定 JSON 数据到 account 对象
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名或邮箱是否已被注册
	var existingAccount models.AccountInfo
	if err := config.DB.Where("id = ? OR email = ?", account.ID, account.Email).First(&existingAccount).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ID or Email already exists"})
		return
	}

	// 创建新用户
	if err := config.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account registered successfully"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var loginRequest struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户并验证密码
	var account models.AccountInfo
	if err := config.DB.Where("id = ? AND password = ?", loginRequest.ID, loginRequest.Password).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"nickname": account.Nickname,
		"email":    account.Email,
	})
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

// AddFriend 添加好友
func AddFriend(c *gin.Context) {
	var request struct {
		UserID   uint   `json:"user_id"`
		FriendID uint   `json:"friend_id"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact := models.Contacts{
		ContactID:   request.FriendID,
		Remark:      request.Remark,
		IsGroupChat: false,
	}

	if err := config.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	var group models.GroupChatInfo
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group created successfully"})
}

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	var message struct {
		SenderID uint `json:"sender_id"`
		//ReceiverID uint   `json:"receiver_id"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := models.MessageInfo{
		SendAccountID: message.SenderID,
		// ReceiverID:    message.ReceiverID,
		Content: message.Content,
	}

	if err := config.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

// GetMessages 获取消息记录
func GetMessages(c *gin.Context) {
	receiverID := c.Param("id")

	var messages []models.MessageInfo
	if err := config.DB.Where("receiver_id = ?", receiverID).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// UpdateSystemSetting 更新系统设置
func UpdateSystemSetting(c *gin.Context) {
	var settings models.SystemSetting
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "System settings updated successfully"})
}

// SubmitApply 提交申请
func SubmitApply(c *gin.Context) {
	var applyInfo models.ApplyInfo
	if err := c.ShouldBindJSON(&applyInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&applyInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
}

// ReviewApply 审核申请
func ReviewApply(c *gin.Context) {
	var review struct {
		ApplyID uint   `json:"apply_id"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.ApplyInfo{}).Where("apply_id = ?", review.ApplyID).Update("status", review.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application reviewed successfully"})
}
