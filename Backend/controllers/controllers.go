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

/*登录注册*/

// Login 登录处理
func Login(c *gin.Context) {
	var input struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	// 解析 JSON 请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Json输入格式错误",
		})
		return
	}

	// 查询数据库中账号信息
	var account models.AccountInfo
	if err := global.Db.Where("ID = ? AND password = ?", input.Account, input.Password).First(&account).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "账号名称或密码错误",
		})
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

	// 构建返回的 JSON 数据
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "登录成功",
		"avatar":     account.Avatar,
		"nickname":   account.Nickname,
		"account_id": account.ID,
		"data": gin.H{
			"account": account.ID,
			"token":   token,
		},
	})
}

// Register 注册处理
func Register(c *gin.Context) {
	var input struct {
		Avatar   string `json:"avatar"`
		ID       string `json:"id"`
		Nickname string `json:"nickname"`
		Gender   string `json:"gender"`
		Birthday string `json:"birthday"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
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

	// 创建新用户
	newUser := models.AccountInfo{
		Avatar:   input.Avatar,
		ID:       input.ID,
		Nickname: input.Nickname,
		Gender:   input.Gender,
		Birthday: input.Birthday,
		Email:    input.Email,
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
		if utils.CheckEmailRegistered(input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "该邮箱未注册", "code": ""})
			return
		}
	case "register":
		if utils.CheckEmailRegistered(input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "该邮箱已注册，不可重复注册", "code": ""})
			return
		}
	case "resetPassword":
		if utils.CheckEmailRegistered(input.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "该邮箱未注册", "code": ""})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的命令", "code": ""})
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

	response := gin.H{
		"success":    true,
		"avatar":     account.Avatar,
		"nickname":   account.Nickname,
		"account_id": account.ID,
		"message":    "登录成功",
	}
	c.JSON(http.StatusOK, response)
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

/*通信录*/

// FriendRequests 获取当前用户的所有好友请求列表
func FriendRequests(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少用户ID"})
		return
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID格式错误"})
		return
	}

	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "id用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "用户已经注销"})
		return
	}

	// 获取所有当前用户收到的好友申请
	var receivedApplyInfos []models.ApplyInfo
	err = global.Db.Where("receiver_id = ? AND apply_type = ?", accountID, "friend").Find(&receivedApplyInfos).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "无法获取好友申请"})
		return
	}

	// 获取所有当前用户发出的好友申请
	var sentApplyInfos []models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND apply_type = ?", accountID, "friend").Find(&sentApplyInfos).Error
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

// FriendRequestPend 处理好友请求
func FriendRequestPend(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求体格式错误"})
		return
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求体格式错误"})
		return
	}
	var input struct {
		AccountID uint `json:"account_id"`
		Accept    bool `json:"accept"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求体格式错误"})
		return
	}

	var me models.AccountInfo
	var friend models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.AccountID).First(&friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "申请好友的的用户不存在"})
		return
	}
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "当前id用户不存在"})
		return
	}

	if me.Deactivate == true {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "当前id用户已经注销"})
		return
	}
	if friend.Deactivate == true {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "申请好友的的用户已经注销"})
		return
	}

	var apply models.ApplyInfo
	if err := global.Db.Where("receiver_id = ? AND sender_id = ? AND apply_type = ? AND status = ?", accountID, input.AccountID, "friend", "pending").First(&apply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "申请不存在"})
		return
	}

	if input.Accept {
		apply.Status = "accepted"
	} else {
		apply.Status = "rejected"
	}
	if err := global.Db.Save(&apply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "无法更新申请状态"})
		return
	}
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

