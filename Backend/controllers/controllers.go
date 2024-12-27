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

// -------------------------------------------------------------------------------
/*登录注册*/

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

	// 构建返回的 JSON 数据
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "登录成功",
		"avatar":     account.Avatar,
		"nickname":   account.Nickname,
		"account_id": account.AccountID,
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

// ---------------------------------------------------------------------------
/*通信录*/

// SearchStrangers 搜索添加陌生人
func SearchStrangers(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中缺少用户ID"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
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
	var input struct {
		Key string `json:"key"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少关键字参数或关键字为空"})
		return
	}

	var strangers []gin.H
	var accounts []models.AccountInfo
	err = global.Db.Where("id LIKE ? OR phone LIKE ? OR email LIKE ? OR account_id LIKE ?", "%"+input.Key+"%", "%"+input.Key+"%", "%"+input.Key+"%", "%"+input.Key+"%").Find(&accounts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户信息失败"})
		return
	}
	for _, account := range accounts {
		strangers = append(strangers, gin.H{
			"id":       account.AccountID,
			"nickname": account.Nickname,
			"avatar":   account.Avatar,
		})
	}

	var groupchats []models.GroupChatInfo
	err = global.Db.Where("contact_id LIKE ?", "%"+input.Key+"%").Find(&groupchats).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊信息失败"})
		return
	}
	for _, groupchat := range groupchats {
		strangers = append(strangers, gin.H{
			"group_id": groupchat.GroupID,
			"nickname": groupchat.GroupName,
			"avatar":   groupchat.GroupAvatar,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "查询成功", "data": strangers})
}

//-----------------------------------------------------------------------------

/*好友请求*/

// GetFriendRequests 获取当前用户的所有好友请求列表
func GetFriendRequests(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中缺少用户ID"})
		return
	}
	fmt.Println(ID)
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取其他人加当前id用户的好友申请失败"})
		return
	}

	// 获取所有当前用户发出的好友申请
	var sentApplyInfos []models.ApplyInfo
	err = global.Db.Where("(sender_id = ? AND apply_type = ?", accountID, "friend").Find(&sentApplyInfos).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取当前id用户加其他人的好友申请失败"})
		return
	}

	var friendRequests []map[string]interface{}

	// 其他人加当前用户id
	for _, applyInfo := range receivedApplyInfos {
		var senderInfo models.AccountInfo
		err := global.Db.Where("account_id = ?", applyInfo.SenderID).First(&senderInfo).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "无法获取申请者信息"})
			return
		}

		friendRequest := map[string]interface{}{
			"apply_id":    applyInfo.ApplyID,
			"avatar":      senderInfo.Avatar,
			"name":        senderInfo.Nickname,
			"sender_id":   applyInfo.SenderID,
			"receiver_id": applyInfo.ReceiverID,
			"reason":      applyInfo.Reason,
			"status":      applyInfo.Status,
			"time":        time.Now().Format("2006-01-02 15:04:05"),
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
			"apply_id":    applyInfo.ApplyID,
			"avatar":      receiverInfo.Avatar,
			"name":        receiverInfo.Nickname,
			"sender_id":   applyInfo.SenderID,
			"receiver_id": applyInfo.ReceiverID,
			"reason":      applyInfo.Reason,
			"status":      applyInfo.Status,
			"time":        time.Now().Format("2006-01-02 15:04:05"),
		}
		friendRequests = append(friendRequests, friendRequest)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功！", "friend_requests": friendRequests})
}

// FriendRequestPend 处理好友请求(其他人申请当前id用户为好友)
func FriendRequestPend(c *gin.Context) {
	// 从请求头获取用户ID
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	// 转换用户ID为整数
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
		Accept    bool `json:"accept"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求体格式错误"})
		return
	}

	// 判断当前id用户和申请的用户是否存在
	var me models.AccountInfo
	var friend models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "当前id用户不存在"})
		return
	}
	if err := global.Db.Where("account_id = ?", input.AccountID).First(&friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "申请好友的的用户不存在"})
		return
	}

	//判断是否注销
	if me.Deactivate == true {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "当前id用户已经注销"})
		return
	}
	if friend.Deactivate == true {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "申请好友的的用户已经注销"})
		return
	}

	// 判断是否存在好友关系
	var apply models.ApplyInfo
	if err := global.Db.Where("receiver_id = ? AND sender_id = ? AND apply_type = ? AND status = ?", accountID, input.AccountID, "friend", "pending").First(&apply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "申请不存在"})
		return
	}

	//处理申请
	if input.Accept {
		apply.Status = "accepted"
	} else {
		apply.Status = "rejected"
	}

	// 保存到数据库
	if err := global.Db.Save(&apply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "无法更新申请状态"})
		return
	}
}

