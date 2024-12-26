package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 登录和注册相关路由
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.POST("/resetPassword", controllers.ResetPassword)
	r.POST("/sendSmsCode", controllers.SendSmsCode)
	r.POST("/smslogin", controllers.SmsLogin)

	//好友请求部分
	r.GET("/contactList/friendRequests", controllers.GetFriendRequests)
	r.POST("/contactList/friendRequests/pend", controllers.FriendRequestPend)
	r.POST("/contactList/friendRequests/addFriend", controllers.AddFriend)

	//群聊请求部分
	r.GET("/contactList/groupRequests", controllers.GetGroupRequests)
	r.POST("/contactList/groupRequests/applyPend", controllers.DealGroupApplyRequest)
	r.POST("/contactList/groupRequests/invitePend", controllers.DealGroupInviteRequest)

	//黑名单部分
	r.GET("/contactList/blackList", controllers.GetBlackList)
	r.POST("/contactList/blackList/remove", controllers.RemoveFromBlacklist)
	r.POST("/contactList/blackList/add", controllers.AddToBlacklist)

	// 好友列表
	r.GET("/contactList/friends", controllers.GetFriends)

	// 分组部分
	r.GET("/contactList/{type}/divides", controllers.GetDivides)
	r.POST("/contactList/{type}/divides/create", controllers.CreateDivide)
	r.DELETE("/contactList/{type}/divides/delete/{fd_name}", controllers.DeleteDivide)
	r.POST("/contactList/{type}/divides/rename", controllers.RenameDivide)
	r.POST("/contactList/{type}/divides/moveIn", controllers.MoveInDivide)

	// 群聊部分
	r.GET("/contactList/groups", controllers.GetGroups)
	r.POST("/contactList/createGroup", controllers.CreateGroup)
	r.POST("/contactList/group/invite", controllers.Invite)
	r.POST("/contactList/group/quit", controllers.Quit)
	r.POST("/contactList/group/dismissGroup", controllers.DisMissGroup)
	r.GET("/contactList/group/{group_id}", controllers.GetGroupInfo)

	// 资料卡片
	r.GET("/contactList/card", controllers.GetProfileCard)

	// 个人主业路由
	r.GET("/Settings/showProfile", controllers.ShowProfile)
	r.POST("/Settings/saveEdit", controllers.SaveEdit)
	r.GET("/Settings/getInfo", controllers.GetUserInfo)
	r.POST("/Settings/getCode", controllers.GetCode)
	r.POST("/Settings/saveEmail", controllers.SaveEmail)
	r.POST("/Settings/savePassword", controllers.SavePassword)
	r.POST("/Settings/isIDAdd", controllers.IsIDAdd)
	r.POST("/Settings/isNicknameAdd", controllers.IsNickNameAdd)
	r.POST("/Settings/deactivation", controllers.ConfirmDeactivation)

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
	r.POST("/logout", controllers.Logout)
}
