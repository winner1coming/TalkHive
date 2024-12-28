package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func SettingsRoutes(r *gin.Engine) {

	// 个人主业路由
	r.GET("/Settings/profile", controllers.ShowProfile)
	r.POST("/Settings/saveEdit", controllers.SaveEdit)
	r.GET("/Settings/getInfo", controllers.GetUserInfo)
	r.POST("/Settings/getCode", controllers.GetCode)
	r.POST("/Settings/saveEmail", controllers.SaveEmail)
	r.POST("/Settings/savePassword", controllers.SavePassword)
	r.POST("/Settings/isIDAdd", controllers.IsIDAdd)
	r.POST("/Settings/isNicknameAdd", controllers.IsNickNameAdd)
	r.POST("/Settings/deactivate", controllers.ConfirmDeactivation)

	// 系统设置
	r.POST("/Settings/changeTheme", controllers.ChangeTheme)
	r.POST("/Settings/fontSize", controllers.ChangeFontsize)
	r.POST("/Settings/fontStyle", controllers.ChangeFontstyle)
	r.POST("/Settings/isNotice", controllers.IsNotice)
	r.POST("/Settings/isNoticeGroup", controllers.IsNoticeGroup)
	r.POST("/Settings/changeSound", controllers.ChangeSound)
	r.POST("/Settings/submitSound", controllers.SubmitSound) // 需要加入新的表
	r.POST("/Settings/changeBackground", controllers.ChangeBackground)
	r.GET("/systemSetting", controllers.GetSystemSetting)

	// 退出登录
	r.POST("/Logout", controllers.Logout)
}
