package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	// 用户账号相关
	r.POST("/register", controllers.Register)               // 用户注册
	r.POST("/login", controllers.Login)                     // 用户登录
	r.POST("/sendSmsCode", controllers.SendSmsCode)         // 发送短信验证码
	r.POST("/smsLogin", controllers.SmsLogin)               // 验证短信验证码
	r.POST("/resetPassword", controllers.ResetPassword)     // 重置密码
	r.GET("/profile/:id", controllers.GetProfile)           // 获取用户资料
	r.PUT("/profile/update/:id", controllers.UpdateProfile) // 更新用户资料

}
