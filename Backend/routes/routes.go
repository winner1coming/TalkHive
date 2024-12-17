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
	// 聊天列表
	/*

	 */

	// 通信录相关路由
	r.GET("/contactList/friendRequests/:id", controllers.FriendRequests)
	r.POST("/contactList/friendRequests/pend/:id", controllers.FriendRequestPend)
	//r.POST("/contactList/friendRequests/accept/:id", controllers.AcceptFriendRequest)
	//r.POST("/contactList/friendRequests/refuse/:id", controllers.RejectFriendRequest)

	r.GET("/contactList/groupRequests/:id", controllers.GetGroupRequests)
	r.POST("/contactList/groupRequests/applyPend/:id", controllers.DealGroupApplyRequest)
	r.POST("/contactList/groupRequests/invitePend/:id", controllers.DealGroupInviteRequest)
	r.GET("/contactList/blackList/:id", controllers.GetBlackList)
	r.POST("/contactList/blackList/remove/:id", controllers.RemoveFromBlacklist)
	r.GET("/contactList/friends/:id", controllers.GetFriends)
	r.GET("/contactList/groups/:id", controllers.GetGroups)
	r.POST("/contactList/CreateGroup/:id", controllers.CreateGroup)
	r.POST("/contactList/DeleteGroup/:id", controllers.DeleteGroup)
	r.GET("/contactList/profileCard/:id", controllers.ProfileCard)
	/*
		删除群聊路由
	*/

	// 个人主业路由
	r.GET("/Settings/profile/:id", controllers.ShowProfile)
	r.POST("/Settings/saveEdit/:id", controllers.SaveEdit)
	r.GET("/Settings/getPhone/:id", controllers.GetPhone)
	r.POST("/Settings/savePhone/:id", controllers.SavePhone)
	r.POST("/Settings/getCode/:id", controllers.GetCode)
	r.POST("/Settings/savePassword/:id", controllers.SavePassword)
	r.POST("/Settings/isIDAdd/:id", controllers.IsIDAdd)
	r.POST("/Settings/isPhoneAdd/:id", controllers.IsPhoneAdd)
	r.POST("/Settings/deactivation/:id", controllers.ConfirmDeactivation)
}
