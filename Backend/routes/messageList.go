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
	r.GET("/messages/:id", controllers.GetMessages)   // 获取消息记录
	r.POST("/messages/send", controllers.SendMessage) // 发送消息

	// 系统设置相关
	r.PUT("/settings/system", controllers.UpdateSystemSetting) // 更新系统设置

	// 申请通知相关
	r.POST("/apply/submit", controllers.SubmitApply) // 提交申请
	r.PUT("/apply/review", controllers.ReviewApply)  // 审核申请

	return r
}
