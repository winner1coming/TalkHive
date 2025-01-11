package models

import (
	"time"
)

// MessageInfo 消息表
type MessageInfo struct {
	MessageID     uint   `gorm:"primaryKey" json:"message_id"`
	CreateTime    string `json:"create_time"`
	SendAccountID uint   `json:"send_account_id"`
	TargetID      uint   `json:"receive_account_id"`
	Content       string `json:"content"`
	Type          string `json:"type"`
	ChatID        uint   `json:"chat_id" gorm:"foreignKey:ChatID;references:ChatID"`
	IsRead        bool   `json:"is_read"`
}

// ChatInfo 聊天记录表
type ChatInfo struct {
	ChatID     uint   `gorm:"primaryKey" json:"chat_id"`
	CreateTime string `json:"create_time"`
	AccountID  uint   `json:"account_id"`
	TargetID   uint   `json:"target_id"`
	IsGroup    bool   `json:"is_group"`
}

// AccountInfo 账号信息表
type AccountInfo struct {
	AccountID                uint   `gorm:"primaryKey" json:"account_id"`
	ID                       string `gorm:"unique" json:"id"`
	Password                 string `json:"password"`
	Phone                    string `json:"phone"`
	Email                    string `json:"email"`
	Avatar                   string `json:"avatar"`
	Nickname                 string `json:"nickname"`
	Signature                string `json:"signature"`
	Gender                   string `json:"gender"`
	Birthday                 string `json:"birthday"`
	Status                   string `json:"status"`
	FriendPermissionID       bool   `json:"friend_permission"`
	FriendPermissionNickName bool   `json:"friend_permission_nick_name"`
	LastLogout               string `json:"last_logout"`
	Deactivate               bool   `json:"deactivate"`
}

// Contacts 好友/群聊表
type Contacts struct {
	OwnerID          uint   `json:"owner_id"`
	ContactID        uint   `gorm:"primaryKey" json:"contact_id"`
	IsBlacklist      bool   `json:"is_blacklist"`
	IsPinned         bool   `json:"is_pinned"`
	Divide           string `json:"divide"`
	IsMute           bool   `json:"is_mute"`
	IsBlocked        bool   `json:"is_blocked"`
	IsGroupChat      bool   `json:"is_group_chat"`
	Remark           string `json:"remark"`
	UnreadMessageNum uint   `json:"unread_message_num"`
}

// SystemSetting 系统环境设置表
type SystemSetting struct {
	AccountID   uint   `gorm:"primaryKey" json:"account_id"`
	Background  string `json:"background"`
	FontStyle   string `json:"font_style"`
	FontSize    uint   `json:"font_size"`
	Theme       string `json:"theme"`
	Sound       string `json:"sound"`
	Notice      bool   `json:"notice"`
	NoticeGroup bool   `json:"noticeGroup"`
}

// ApplyInfo 申请通知表
type ApplyInfo struct {
	ApplyID    uint   `gorm:"primaryKey;autoIncrement" json:"apply_id"`
	ApplyType  string `json:"apply_type"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	GroupID    uint   `json:"group_id"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
}

// GroupChatInfo 群聊	总表
type GroupChatInfo struct {
	GroupID           uint   `gorm:"primaryKey" json:"group_id"`
	GroupOwner        uint   `json:"group_owner"`
	GroupAvatar       string `json:"group_avatar"`
	GroupName         string `json:"group_name"`
	GroupIntroduction string `json:"group_introduction"`
	AllowInvite       bool   `json:"allow_invite"`
	AllowIDSearch     bool   `json:"allow_id_search"`
	ALlowNameSearch   bool   `json:"allow_name_search"`
	IsAllBanned       bool   `json:"is_all_banned"`
}

// GroupMemberInfo 群成员信息表
type GroupMemberInfo struct {
	AccountID     uint   `gorm:"primaryKey" json:"account_id"`
	GroupID       uint   `gorm:"primaryKey" json:"group_id"`
	GroupNickname string `json:"group_nickname"`
	IsBanned      bool   `json:"is_banned"`
	GroupRole     string `json:"group_role"`
}

// Notes 笔记表
type Notes struct {
	NoteID    uint      `gorm:"primaryKey" json:"note_id"`
	Type      string    `json:"type"`
	CachePath string    `json:"cache_path"`
	AccountID uint      `json:"account_id"`
	IsShow    bool      `json:"is_show"`
	SaveTime  time.Time `json:"save_time"`
	Name      string    `json:"name"`
}

// Favorites 收藏表
type Favorites struct {
	TableName string `gorm:"primaryKey" json:"table_name"`
	ID        uint   `json:"id"`
	AccountID uint   `json:"account_id"`
}

// Codes 代码表
type Codes struct {
	CodeID    uint      `gorm:"primaryKey" json:"code_id"`
	Name      string    `json:"name"`
	SaveTime  time.Time `json:"save_time"`
	CachePath string    `json:"cache_path"`
	Suffix    string    `json:"suffix"`
	AccountID uint      `json:"account_id"`
	IsShow    bool      `json:"is_show"`
}

// DDLS DDL表
type DDLS struct {
	DDLID       uint      `gorm:"primaryKey" json:"ddl_id"`
	AccountID   uint      `json:"account_id"`
	DDLDate     time.Time `json:"ddl_date"`
	Task        string    `json:"task"`
	IsCompleted bool      `json:"is_completed"`
	Urgency     bool      `json:"urgency"`
}

// Recycle 回收站表
type Recycle struct {
	RecycleID   uint      `gorm:"primaryKey" json:"recycle_id"`
	RecycleTime time.Time `json:"recycle_time"`
	RecycleType string    `json:"recycle_type"`
	AccountID   uint      `json:"account_id"`
}

// GroupDivide 群聊分组表
type GroupDivide struct {
	GroupDivideId uint   `gorm:"primaryKey" json:"group_divide_id"`
	GDName        string `json:"gd_name"`
	AccountID     uint   `json:"account_id"`
}

// FriendDivide 好友分组
type FriendDivide struct {
	FriendDivideId uint   `gorm:"primaryKey" json:"friend_divide_id"`
	FDName         string `json:"fd_name"`
	AccountID      uint   `json:"account_id"`
}

// NoteDivide 笔记分类表
type NoteDivide struct {
	NDName    string `gorm:"primaryKey" json:"nd_name"`
	AccountID uint   `json:"account_id"`
}

// Links 网页链接器表
type Links struct {
	URL     string `gorm:"primaryKey" json:"url"`
	URLName string `json:"url_name"`
	Icon    string `json:"icon"`
}

// DeleteInfo 删除消息表
type DeleteInfo struct {
	TargetID uint   `gorm:"primaryKey" json:"target_id"`
	Range    string `json:"range"`
}
