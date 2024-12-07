package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

// 设置路由
func SetupRoutes(r *gin.Engine) {
	// 登录和注册相关路由
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.POST("/resetPassword", controllers.ResetPassword)
	r.POST("/sendSmsCode", controllers.SendSmsCode)
	r.POST("/smsLogin", controllers.SmsLogin)

	// 通信录相关路由
	r.GET("/contactList/friendRequests/:id", controllers.FriendRequests)
	r.POST("/contactList/friendRequests/accept/:id", controllers.AcceptFriendRequest)
	r.POST("/contactList/friendRequests/refuse/:id", controllers.RejectFriendRequest)

}
