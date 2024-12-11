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

	// 安全设置相关
	r.GET("/Settings/getPhone", controllers.GetPhone)                 // 1.展示手机号：√×
	r.POST("/Settings/getCode", controllers.GetCode)                  // 2.获取验证码：√×
	r.POST("/Settings/savePhone", controllers.SavePhone)              // 3.保存修改手机号：√×
	r.POST("/Settings/savePassword", controllers.SavePassword)        // 4.保存密码的更改：√×
	r.POST("/Settings/idAdd", controllers.IsIDAdd)                    // 5.设置是否允许通过ID查找：√×
	r.POST("/Settings/phoneAdd", controllers.IsPhoneAdd)              // 6.设置是否允许通过手机号查找：√×
	r.POST("/Settings/deactivation", controllers.ConfirmDeactivation) // 7.注销账号

	// 系统设置相关
	r.POST("/apply/submit", controllers.UpdateSystemSettings) // 更新系统设置

	return r
}