// GetGroupRequests 获取群组申请列表
func GetGroupRequests(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID为空"})
		return
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID格式错误"})
		return
	}

	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户已被注销"})
		return
	}

	// 获取当前用户的所有群聊申请（包括作为群主的和作为申请者的）
	var applyInfos []models.ApplyInfo
	err = global.Db.Where("(receiver_id = ? OR sender_id = ? ) AND apply_type = ?", accountID, accountID, "group_apply").Find(&applyInfos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库查询失败"})
		return
	}

	var groupRequests []map[string]interface{}
	for _, applyInfo := range applyInfos {
		var groupChat models.GroupChatInfo
		err := global.Db.Where("group_id = ?", applyInfo.GroupID).First(&groupChat).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "这个群没找到"})
		}

		var groupMember models.GroupMemberInfo
		err = global.Db.Where("account_id = ? AND group_id = ?", applyInfo.SenderID, applyInfo.GroupID).First(&groupMember).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve group member information"})
			return
		}

		// 获取申请人的信息（头像、昵称等）
		var accountInfo models.AccountInfo
		err = global.Db.Where("account_id = ?", applyInfo.SenderID).First(&accountInfo).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account information"})
			return
		}

		// 根据申请的状态生成群聊请求的响应数据
		groupRequest := map[string]interface{}{
			"avatar":        groupChat.GroupAvatar,                    // 群聊头像
			"group_name":    groupChat.GroupName,                      // 群名称
			"account_name":  accountInfo.Nickname,                     // 申请人或群主的昵称
			"other_side_id": applyInfo.SenderID,                       // 申请人或群主的ID
			"group_id":      applyInfo.GroupID,                        // 群聊ID
			"reason":        applyInfo.Reason,                         // 申请理由
			"apply_type":    applyInfo.ApplyType,                      // 申请类型
			"status":        applyInfo.Status,                         // 状态：pending、accepted、rejected等
			"time":          time.Now().Format("2006-01-02 15:04:05"), // 申请时间
		}

		groupRequests = append(groupRequests, groupRequest)
	}

	// 返回群聊请求列表
	c.JSON(http.StatusOK, gin.H{"群聊申请请求": groupRequests})
}

// DealGroupApplyRequest 处理群聊申请请求
func DealGroupApplyRequest(c *gin.Context) {
	userID := c.Param("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID无效"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"` // 对方的账号ID
		GroupID   uint `json:"group_id"`   // 群聊ID
		Accept    bool `json:"accept"`     // 是否接受邀请
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json解析失败"})
		return
	}

	// 检查申请记录是否存在，并获取相关信息
	var applyInfo models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND group_id = ? AND apply_type = ?", input.AccountID, userIDInt, input.GroupID, "group").First(&applyInfo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据库中没有这个群申请！"})
		return
	}

	// 根据是否接受邀请，更新申请表和群聊成员表
	if input.Accept {
		// 接受邀请，更新申请状态为 accepted
		applyInfo.Status = "accepted"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新申请表信息失败"})
			return
		}

		// 获取群聊名称作为昵称
		var groupChat models.GroupChatInfo
		err := global.Db.Where("group_id = ?", input.GroupID).First(&groupChat).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "群聊不存在"})
			return
		}

		// 将用户添加到群聊成员表中
		groupMember := models.GroupMemberInfo{
			AccountID:     uint(userIDInt),
			GroupID:       input.GroupID,
			GroupNickname: groupChat.GroupName, // 使用群聊名称作为群聊昵称
			IsBanned:      false,               // 初始加入群聊时不禁言
		}
		if err := global.Db.Create(&groupMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "将用户添加到群聊中失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "同意群聊加入申请"})
	} else {
		// 拒绝邀请，更新申请状态为 rejected
		applyInfo.Status = "rejected"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新申请表信息失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "拒绝群聊加入申请"})
	}
}