// AddFriend 添加好友
func AddFriend(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var me models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&me).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户不存在"})
		return
	}
	if me.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户已注销，无法进行操作"})
		return
	}
	var input struct {
		AccountID uint   `json:"account_id"` // 被申请人ID
		Reason    string `json:"reason"`     // 申请理由
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}
	var other models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.AccountID).First(&other).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不存在"})
		return
	}
	if other.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户已注销，无法添加好友"})
		return
	}

	// 检查是否已经是好友
	var contacts models.Contacts
	err = global.Db.Where("(owner_id = ? AND contact_id = ?) OR (owner_id = ? AND contact_id = ?)", me.AccountID, other.AccountID, other.AccountID, me.AccountID).First(&contacts).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "你们已经是好友，请勿重复添加"})
		return
	}

	// 检查是否已发送过好友申请
	var existingApply models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND apply_type = ?", me.AccountID, other.AccountID, "friend").First(&existingApply).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "好友申请已发送，请勿重复申请"})
		return
	}

	newApply := models.ApplyInfo{
		ApplyType:  "friend",        // 申请类型
		SenderID:   me.AccountID,    // 发送者ID
		ReceiverID: other.AccountID, // 接收者ID
		Status:     "pending",       // 初始状态
		Reason:     input.Reason,    // 申请理由
	}

	if err := global.Db.Create(&newApply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "好友申请保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "好友申请已发送", "apply_id": newApply.ApplyID})
}

// ---------------------------------------------------------------------------
/*群聊请求*/

// DealGroupInviteRequest  处理群聊邀请请求
func DealGroupInviteRequest(c *gin.Context) {
	// 从请求头获取用户ID
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	// 转换用户ID为整数
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "json解析失败"})
		return
	}

	// 数据库中是否有这个群聊邀请
	var applyInfo models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND group_id = ? AND apply_type = ?", input.AccountID, accountID, input.GroupID, "groupInvitation").First(&applyInfo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "数据库中没有这个群聊邀请申请记录！"})
		return
	}

	if input.Accept { // 接受邀请，更新申请状态为 accepted
		applyInfo.Status = "accepted"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新申请表信息失败"})
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
			AccountID:     uint(accountID),
			GroupID:       input.GroupID,
			GroupNickname: groupChat.GroupName, // 使用群聊名称作为群聊昵称
			IsBanned:      false,               // 初始加入群聊时不禁言
		}
		if err := global.Db.Create(&groupMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "将用户添加到群聊中失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "已接受群聊邀请"})
	} else { // 拒绝邀请，更新申请状态为 rejected
		applyInfo.Status = "rejected"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新申请表信息失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "已拒绝群聊邀请"})
	}
}

