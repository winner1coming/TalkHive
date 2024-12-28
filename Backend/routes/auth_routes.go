package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

// AuthRoutes 登录和注册相关路由
func AuthRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.POST("/resetPassword", controllers.ResetPassword)
	r.POST("/sendSmsCode", controllers.SendSmsCode)
	r.POST("/smslogin", controllers.SmsLogin)
}
