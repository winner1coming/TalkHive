package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func SetupAddressBookRoutes(router *gin.Engine) *gin.RouterGroup {
	// 创建通讯录分组
	r := router.Group("/addressBook")

	// 发送好友申请/入群申请
	r.POST("/sendApplication", controllers.SendApplication)

	// 接受好友申请/通过好友申请
	r.POST("/acceptApplication", controllers.AcceptApplication)

	// 删除好友/被删除
	r.DELETE("/removeFriend", controllers.RemoveFriend)

	// 获取好友申请状态
	r.GET("/applicationStatus", controllers.GetApplicationStatus)

	// 入群申请通过
	r.POST("/acceptGroupApplication", controllers.AcceptGroupApplication)

	// 退出群聊/被踢出群聊
	r.POST("/leaveGroup", controllers.LeaveGroup)

	// 获取群聊申请状态
	r.GET("/groupApplicationStatus", controllers.GetGroupApplicationStatus)

	// 添加分组
	r.POST("/addGroup", controllers.AddGroup)

	// 删除分组
	r.DELETE("/deleteGroup", controllers.DeleteGroup)

	// 重命名分组
	r.PUT("/renameGroup", controllers.RenameGroup)

	// 获取分组人数统计
	r.GET("/groupStatistics", controllers.GetGroupStatistics)

	// 获取总好友数和群聊数
	r.GET("/summary", controllers.GetSummary)

	// 搜索好友或群聊
	r.GET("/search", controllers.SearchContact)

	return r
}