// DealGroupApplyRequest 处理群聊申请请求
func DealGroupApplyRequest(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}
	AccountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID转换失败"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
		GroupID   uint `json:"group_id"`
		Accept    bool `json:"accept"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json解析失败"})
		return
	}

	// 检查申请记录是否存在
	var applyInfo models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND receiver_id = ? AND group_id = ? AND apply_type = ?", input.AccountID, AccountID, input.GroupID, "groupApply").First(&applyInfo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "数据库中没有这个群申请！"})
		return
	}

	if input.Accept { // 接受邀请
		applyInfo.Status = "accepted"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新申请表信息失败"})
			return
		}

		// 查询该群聊
		var groupChat models.GroupChatInfo
		err := global.Db.Where("group_id = ?", input.GroupID).First(&groupChat).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
			return
		}

		// 将用户添加到群聊成员表中
		groupMember := models.GroupMemberInfo{
			AccountID:     uint(AccountID),
			GroupID:       input.GroupID,
			GroupNickname: groupChat.GroupName, // 使用群聊名称作为群聊昵称
			IsBanned:      false,               // 初始加入群聊时不禁言
		}
		if err := global.Db.Create(&groupMember).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "将用户添加到群聊中失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "同意群聊加入申请"})
	} else { // 拒绝邀请
		applyInfo.Status = "rejected"
		if err := global.Db.Save(&applyInfo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "更新申请表信息失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "拒绝群聊加入申请"})
	}
}

// GetGroupRequests 获取群组申请列表
func GetGroupRequests(c *gin.Context) {
	// 从请求头获取用户ID
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	// 转换用户ID为整数
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}

	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id用户已被注销"})
		return
	}

	// 获取当前用户的所有群聊申请（包括作为群主的和作为申请者的，以及申请类型的和邀请类型的）
	var applyInfos []models.ApplyInfo
	err = global.Db.Where("(receiver_id = ? OR sender_id = ? ) AND (apply_type = ? OR apply_type = ?)", accountID, accountID, "groupInvitation", "groupApply").Find(&applyInfos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "数据库查询失败"})
		return
	}

	var groupRequests []map[string]interface{}
	for _, applyInfo := range applyInfos {
		var groupChat models.GroupChatInfo
		// 判断此群是否存在
		err := global.Db.Where("group_id = ?", applyInfo.GroupID).First(&groupChat).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "这个群没找到"})
			return
		}

		var groupMember models.GroupMemberInfo
		// 申请人SenderID是否已经在群聊中
		err = global.Db.Where("account_id = ? AND group_id = ?", applyInfo.SenderID, applyInfo.GroupID).First(&groupMember).Error
		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "申请人已经在群聊中"})
			return
		}

		// 获取申请人的信息
		var accountInfo models.AccountInfo
		err = global.Db.Where("account_id = ?", applyInfo.SenderID).First(&accountInfo).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取申请人信息失败"})
			return
		}
		if accountInfo.Deactivate == true {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "申请人已被注销"})
			return
		}

		groupRequest := map[string]interface{}{
			"apply_id":     applyInfo.ApplyID,                        // 申请ID
			"avatar":       groupChat.GroupAvatar,                    // 群聊头像
			"group_name":   groupChat.GroupName,                      // 群名称
			"account_name": accountInfo.Nickname,                     // 申请人或群主的昵称
			"sender_id":    applyInfo.SenderID,                       // 申请人或群主的ID
			"receiver_id":  applyInfo.ReceiverID,                     // 接收者ID
			"group_id":     applyInfo.GroupID,                        // 群聊ID
			"reason":       applyInfo.Reason,                         // 申请理由
			"apply_type":   applyInfo.ApplyType,                      // 申请类型
			"status":       applyInfo.Status,                         // 状态：pending、accepted、rejected等
			"time":         time.Now().Format("2006-01-02 15:04:05"), // 申请时间
		}

		groupRequests = append(groupRequests, groupRequest)
	}

	c.JSON(http.StatusOK, gin.H{"success": false, "message": "成功!", "群聊申请请求": groupRequests})
}

// AddGroup 申请加入群聊
func AddGroup(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户不存在"})
		return
	}
	if user.Deactivate {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "当前用户已注销，无法进行操作"})
		return
	}
	var input struct {
		GroupID uint   `json:"group_id"` // 群组ID
		Reason  string `json:"reason"`   // 申请理由
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查用户是否已经是群成员
	var membership models.GroupMemberInfo
	err = global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, user.AccountID).First(&membership).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "你已经是该群的成员，请勿重复申请"})
		return
	}

	// 检查是否已发送过加入申请
	var existingApply models.ApplyInfo
	err = global.Db.Where("sender_id = ? AND group_id = ? AND apply_type = ?", user.AccountID, input.GroupID, "groupInvitation").First(&existingApply).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "加入群组申请已发送，请勿重复申请"})
		return
	}

	// 构造新的群组申请记录
	newApply := models.ApplyInfo{
		ApplyType:  "groupInvitation", // 申请类型
		SenderID:   user.AccountID,    // 申请者ID
		ReceiverID: group.GroupOwner,  // 群主作为接收方
		Status:     "pending",         // 初始状态
		Reason:     input.Reason,      // 申请理由
	}

	if err := global.Db.Create(&newApply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群组申请保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "加入群组申请已发送", "apply_id": newApply.ApplyID})
}

// ---------------------------------------------------------------------------
/*黑名单*/

// GetBlackList 获取当前用户的黑名单列表
func GetBlackList(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	// 转换用户ID为整数
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}

	// 查询Contacts表，查找当前用户的黑名单联系人
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ? AND is_blacklist = ? AND is_group_chat = ?", accountID, true, false).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询黑名单失败"})
		return
	}

	// 获取黑名单用户的详细信息
	var blackList []gin.H
	for _, contact := range contacts {
		var account models.AccountInfo
		err := global.Db.Where("account_id = ?", contact.ContactID).First(&account).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "查询用户信息失败"})
			return
		}

		blackList = append(blackList, gin.H{
			"avatar":     account.Avatar,
			"name":       account.Nickname,
			"account_id": account.AccountID,
			"signature":  account.Signature,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"blacklist": blackList,
	})
}

// RemoveFromBlacklist 从黑名单中移除用户
func RemoveFromBlacklist(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID为空"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	var contact models.Contacts
	result := global.Db.Where("owner_id = ? AND contact_id = ? AND is_group_chat = ?", accountID, input.AccountID, false).First(&contact)

	// 如果没有找到该记录
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "没有找到这个关系"})
		return
	}

	// contact_id不在黑名单中
	if !contact.IsBlacklist {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "黑名单已经不在黑名单中"})
		return
	}

	contact.IsBlacklist = false
	if err := global.Db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "黑名单移除错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "黑名单移除成功"})
}

// AddToBlacklist 将用户添加到黑名单
func AddToBlacklist(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	AccountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "转换失败"})
		return
	}

	var input struct {
		AccountID uint `json:"account_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Json解析失败"})
		return
	}

	var contact models.Contacts
	result := global.Db.Where("owner_id = ? AND contact_id = ? AND is_group_chat = ?", AccountID, input.AccountID, false).First(&contact)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "没有找到这个关系"})
		return
	}
	if contact.IsBlacklist {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "黑名单已经存在"})
		return
	}

	contact.IsBlacklist = true
	if err := global.Db.Save(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加黑名单错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "添加黑名单成功"})
}

