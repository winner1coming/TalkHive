package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

// 路由会把 HTTP 请求指向不同的控制器函数
// 例如 controllers.Register、controllers.Login
// 这些控制器函数将包含具体的业务逻辑，比如数据库操作、数据验证、响应生成等
func SetupRouter_testing() *gin.Engine {
	r := gin.Default()

	// 用户账号相关
	// 路由路径 处理函数
	r.POST("/register", controllers.Register)               // 用户注册
	r.POST("/login", controllers.Login)                     // 用户登录
	r.GET("/profile/:id", controllers.GetProfile)           // 获取用户资料
	r.PUT("/profile/update/:id", controllers.UpdateProfile) // 更新用户资料

	// 好友相关
	r.POST("/friend/add", controllers.AddFriend) // 添加好友
	// 你可以添加更多的好友操作，如删除好友、查看好友列表等

	// 群聊相关
	r.GET("/group/:id", controllers.GetMessages)     // 获取群聊信息
	r.POST("/group/create", controllers.CreateGroup) // 创建群聊
	// 你可以根据需求继续添加更多群聊相关操作，如删除群聊、查看群聊成员等

	// 消息相关
	r.GET("/messages/:id", controllers.GetMessages)   // 获取消息记录
	r.POST("/messages/send", controllers.SendMessage) // 发送消息
	// 可以添加更多的消息相关功能，如删除消息、查看特定时间段的消息等

	// 申请通知相关
	r.POST("/apply/submit", controllers.SubmitApply) // 提交申请
	r.PUT("/apply/review", controllers.ReviewApply)  // 审核申请

	// 系统设置相关
	r.PUT("/settings/update", controllers.UpdateSystemSetting) // 更新系统设置

	// 其他路由...
	return r
}
