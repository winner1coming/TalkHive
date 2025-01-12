package routes

import (
	"TalkHive/controllers"

	"github.com/gin-gonic/gin"
)

func SetupGroupChatRoutes(router *gin.Engine) *gin.RouterGroup {
	r := router.Group("/contactList")

	//r.GET("/groups", workSpace.GetGroups)                                 // 获取群聊列表
	//r.POST("/groups/createGroup", workSpace.CreateGroup)                  //创建群聊
	//r.POST("/groups/dismissGroup", workSpace.DismissGroup)                // 解散群聊
	//r.POST("/groups/friendsNotInGroup", workSpace.FetchFriendsNotInGroup) // 搜索不在群聊内的好友
	//r.POST("/groups/invite", workSpace.InviteMember)                      // 邀请他人入群
	//r.POST("/groups/exit", workSpace.ExitGroup)                           // 退出群聊
	//r.POST("/groups/groupInfo/:group_id", workSpace.GetGroupInfo)         // 获取群聊信息
	//r.POST("/groups/changeNickName", workSpace.ChangeNickName)            // 更改我在群内的昵称
	// r.POST("/groups/setAllowInvite", controllers.SetAllowInvite)         // 设置是否允许群成员邀请他
	r.POST("/groups/setAllowSearch", controllers.SetAllowSearch) // 设置是否允许群成员邀请他：ID
	// r.POST("/groups/setAllowNameSearch", controllers.SetAllowNameSearch) // 设置是否允许群成员邀请他：名字
	// r.POST("/groups/setAllBanned", controllers.SetAllBanned)             // 全体禁言
	// r.POST("/groups/changeGroupAvatar", controllers.ChangeGroupAvatar)   // 更改群头像

	return r
}