// ---------------------------------------------------------------------------
/*好友列表*/

// GetFriends 获取好友列表
func GetFriends(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
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
	err = global.Db.Where("owner_id = ?", accountID).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取当前ID的好友列表失败"})
		return
	}

	var friends []FriendResponse
	for _, contact := range contacts {
		if contact.IsBlacklist { // 黑名单不显示
			continue
		}
		var accountInfo models.AccountInfo
		err := global.Db.First(&accountInfo, contact.ContactID).Error
		if err != nil { // 如果某个联系人信息获取失败，继续下一个联系人
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取当前好友信息失败"})
			continue
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取好友列表成功", "data": friends})
}

// ---------------------------------------------------------------------------
/*分组部分*/

// GetDivides 获取分组名称
func GetDivides(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}

	// 验证用户是否存在
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// 获取URL中的type参数
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" { // 参数不匹配
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 当前id用户创建分组名称
	type response struct {
		Divide string `json:"divide"` // 分组
	}

	var isGroupChat bool
	if groupType == "groups" {
		isGroupChat = true // 群聊分组
	} else {
		isGroupChat = false // 好友分组
	}

	// 查询用户分组信息
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ? AND is_group_chat = ?", accountID, isGroupChat).Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取当前ID的群聊/好友分组失败"})
		return
	}

	// 存储去重后的分组信息
	var groups []response
	seenGroups := make(map[string]bool) // 用来去重的map

	for _, contact := range contacts {
		if _, exists := seenGroups[contact.Divide]; exists { // 分组去重
			continue
		}

		seenGroups[contact.Divide] = true
		group := response{
			Divide: contact.Divide,
		}
		groups = append(groups, group)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功", "groups": groups})
}

