package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func SetupMessageListRoutes(router *gin.Engine) *gin.RouterGroup {

	r := router.Group("/messageList")

	// 好友相关
	r.POST("/friend/add", controllers.AddFriend) // 添加好友

	// 群聊相关
	r.GET("/group/:id", controllers.GetMessages)     // 获取群聊信息
	r.POST("/group/create", controllers.CreateGroup) // 创建群聊

	// 消息相关
	r.GET("/messages/:id", controllers.GetMessages)   // 3.9.获取聊天记录
	r.POST("/messages/send", controllers.SendMessage) // 3.10.发送消息

	// 申请通知相关
	r.POST("/apply/submit", controllers.SubmitApply) // 提交申请
	r.PUT("/apply/review", controllers.ReviewApply)  // 审核申请

	return r
}
