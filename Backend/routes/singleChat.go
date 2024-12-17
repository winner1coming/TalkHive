package routes

import (
	"TalkHive/controllers/singleChat"
	"github.com/gin-gonic/gin"
)

func SetupSingleChatRoutes(router *gin.Engine) *gin.RouterGroup {
	r := router.Group("/singleChat")

	// 此处的单聊路由设计，是针对群聊路由的补充
	// 二者重复之处不予考虑
	r.GET("/chat/:id/friend/info", singleChat.GetFriendInfo)               // 查看好友个人信息
	r.PUT("/chat/:id/friend/remark", singleChat.SetFriendRemark)           // 设置好友备注
	r.POST("/chat/:id/friend/card", singleChat.SendFriendCard)             // 发送好友名片
	r.POST("/chat/:id/friend/blacklist", singleChat.AddToBlacklist)        // 加入黑名单
	r.DELETE("/chat/:id/friend/blacklist", singleChat.RemoveFromBlacklist) // 移除黑名单
	r.DELETE("/chat/:id/friend", singleChat.DeleteContact)                 // 删除联系人
	r.POST("/chat/:id/mute", singleChat.SetMute)                           // 3.7.设置免打扰
	r.POST("/chat/:id/pin", singleChat.SetPin)                             // 3.4.设置置顶
	r.DELETE("/chat/:id/records", singleChat.DeleteChatRecords)            // 删除聊天记录

	return r
}