// CreateDivide 创建分组
func CreateDivide(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}

	// 验证用户是否存在
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	// 获取URL中的type参数
	groupType := c.Param("type") // 'friends' 或 'groups'
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 分组名称
	var input struct {
		Divide string `json:"divide" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "json绑定失败"})
		return
	}

	// 定义是否是群聊
	var isGroupChat bool
	if groupType == "groups" {
		isGroupChat = true // 创建群聊分组
	} else {
		isGroupChat = false // 创建好友分组
	}

	// 检查分组名称是否已存在
	var existingContact models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, input.Divide, isGroupChat).First(&existingContact).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "分组名称已存在"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询分组名称失败"})
		return
	}

	// 创建新分组
	newContact := models.Contacts{
		OwnerID:     uint(accountID),
		Divide:      input.Divide,
		IsGroupChat: isGroupChat,
	}
	if err := global.Db.Create(&newContact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建分组失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建分组成功"})
}

// DeleteDivide 删除分组
func DeleteDivide(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	err = global.Db.Where("account_id = ?", accountID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户已注销"})
		return
	}

	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 获取URL中的divide参数（要删除的分组名称）
	divideName := c.Param("divide")
	if divideName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "分组名称不能为空"})
		return
	}

	// 定义是否是群聊
	var isGroupChat bool
	if groupType == "groups" {
		isGroupChat = true
	} else {
		isGroupChat = false
	}

	// 检查分组是否存在
	var group models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, divideName, isGroupChat).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "分组不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询分组失败"})
		}
		return
	}

	// 将分组内的成员移入未分类分组（"未分类" 分组名称）
	uncategorizedDivide := "未分类"
	var uncategorizedGroup models.Contacts

	// 检查是否已经存在 "未分类" 分组
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, uncategorizedDivide, isGroupChat).First(&uncategorizedGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果 "未分类" 分组不存在，则创建一个
			uncategorizedGroup = models.Contacts{
				OwnerID:     uint(accountID),
				Divide:      uncategorizedDivide,
				IsGroupChat: isGroupChat,
			}
			if err := global.Db.Create(&uncategorizedGroup).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建未分类分组失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询未分类分组失败"})
			return
		}
	}

	// 将分组内的所有联系人或成员移动到 "未分类" 分组
	err = global.Db.Model(&models.Contacts{}).
		Where("owner_id = ? AND divide = ?", accountID, divideName).
		Update("divide", uncategorizedDivide).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "移动分组成员到未分类失败"})
		return
	}

	// 删除指定分组
	if err := global.Db.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除分组失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "分组删除成功，成员已移入未分类分组"})
}

// RenameDivide 重命名分组
func RenameDivide(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}

	// 获取URL中的type参数（friends 或 groups）
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 解析请求体中的参数
	var requestBody struct {
		OldFdName string `json:"old_fd_name"` // 旧分组名称
		NewFdName string `json:"new_fd_name"` // 新分组名称
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "JSON绑定失败"})
		return
	}
	if requestBody.OldFdName == "" || requestBody.NewFdName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "old_fd_name和new_fd_name参数不能为空"})
		return
	}

	// 确定是否是群聊
	var isGroupChat bool
	if groupType == "groups" {
		isGroupChat = true
	} else {
		isGroupChat = false
	}

	// 检查旧分组是否存在
	var oldGroup models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, requestBody.OldFdName, isGroupChat).First(&oldGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "旧分组不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询旧分组失败"})
		}
		return
	}

	// 检查新分组名称是否已存在
	var newGroup models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, requestBody.NewFdName, isGroupChat).First(&newGroup).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "新分组名称已存在"})
		return
	}

	// 更新旧分组的名称为新分组名称
	err = global.Db.Model(&oldGroup).Update("divide", requestBody.NewFdName).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新分组名称失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "分组名称更新成功"})
}

// MoveInDivide 将传入的成员的分组改为相应分组
func MoveInDivide(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "http的Header中用户ID为空"})
		return
	}

	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID解析失败"})
		return
	}
	var user models.AccountInfo
	if err := global.Db.Where("id = ?", accountID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	if user.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户注销"})
		return
	}

	// 获取URL中的type参数（friends 或 groups）
	groupType := c.Param("type")
	if groupType != "friends" && groupType != "groups" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "type参数无效"})
		return
	}

	// 解析请求体中的参数
	var requestBody struct {
		TID    string `json:"tid"`    // 成员ID或者群号
		Divide string `json:"divide"` // 目标分组名称
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "JSON绑定失败"})
		return
	}
	if requestBody.TID == "" || requestBody.Divide == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "tid和divide参数不能为空"})
		return
	}

	// 确定是否是群聊
	var isGroupChat bool
	if groupType == "groups" {
		isGroupChat = true
	} else {
		isGroupChat = false
	}

	//检测进行分组的这个人是否存在，是否注销
	if !isGroupChat {
		var user1 models.AccountInfo
		if err := global.Db.Where("account_id = ?", requestBody.TID).First(&user1).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "需要进行分组的用户不存在"})
			return
		}
		if user1.Deactivate == true {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "需要进行分组的用户已经注销"})
			return
		}
	}

	// 检查目标分组是否存在
	var targetGroup models.Contacts
	err = global.Db.Where("owner_id = ? AND divide = ? AND is_group_chat = ?", accountID, requestBody.Divide, isGroupChat).First(&targetGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "目标分组不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询目标分组失败"})
		}
		return
	}

	// 检查当前id用户的这个好友/群聊是否存在
	var member models.Contacts
	err = global.Db.Where("owner_id = ? AND tid = ? AND is_group_chat = ?", accountID, requestBody.TID, isGroupChat).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "成员不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询成员失败"})
		}
		return
	}

	err = global.Db.Model(&member).Update("divide", requestBody.Divide).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新成员分组失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成员已成功移入目标分组"})
}

// ---------------------------------------------------------------------------
/*群聊部分*/

// GetGroups 获取群聊列表
func GetGroups(c *gin.Context) {
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

	// 查询用户所属的群聊信息
	var contacts []models.Contacts
	err = global.Db.Where("owner_id = ? AND is_group_chat = ? AND is_blocked = ?", accountID, true, false).
		Find(&contacts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊失败"})
		return
	}

	var groupList []map[string]interface{}

	for _, contact := range contacts {
		var group models.GroupChatInfo
		err := global.Db.Where("group_id = ?", contact.ContactID).First(&group).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询群聊失败"})
			continue
		}

		// 将群聊信息添加到结果数组
		groupList = append(groupList, map[string]interface{}{
			"avatar":     group.GroupAvatar,       // 群头像
			"account_id": group.GroupID,           // 群ID
			"signature":  group.GroupIntroduction, // 群介绍
			"remark":     contact.Remark,          // 群名称或备注
			"tag":        contact.Divide,          // 分组名称
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": groupList})
}

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
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

	var input struct {
		GroupAvatar string `json:"group_avatar"`
		GroupName   string `json:"group_name"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数解析失败"})
		return
	}
	if input.GroupAvatar == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊头像不能为空"})
		return
	}
	if input.GroupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊名称不能为空"})
		return
	}

	group := models.GroupChatInfo{
		// GroupID自增自动生成
		GroupAvatar: input.GroupAvatar,
		GroupName:   input.GroupName,
	}
	if err := global.Db.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建群聊失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建群聊成功"})
}

