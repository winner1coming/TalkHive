package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 用户账号相关
	r.POST("/register", controllers.Register)           // 用户注册
	r.POST("/login", controllers.Login)                 // 用户登录
	r.GET("/Settings/profile", controllers.ShowProfile) // 展示用户个人主页资料：√×
	r.POST("/Settings/saveEdit", controllers.SaveEdit)  // 保存用户编辑后的资料：√×

	return r
}