// DealGroupInviteRequest  处理群聊邀请请求
func DealGroupInviteRequest(c *gin.Context) {
	// 获取当前用户ID
	userID := c.Param("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID无效"})
		return
	}

	// 解析请求体中的数据
	var input struct {
		AccountID uint `json:"account_id"` // 申请人的账号 ID
		GroupID   uint `json:"group_id"`   // 群聊 ID
		Accept    bool `json:"accept"`     // 是否接受邀请
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json解析失败"})
		return
	}

	// 检查申请记录是否存在，并获取相关信息
	var applyInfo models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND group_id = ? AND apply_type = ?", input.AccountID, userIDInt, input.GroupID, "groupInvitation").First(&applyInfo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据库中没有这个群聊邀请申请记录！"})
		return
	}

	// 根据是否接受邀请，更新申请表和群聊成员表
	if input.Accept {
		// 接受邀请，更新申请状态为 accepted
		applyInfo.Status = "accepted"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新申请表信息失败"})
			return
		}

		// 获取群聊信息（群名称和群聊信息）
		var groupChat models.GroupChatInfo
		err := global.Db.Where("group_id = ?", input.GroupID).First(&groupChat).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "群聊不存在"})
			return
		}

		// 将用户添加到群聊成员表中
		groupMember := models.GroupMemberInfo{
			AccountID:     uint(userIDInt),
			GroupID:       input.GroupID,
			GroupNickname: groupChat.GroupName, // 使用群聊名称作为群聊昵称
			IsBanned:      false,               // 初始加入群聊时不禁言
		}
		if err := global.Db.Create(&groupMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "将用户添加到群聊中失败"})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "已接受群聊邀请"})
	} else {
		// 拒绝邀请，更新申请状态为 rejected
		applyInfo.Status = "rejected"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新申请表信息失败"})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "已拒绝群聊邀请"})
	}
}

// GetBlackList 获取当前用户的黑名单列表
func GetBlackList(c *gin.Context) {
	userID := c.Param("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id无效"})
		return
	}

	// 查询Contacts表，查找当前用户的黑名单联系人
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ? AND is_blacklist = ? AND is_group_chat = ?", userIDInt, true, false).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询黑名单失败"})
		return
	}

	// 获取黑名单用户的详细信息
	var blackList []gin.H
	for _, contact := range contacts {
		var account models.AccountInfo
		err := global.Db.Where("account_id = ?", contact.ContactID).First(&account).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户信息失败"})
			return
		}

		blackList = append(blackList, gin.H{
			"avatar":     account.Avatar,    // 头像
			"name":       account.Nickname,  // 昵称
			"account_id": account.AccountID, // 账户ID
			"signature":  account.Signature, // 电话

		})
	}

	c.JSON(http.StatusOK, gin.H{
		"blacklist": blackList,
	})
}

// RemoveFromBlacklist 从黑名单中移除用户
func RemoveFromBlacklist(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id不为空"})
		return
	}

	ownerID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}
	var input struct {
		AccountID uint `json:"account_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var contact models.Contacts
	result := global.Db.Where("owner_id = ? AND contact_id = ? AND is_group_chat = ?", ownerID, input.AccountID, false).First(&contact)

	// 如果没有找到该记录
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "没有找到这个关系"})
		return
	}

	// contact_id不在黑名单中
	if !contact.IsBlacklist {
		c.JSON(http.StatusOK, gin.H{"message": "黑名单已经不在黑名单中"})
		return
	}

	contact.IsBlacklist = false
	if err := global.Db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "黑名单移除错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "黑名单移除成功"})
}

// GetFriends 获取好友列表
func GetFriends(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id不为空"})
		return
	}

	ownerID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	// 定义好友返回结构体
	type FriendResponse struct {
		Avatar    string `json:"avatar"`
		AccountID uint   `json:"account_id"`
		Remark    string `json:"remark"`
		Status    string `json:"status"`
		Signature string `json:"signature"`
		Tag       string `json:"tag"`
	}

	// 查询联系人
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ?", ownerID).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取当前id的好友表失败"})
		return
	}

	var friends []FriendResponse
	for _, contact := range contacts {
		var accountInfo models.AccountInfo
		err := global.Db.First(&accountInfo, contact.ContactID).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取当前id的好友信息失败"})
		}
		friend := FriendResponse{
			Avatar:    accountInfo.Avatar,
			AccountID: accountInfo.AccountID,
			Remark:    contact.Remark,
			Status:    accountInfo.Status,
			Signature: accountInfo.Signature,
			Tag:       contact.Divide,
		}

		friends = append(friends, friend)
	}

	c.JSON(http.StatusOK, friends)
}

// GetGroups 获取群聊列表
func GetGroups(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id解析获取失败"})
		return
	}

	ownerID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	// 群聊返回结构体
	type GroupResponse struct {
		GroupID     uint   `json:"group_id"`
		GroupName   string `json:"group_name"`
		GroupAvatar string `json:"group_avatar"`
	}

	// 查询群聊列表
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ? AND is_group_chat = ?", ownerID, true).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching groups"})
		return
	}

	var groups []GroupResponse
	for _, contact := range contacts {
		var groupInfo models.GroupChatInfo
		err := global.Db.Model(&models.GroupChatInfo{}).First(&groupInfo, contact.ContactID).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "当前群号的群信息查找失败"})
		}

		// 组装返回数据
		group := GroupResponse{
			GroupID:     groupInfo.GroupID,
			GroupName:   groupInfo.GroupName,
			GroupAvatar: groupInfo.GroupAvatar,
		}
		groups = append(groups, group)
	}

	c.JSON(http.StatusOK, groups)
}

// CreateGroup 当前id用户创建群聊
func CreateGroup(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id解析获取失败"})
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
	}

	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户已注销"})
		return
	}

	var input struct {
		GroupName   string `json:"group_name" binding:"required"`
		GroupAvatar string `json:"group_avatar" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var group models.GroupChatInfo
	group.GroupName = input.GroupName
	group.GroupAvatar = input.GroupAvatar
	group.GroupOwner = uint(accountID)

	if err := global.Db.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"GroupID": group.GroupID, "message": "创建群聊失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "创建群聊成功"})
	}

}