// DisMissGroup 解散群聊
func DisMissGroup(c *gin.Context) {
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

	var input struct {
		GroupID uint `json:"group_id"`
	}
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 非群主无法解除
	if uint(accountID) != group.GroupOwner {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主才能解散群聊"})
		return
	}

	if err := global.Db.Where("group_id = ?", input.GroupID).Delete(&models.GroupChatInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "解散群聊失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "解散群聊成功"})
}

// Invite 邀请其他人加入群聊
func Invite(c *gin.Context) {
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

	var input struct {
		AccountID uint `json:"account_id"`
		GroupID   uint `json:"group_id"`
	}

	//检测被邀请的用户是否存在
	var other models.AccountInfo
	if err := global.Db.Where("account_id = ?", input.AccountID).First(&other).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "想邀请的用户不存在"})
		return
	}
	if other.Deactivate == true {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "想邀请的用户已注销"})
		return
	}

	// 群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 是否已经在群聊中
	var member models.GroupMemberInfo
	if err := global.Db.Where("account_id = ? AND group_id = ?", input.AccountID, input.GroupID).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "该成员已经在群聊中"})
		return
	}

	newMember := models.GroupMemberInfo{
		GroupID:   input.GroupID,
		AccountID: input.AccountID,
		GroupRole: "group_ordinary",
	}
	if err := global.Db.Create(&newMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "邀请失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "邀请成功"})
}

