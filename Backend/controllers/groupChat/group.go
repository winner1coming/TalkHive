package groupChat

import (
	"TalkHive/global"
	"TalkHive/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 获取数据库连接
var db = global.Db

// GetGroupInfo 显示群头像、名称、成员详情等
func GetGroupInfo(c *gin.Context) {
	groupID := c.Param("id")
	var group models.GroupChatInfo
	var members []models.GroupMemberInfo

	// 查询群聊信息
	if err := db.Where("group_id = ?", groupID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// 查询群聊成员信息
	if err := db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Members not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group_info": group,
		"members":    members,
	})
}

// SendTextMessage 发送文本消息
func SendTextMessage(c *gin.Context) {
	var message models.MessageInfo
	var sendAccountID uint // 假设从登录信息中获取发送者的 AccountID

	// 绑定请求体到 message 结构体
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 设置消息类型、群ID以及发送时间
	message.Type = "text"
	message.CreateTime = time.Now()
	message.SendAccountID = sendAccountID

	// 保存消息到数据库
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

// SendImageMessage 发送图片消息
func SendImageMessage(c *gin.Context) {
	var message models.MessageInfo
	var sendAccountID uint // 假设从登录信息中获取发送者的 AccountID

	// 绑定请求体到 message 结构体
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 设置消息类型、群ID以及发送时间
	message.Type = "image"
	message.CreateTime = time.Now()
	message.SendAccountID = sendAccountID
	message.Content = message.Content // 图片路径或URL

	// 保存消息到数据库
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image sent"})
}

// SendVideoMessage 发送视频消息
func SendVideoMessage(c *gin.Context) {
	var message models.MessageInfo
	var sendAccountID uint // 假设从登录信息中获取发送者的 AccountID

	// 绑定请求体到 message 结构体
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 设置消息类型、群ID以及发送时间
	message.Type = "video"
	message.CreateTime = time.Now()
	message.SendAccountID = sendAccountID
	message.Content = message.Content // 视频路径或URL

	// 保存消息到数据库
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video sent"})
}

// SendFileMessage 发送文件消息
func SendFileMessage(c *gin.Context) {
	var message models.MessageInfo
	var sendAccountID uint // 假设从登录信息中获取发送者的 AccountID

	// 绑定请求体到 message 结构体
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 设置消息类型、群ID以及发送时间
	message.Type = "file"
	message.CreateTime = time.Now()
	message.SendAccountID = sendAccountID
	message.Content = message.Content // 文件路径或URL

	// 保存消息到数据库
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File sent"})
}

// SendCodeFile 发送代码文件消息
func SendCodeFile(c *gin.Context) {
	var message models.MessageInfo
	var sendAccountID uint // 假设从登录信息中获取发送者的 AccountID

	// 绑定请求体到 message 结构体
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 设置消息类型、群ID以及发送时间
	message.Type = "code"
	message.CreateTime = time.Now()
	message.SendAccountID = sendAccountID
	message.Content = message.Content // 代码内容或文件路径

	// 保存消息到数据库
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code file sent"})
}

// SendEmojiMessage 发送表情包消息
func SendEmojiMessage(c *gin.Context) {
	var message models.MessageInfo
	var sendAccountID uint // 假设从登录信息中获取发送者的 AccountID

	// 绑定请求体到 message 结构体
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 设置消息类型、群ID以及发送时间
	message.Type = "emoji"
	message.CreateTime = time.Now()
	message.SendAccountID = sendAccountID
	message.Content = message.Content // 表情包内容或路径

	// 保存消息到数据库
	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emoji sent"})
}

// GetAnnouncements 查找群公告
func GetAnnouncements(c *gin.Context) {
	groupID := c.Param("id")

	// 查找群公告，假设公告的消息类型为 "announcement"
	var announcements []models.MessageInfo
	if err := db.Where("group_id = ? AND type = ?", groupID, "announcement").Order("create_time desc").Find(&announcements).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No announcements found"})
		return
	}

	c.JSON(http.StatusOK, announcements)
}

// CreateAnnouncement 创建群公告
func CreateAnnouncement(c *gin.Context) {
	groupID := c.Param("id")
	var request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 验证用户是否为群管理员
	accountID := c.MustGet("account_id").(uint) // 假设账户ID已通过JWT等方式传递
	var groupMember models.GroupMemberInfo
	if err := db.Where("group_id = ? AND account_id = ? AND is_banned = false", groupID, accountID).First(&groupMember).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not an admin"})
		return
	}

	// 创建群公告，公告消息类型为 "announcement"
	announcement := models.MessageInfo{
		SendAccountID: accountID,
		Content:       fmt.Sprintf("Title: %s\nContent: %s", request.Title, request.Content),
		Type:          "announcement",
		CreateTime:    time.Now(),
	}

	// 保存公告
	if err := db.Create(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create announcement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement created successfully"})
}

// GetChatHistory 查找聊天记录
func GetChatHistory(c *gin.Context) {
	groupID := c.Param("id")
	var messages []models.MessageInfo

	// 查找该群的聊天记录
	if err := db.Where("group_id = ?", groupID).Order("create_time asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No chat history found"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// SearchChatHistory 带条件筛选聊天记录
func SearchChatHistory(c *gin.Context) {
	groupID := c.Param("id")
	var request struct {
		SenderID    uint   `json:"sender_id,omitempty"`
		MessageType string `json:"message_type,omitempty"`
		StartDate   string `json:"start_date,omitempty"` // 可选筛选条件
		EndDate     string `json:"end_date,omitempty"`   // 可选筛选条件
	}

	// 解析请求数据
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// 构建查询条件
	query := db.Where("group_id = ?", groupID)

	if request.SenderID != 0 {
		query = query.Where("send_account_id = ?", request.SenderID)
	}
	if request.MessageType != "" {
		query = query.Where("type = ?", request.MessageType)
	}
	if request.StartDate != "" {
		query = query.Where("create_time >= ?", request.StartDate)
	}
	if request.EndDate != "" {
		query = query.Where("create_time <= ?", request.EndDate)
	}

	// 获取筛选后的聊天记录
	var messages []models.MessageInfo
	if err := query.Order("create_time asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching chat history found"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
