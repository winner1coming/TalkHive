package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func ChatListRoutes(r *gin.Engine) {
	// 聊天列表
	r.GET("/chatList", controllers.GetChatList)
	r.POST("/chatlist/createChat", controllers.GetChat)
	r.GET("/chatlist/search/:keyword", controllers.SearchChats)
	r.POST("/chatlist/pin", controllers.PinChat)
	r.POST("/messages/read", controllers.ReadMessages)
	r.DELETE("/chatlist/:tid", controllers.DeleteChat)
	r.POST("/chatlist/mute", controllers.SetMute)
	r.POST("/chatlist/block", controllers.BlockChat)

	//聊天消息
	r.GET("/messages/:tid", controllers.GetMessages)
	r.POST("/messages/send", controllers.SendMessage)
	r.POST("/messages/collect", controllers.CollectMessage)
	r.POST("/messages/reply", controllers.ReplyMessage)
	r.POST("/messages/forward", controllers.ForwardMessage)
}
