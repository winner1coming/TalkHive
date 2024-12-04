package routes

import (
	"TalkHive/controllers/groupChat"
	"github.com/gin-gonic/gin"
)

func SetupGroupChatRoutes(router *gin.Engine) *gin.RouterGroup {
	r := router.Group("/groupChat")

	// 群聊基本功能
	r.GET("/group/:id", groupChat.GetGroupInfo)                 // 显示群头像、名称、成员详情等
	r.POST("/group/:id/send/text", groupChat.SendTextMessage)   // 发送文本消息
	r.POST("/group/:id/send/image", groupChat.SendImageMessage) // 发送图片
	r.POST("/group/:id/send/video", groupChat.SendVideoMessage) // 发送视频
	r.POST("/group/:id/send/file", groupChat.SendFileMessage)   // 发送文件
	r.POST("/group/:id/send/code", groupChat.SendCodeFile)      // 发送代码文件
	r.POST("/group/:id/send/emoji", groupChat.SendEmojiMessage) // 发送表情包

	// 消息管理
	r.POST("/group/:id/message/copy", groupChat.CopyMessage)             // 复制消息
	r.POST("/group/:id/message/favorite", groupChat.FavoriteMessage)     // 收藏消息
	r.POST("/group/:id/message/reply", groupChat.ReplyMessage)           // 回复消息
	r.POST("/group/:id/message/forward", groupChat.ForwardMessage)       // 转发消息
	r.DELETE("/group/:id/message/delete", groupChat.DeleteMessage)       // 删除消息
	r.POST("/group/:id/message/select", groupChat.BatchMessageOperation) // 多选消息（收藏、转发、删除）

	// 群公告与聊天记录
	r.GET("/group/:id/announcement", groupChat.GetAnnouncements)         // 查看群公告
	r.POST("/group/:id/announcement", groupChat.CreateAnnouncement)      // 创建群公告（管理员权限）
	r.GET("/group/:id/chat-history", groupChat.GetChatHistory)           // 查找聊天记录
	r.GET("/group/:id/chat-history/search", groupChat.SearchChatHistory) // 带条件筛选聊天记录

	// 成员互动
	r.POST("/group/:id/at-all", groupChat.AtAllMembers)      // 管理员@所有人
	r.POST("/group/:id/at-member", groupChat.AtSingleMember) // 普通成员@单人

	// 群聊详情
	r.GET("/group/:id/details", groupChat.GetGroupDetails)           // 显示群聊详情
	r.POST("/group/:id/member/search", groupChat.SearchMember)       // 搜索群成员
	r.POST("/group/:id/invite", groupChat.InviteMember)              // 邀请好友加入
	r.POST("/group/:id/member/nickname", groupChat.SetGroupNickname) // 设置群昵称
	r.POST("/group/:id/note", groupChat.SetGroupNote)                // 设置群聊备注
	r.POST("/group/:id/settings", groupChat.UpdateGroupSettings)     // 设置消息免打扰、置顶
	r.DELETE("/group/:id/exit", groupChat.ExitGroup)                 // 退出群聊

	// 群聊管理（管理员功能）
	r.POST("/group/:id/manage/avatar", groupChat.UpdateGroupAvatar)            // 修改群头像
	r.POST("/group/:id/manage/name", groupChat.UpdateGroupName)                // 修改群名称
	r.POST("/group/:id/manage/permissions", groupChat.UpdateJoinPermissions)   // 设置入群权限
	r.POST("/group/:id/manage/intro", groupChat.UpdateGroupIntro)              // 编辑群简介
	r.POST("/group/:id/manage/mute", groupChat.SetMuteStatus)                  // 设置禁言（个体或全员）
	r.DELETE("/group/:id/manage/member", groupChat.RemoveMember)               // 移除群成员
	r.DELETE("/group/:id/manage", groupChat.DisbandGroup)                      // 解散群聊
	r.POST("/group/:id/manage/applications", groupChat.HandleJoinApplications) // 管理入群申请

	// 群主特有功能
	r.POST("/group/:id/owner/transfer", groupChat.TransferOwnership) // 转让群主
	r.POST("/group/:id/owner/add-admin", groupChat.AddAdministrator) // 添加管理员（最多五人）

	// 其他路由...
	return r
}
