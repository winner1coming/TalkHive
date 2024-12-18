package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func SetupPersonalHomeRoutes(router *gin.Engine) *gin.RouterGroup {

	r := router.Group("/personalHome")

	// 资料编辑相关
	r.POST("/friend/add", controllers.UpdateAvatar)  // 更新用户头像
	r.GET("/group/:id", controllers.UpdateNickname)  // 修改用户昵称
	r.POST("/group/create", controllers.EditProfile) // 编辑其他用户自来哦：性别、签名……

	return r
}
