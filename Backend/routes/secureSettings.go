package routes

import (
	"TalkHive/controllers/settings"
	"github.com/gin-gonic/gin"
)

func SetupSecureSettingsRoutes(router *gin.Engine) *gin.RouterGroup {

	r := router.Group("/Settings")

	// 安全设置相关
	r.GET("/getPhone", settings.GetPhone)                 // 1.展示手机号：√×
	r.POST("/getCode", settings.GetCode)                  // 2.获取验证码：√×
	r.POST("/savePhone", settings.SavePhone)              // 3.保存修改手机号：√×
	r.POST("/savePassword", settings.SavePassword)        // 4.保存密码的更改：√×
	r.POST("/idAdd", settings.IsIDAdd)                    // 5.设置是否允许通过ID查找：√×
	r.POST("/phoneAdd", settings.IsPhoneAdd)              // 6.设置是否允许通过手机号查找：√×
	r.POST("/deactivation", settings.ConfirmDeactivation) // 7.注销账号：√×

	return r
}
