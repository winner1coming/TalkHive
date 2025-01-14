package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func ChatListRoutes(r *gin.Engine) {
	// 聊天列表
	r.GET("/chatlist", controllers.GetChatList)
	r.POST("/chatlist/createChat", controllers.CreateChat)
	r.GET("/chatlist/search/:keyword", controllers.SearchChats)
	r.POST("/chatlist/pin", controllers.PinChat)
	r.POST("/messages/read", controllers.ReadMessages)
	r.POST("/chatlist/delete", controllers.DeleteChat)
	r.POST("/chatlist/mute", controllers.SetMute)
	r.POST("/chatlist/block", controllers.BlockChat)

	//聊天消息
	r.POST("/messages", controllers.GetMessages)
	r.POST("/messages/send", controllers.SendMessage)
	r.POST("/messages/sendFile", controllers.SendFile)
	r.POST("/messages/collect", controllers.CollectMessage)
	r.POST("/messages/delete", controllers.DeleteMessage)
}
