package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

// ContactListRoutes 通讯录路由
func ContactListRoutes(r *gin.Engine) {
	// 添加、搜索好友
	r.POST("/stranger/search", controllers.SearchStrangers)
	r.POST("/contactList/remark", controllers.ChangeRemark)
	r.POST("/contactList/search", controllers.SearchContacts)

	//好友请求部分
	r.GET("/contactList/friendRequests", controllers.GetFriendRequests)
	r.POST("/contactList/friendRequests/pend", controllers.FriendRequestPend)
	r.POST("/contactList/friendRequests/addFriend", controllers.AddFriend)

	//群聊请求部分
	r.GET("/contactList/groupRequests", controllers.GetGroupRequests)
	r.POST("/contactList/groupRequests/applyPend", controllers.DealGroupApplyRequest)
	r.POST("/contactList/groupRequests/invitationPend", controllers.DealGroupInviteRequest)
	r.POST("/contactList/groupRequests/addGroup", controllers.AddGroup)

	//黑名单部分
	r.GET("/contactList/blackList", controllers.GetBlackList)
	r.POST("/contactList/blackList/remove", controllers.RemoveFromBlacklist)
	r.POST("/contactList/blackList/add", controllers.AddToBlacklist)

	// 好友列表
	r.GET("/contactList/friends", controllers.GetFriends)
	r.POST("/contactList/friends/delete", controllers.DeleteFriend)

	// 分组部分
	r.GET("/contactList/:type/divides", controllers.GetDivides)
	r.POST("/contactList/:type/divides/create", controllers.CreateDivide)
	r.POST("/contactList/:type/divides/delete", controllers.DeleteDivide)
	r.POST("/contactList/:type/divides/rename", controllers.RenameDivide)
	r.POST("/contactList/:type/divides/moveIn", controllers.MoveInDivide)

	// 群聊部分
	r.GET("/contactList/groups", controllers.GetGroups)
	r.POST("/contactList/groups/createGroup", controllers.CreateGroup)
	r.POST("/contactList/groups/dismissGroup", controllers.DisMissGroup)
	r.POST("/contactList/groups/friendsNotInGroup", controllers.FetchFriendsNotInGroup)
	r.POST("/contactList/groups/invite", controllers.Invite)
	r.POST("/contactList/groups/quit", controllers.Quit)
	r.GET("/contactList/groups/groupInfo/:group_id", controllers.GetGroupInfo)
	r.POST("/contactList/groups/changeNickname", controllers.ChangeNickname)
	r.POST("/contactList/groups/setAllowInvite", controllers.SetAllowInvite)
	r.POST("/contactList/groups/setAllowIdSearch", controllers.SetAllowIdSearch)
	r.POST("/contactList/groups/setAllowNameSearch", controllers.SetAllowNameSearch)
	r.POST("/contactList/groups/setAllBanned", controllers.SetAllBanned)
	r.POST("/contactList/groups/banMember", controllers.SetBanned)
	r.POST("/contactList/groups/removeMember", controllers.RemoveMember)
	r.POST("/contactList/groups/setAdmin", controllers.SetAdmin)
	r.POST("/contactList/groups/transferOwner", controllers.TransferOwner)
	r.POST("/contactList/groups/changeAvatar", controllers.ChangeGroupAvatar)
	r.POST("/contactList/groups/changeName", controllers.ChangeGroupName)
	r.POST("/contactList/groups/changeGroupIntroduction", controllers.ChangeGroupIntroduction)
	// 资料卡片
	r.POST("/profileCard/person", controllers.GetPersonProfileCard)
	r.POST("/profileCard/group", controllers.GetGroupProfileCard)
}