// DeleteGroup 当前ID用户删除群聊
func DeleteGroup(c *gin.Context) {

}

// ProfileCard 获取id用户的获取资料卡片
func ProfileCard(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id为空"})
		return
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	//当前用户
	var owner models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&owner).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	if owner.Deactivate == true {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户已注销"})
	}

	var input struct {
		AnotherID int `json:"another_id" binding:"required"`
	}

	//需要查询的用户
	var another models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.AnotherID).First(&another).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
	}
	if another.Deactivate == true {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户已注销"})
	}

	var contact models.Contacts
	if err := global.Db.Where("owner_id = ? and contact_id = ?", owner.AccountID, another.AccountID).First(&contact).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "不是好友关系"})
		return
	}

	responseData := gin.H{
		"account_id": another.AccountID,
		"avatar":     another.Avatar,
		"remark":     contact.Remark,
		"divide":     contact.Divide,
		"nickname":   another.Nickname,
		"signature":  another.Signature,
	}
	c.JSON(http.StatusOK, responseData)
}

// ShowProfile 返回个人信息
func ShowProfile(c *gin.Context) {

	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少账号ID",
		})
		return
	}

	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	// 查询 AccountInfo 表，查找对应的用户信息
	var account models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "账号未找到",
		})
		return
	}

	responseData := gin.H{
		"id":        account.ID,
		"nickname":  account.Nickname,
		"gender":    account.Gender,
		"birthday":  account.Birthday,
		"signature": account.Signature,
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
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少账号ID",
		})
		return
	}

	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	var input struct {
		ID        string `json:"id"`
		Avatar    string `json:"avatar"`
		Nickname  string `json:"nickname"`
		Gender    string `json:"gender"`
		Birthday  string `json:"birthday"`
		Signature string `json:"signature"`
	}

	// 解析请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	// 查询数据库中的用户数据
	var account models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&account).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "账号未找到",
		})
		return
	}

	// 更新用户资料
	account.ID = input.ID
	account.Avatar = input.Avatar
	account.Nickname = input.Nickname
	account.Gender = input.Gender
	account.Birthday = input.Birthday
	account.Signature = input.Signature

	// 保存更新的数据
	if err := global.Db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "保存成功",
	})
}