// Quit 退出群聊
func Quit(c *gin.Context) {
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

	var input struct {
		GroupID string `json:"group_id"`
	}

	if input.GroupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少必要的参数"})
		return
	}

	// 检测该群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 检测该用户是否为群聊成员
	var groupMember models.GroupMemberInfo
	if err := global.Db.Where("AccountID = ? AND group_id = ?", accountID, input.GroupID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该用户不是群聊成员，无法退出"})
		return
	}

	// 退出群聊
	if err := global.Db.Where("AccountID = ? AND group_id = ?", accountID, input.GroupID).Delete(&groupMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "退出群聊失败，请稍后重试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功退出群聊"})
}

// GetGroupInfo 获取群聊的详细信息
func GetGroupInfo(c *gin.Context) {
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

	groupID := c.Param("group_id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "缺少群聊ID参数"})
		return
	}

	// 检测群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", groupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 获取Group_id群成员信息
	var members []models.GroupMemberInfo
	if err := global.Db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取群聊成员失败"})
		return
	}

	memberList := make([]gin.H, len(members))
	for i, member := range members {
		memberList[i] = gin.H{
			"user_id":   member.AccountID,
			"user_name": member.GroupNickname,
		}
	}

	var groupChat models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", groupID).First(&groupChat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取群聊信息失败"})
		return
	}

	// 前后端可能返回不一致，需注意
	c.JSON(http.StatusOK, gin.H{
		"success":            true,
		"message":            "成功",
		"group_id":           groupID,
		"group_name":         group.GroupName,
		"group_introduction": groupChat.GroupIntroduction,
		"members":            memberList,
	})
}

// ChangeNickname 更改用户在群聊内的昵称
func ChangeNickname(c *gin.Context) {
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

	var input struct {
		GroupID       uint   `json:"group_id"`       // 群组ID
		GroupNickname string `json:"group_nickname"` // 新的群昵称
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 检查群聊是否存在
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 检查用户是否是群成员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, accountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户不是该群的成员"})
		return
	}

	// 更新用户的群昵称
	membership.GroupNickname = input.GroupNickname
	if err := global.Db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新昵称失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "昵称更新成功"})
}

