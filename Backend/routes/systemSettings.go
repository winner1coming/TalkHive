package routes

import (
	"TalkHive/controllers/settings"
	"github.com/gin-gonic/gin"
)

func SetupSystemSettingsRoutes(router *gin.Engine) *gin.RouterGroup {

	r := router.Group("/systemSettings")

	r.POST("/Settings/changeTheme", settings.ChangeTheme)     // 1.更改主题颜色：√×
	r.POST("/Settings/fontSize", settings.ChangeFontsize)     // 2.更改字体大小：√×
	r.POST("/Settings/fontStyle", settings.ChangeFontstyle)   // 3.更换字体风格：√×
	r.POST("/Settings/isNotice", settings.UpdateNotification) // 4.更新消息通知设置

	return r
}
