package routes

import (
	"TalkHive/controllers/settings"
	"github.com/gin-gonic/gin"
)

func SetupSystemSettingsRoutes(router *gin.Engine) *gin.RouterGroup {

	r := router.Group("/systemSettings")

	// 系统设置
	r.POST("/Settings/changeTheme", settings.ChangeTheme)           // 1.更改主题颜色：√×
	r.POST("/Settings/fontSize", settings.ChangeFontsize)           // 2.更改字体大小：√×
	r.POST("/Settings/fontStyle", settings.ChangeFontstyle)         // 3.更换字体风格：√×
	r.POST("/Settings/isNotice", settings.IsNotice)                 // 4.更新消息通知设置：√×
	r.POST("/Settings/noticeGroup", settings.IsNoticeGroup)         // 5.更新群消息声音设置：√×
	r.POST("/Settings/changeSound", settings.ChangeSound)           // 6.更换音频文件设置：√×
	r.POST("/Settings/submitSound", settings.SubmitSound)           // 7.上传新的提示音：√×
	r.POST("/Settings/changeBackground", settings.ChangeBackground) // 8.更改背景图片：√×
	r.GET("/systemSetting", settings.GetSystemSetting)              // 9.获取系统设置：√×

	// 退出登录
	r.POST("/Logout", settings.Logout) // 退出登录：√×

	return r
}