// FetchFriendsNotInGroup 搜索不在群聊内的好友
func FetchFriendsNotInGroup(c *gin.Context) {
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
	var input struct {
		GroupID uint `json:"group_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 检查群聊是否存在
	var group models.GroupChatInfo
	if err = global.Db.Where("group_id=?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
		return
	}

	// 获取用户的好友列表
	var friends []models.Contacts
	if err := global.Db.Where("owner_id = ?", accountID).Find(&friends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取好友列表失败"})
		return
	}

	// 获取已在群聊内的成员列表
	var groupMembers []models.GroupMemberInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).Find(&groupMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取群成员列表失败"})
		return
	}

	// 将已在群聊内的成员ID存入map
	memberIDs := make(map[uint]bool)
	for _, member := range groupMembers {
		memberIDs[member.AccountID] = true
	}

	// 筛选不在群聊内的好友
	var friendsNotInGroup []models.Contacts
	for _, friend := range friends {
		if !memberIDs[friend.ContactID] {
			friendsNotInGroup = append(friendsNotInGroup, friend)
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": friendsNotInGroup})
}

// RemoveMember 移除某个群成员
func RemoveMember(c *gin.Context) {
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
	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 查找当前群组信息
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主或管理员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, accountID).First(&membership).Error; err != nil || (membership.GroupRole != "group_admin" && group.GroupOwner != uint(accountID)) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主或管理员可以移除成员"})
		return
	}

	// 检查目标用户是否是群成员
	if err := global.Db.Where("group_id = ? AND member_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员，没法移除"})
		return
	}

	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).Delete(&models.GroupMemberInfo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "移除成员失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "成功移除群成员"})
}

// SetBanned 设置禁言或取消禁言
func SetBanned(c *gin.Context) {
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
	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
		IsBanned  bool `json:"is_banned"`  // 是否禁言
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 查找当前群组信息
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主或管理员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, accountID).First(&membership).Error; err != nil || (membership.GroupRole != "group_ordinary" && group.GroupOwner != uint(accountID)) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主或管理员可以设置禁言"})
		return
	}

	// 检查目标用户是否是群成员
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员"})
		return
	}

	// 更新禁言状态
	membership.IsBanned = input.IsBanned
	if err := global.Db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新禁言状态失败"})
		return
	}

	// 返回成功响应
	action := "设置禁言"
	if !input.IsBanned {
		action = "取消禁言"
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("成功%s", action)})
}

// SetAdmin 设置管理员
func SetAdmin(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
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

	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
		IsAdmin   bool `json:"is_admin"`   // 是否设置为管理员
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主
	if group.GroupOwner != uint(accountID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主可以设置管理员"})
		return
	}

	// 检查目标用户是否是群成员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员"})
		return
	}

	if input.IsAdmin {
		membership.GroupRole = "group_admin"
	} else {
		membership.GroupRole = "group_ordinary"
	}
	if err := global.Db.Save(&membership).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新管理员状态失败"})
		return
	}

	action := "设置为管理员"
	if !input.IsAdmin {
		action = "取消管理员"
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("成功%s", action),
	})
}

// TransferOwner 更换群主
func TransferOwner(c *gin.Context) {
	ID := c.GetHeader("User-ID")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID为空，请检查请求头"})
		return
	}
	accountID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "用户ID无效"})
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

	var input struct {
		GroupID   uint `json:"group_id"`   // 群组ID
		AccountID uint `json:"account_id"` // 目标用户ID
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数格式错误"})
		return
	}

	// 查找当前群组信息
	var group models.GroupChatInfo
	if err := global.Db.Where("group_id = ?", input.GroupID).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标群组不存在"})
		return
	}

	// 检查当前操作用户是否是群主
	if group.GroupOwner != uint(accountID) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "只有群主可以转让群主权限"})
		return
	}

	// 检查目标用户是否是群成员
	var membership models.GroupMemberInfo
	if err := global.Db.Where("group_id = ? AND account_id = ?", input.GroupID, input.AccountID).First(&membership).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "目标用户不是群成员"})
		return
	}

	// 转让群主权限
	if err := global.Db.Model(&group).Update("owner_id", input.AccountID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "群主转让失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "群主权限已成功转让"})
}

//--------------------------------------------------------------------------

/* 获取各类资料卡片*/

// GetProfileCard 获取指定用户的资料卡片
func GetProfileCard(c *gin.Context) {
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

	// 用户ID或者群号
	other := c.Param("tid")
	groupID := c.DefaultPostForm("group_id", "") // 获取POST请求中的group_id参数
	if groupID != "" {
		var groupChat models.GroupChatInfo
		if err := global.Db.Where("group_id = ?", groupID).First(&groupChat).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "群聊不存在"})
			return
		}

		var contact models.Contacts
		if err := global.Db.Where("contact_id = ? AND owner_id = ?", groupID, accountID).First(&contact).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "不是群聊成员"})
			return
		}

		responseData := gin.H{
			"tid":           groupChat.GroupID,
			"avatar":        groupChat.GroupAvatar,
			"remark":        contact.Remark,
			"groupNickname": groupChat.GroupName,
			"tag":           contact.Divide,
			"signature":     groupChat.GroupIntroduction,
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功",
			"data":    responseData,
		})
	} else {
		var another models.AccountInfo
		if err := global.Db.Where("account_id = ?", other).First(&another).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		if another.Deactivate {
			c.JSON(http.StatusNotFound, gin.H{"error": "目标用户已注销"})
			return
		}

		var contact models.Contacts
		if err := global.Db.Where("owner_id = ? AND contact_id = ?", accountID, another.AccountID).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "不是好友关系"})
			return
		}

		responseData := gin.H{
			"tid":       another.AccountID,
			"avatar":    another.Avatar,
			"remark":    contact.Remark,
			"nickname":  another.Nickname,
			"tag":       contact.Divide,
			"signature": another.Signature,
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功",
			"data":    responseData,
		})
	}
}

// -------------------------------------------------------------------------
/*设置*/

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
		ID        string `json:"id"`
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
	user.ID = input.ID
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
		FontSize int `json:"font_size"`
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
		Notice string `json:"notice"`
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
		NoticeGroup string `json:"notice_group"`
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