// GetPhone 获取用户的手机号
func GetPhone(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少账号ID",
		})
		return
	}

	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	// 查询数据库中的用户信息
	var userAccount models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&userAccount).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "账号未找到",
		})
		return
	}

	// 返回用户的手机号、密码和权限信息
	/*
		to be continued 此处的密码可能需要加密解密处理后
	*/

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取手机号成功",
		"data": gin.H{
			"phone":             userAccount.Phone,
			"password":          userAccount.Password,
			"friend_permission": userAccount.FriendPermission,
		},
	})
}

// SavePhone 修改手机号
func SavePhone(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少账号ID",
		})
		return
	}

	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	// 获取请求体中的新手机号
	var input struct {
		NewPhone string `json:"newphone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求体",
		})
		return
	}

	// 校验手机号格式
	if !utils.ValidatePhone(input.NewPhone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "手机号格式不正确",
		})
		return
	}

	// 检查手机号是否已被其他用户使用
	var existingUser models.AccountInfo
	err = global.Db.Where("phone = ?", input.NewPhone).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "该手机号已被其他用户使用",
		})
		return
	}

	// 查询当前用户
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户未找到",
		})
		return
	}

	// 更新手机号
	user.Phone = input.NewPhone
	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新手机号失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "手机号更新成功",
		"data": gin.H{
			"newphone": user.Phone,
		},
	})
}

// GetCode 向id用户的新手机号中发送验证码
func GetCode(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少账号ID",
		})
		return
	}

	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Json中ID解析失败"})
		return
	}

	// 当前用户是否存在
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户未找到",
		})
		return
	}

	// 获取请求体中的新手机号
	var input struct {
		NewPhone string `json:"newphone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求体",
		})
		return
	}

	// 校验手机号格式
	if !utils.ValidatePhone(input.NewPhone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "手机号格式不正确",
		})
		return
	}

	// 检查手机号是否已被其他用户使用
	var existingUser models.AccountInfo
	err = global.Db.Where("phone = ?", input.NewPhone).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "该手机号已被其他用户使用",
		})
		return
	}

	// 生成6位随机验证码
	code := utils.RandomCode(6)

	// 将验证码存储到 Redis 中，设置过期时间为 5 分钟
	cacheKey := global.SmsCodeKey + input.NewPhone
	if err := global.RedisDB.Set(cacheKey, code, 5*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存验证码失败",
		})
		return
	}

	/*
		To be continued 向用户id手机号中发送验证码
	*/
	// 返回发送验证码成功的响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "验证码已发送",
		"code":    code,
	})
}

// SavePassword 保存修改的密码
func SavePassword(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "URL中id为空",
		})
		return
	}

	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID解析失败"})
		return
	}

	// 当前用户是否存在
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户未找到",
		})
		return
	}

	var input struct {
		NewPassword string `json:"newpassword"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的请求体"})
	}

	// 更改
	user.Password = input.NewPassword
	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更改密码失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "更改密码成功",
		"password": user.Password,
	})

}

// IsIDAdd 查询id用户是否能被添加
func IsIDAdd(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "URL中id为空"})
		return
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户未找到"})
		return
	}

	if user.FriendPermission == "ID_true" {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "可以添加"})
	} else if user.FriendPermission == "ID_false" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "不能添加"})
	}
	return
}

// IsPhoneAdd 当前id用户是否能通过手机号进行添加好友
func IsPhoneAdd(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "URL中id为空"})
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID解析失败"})
	}
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户未找到"})
	}

	if user.FriendPermission == "Phone_true" {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "可以添加"})
	} else if user.FriendPermission == "Phone_false" {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "不能添加"})
	}
}

// ConfirmDeactivation 确认当前id用户是否注销
func ConfirmDeactivation(c *gin.Context) {
	urlID := c.Param("id")
	if urlID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "URL中id为空"})
	}
	accountID, err := strconv.Atoi(urlID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID解析失败"})
	}
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户未找到"})
	}

	if user.Deactivate == true {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户已经注销"})
		return
	}

	user.Deactivate = true
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户注销成功"})
}
