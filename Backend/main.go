package main

import (
	"TalkHive/config"
	"TalkHive/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	// 初始化全局配置，包括数据库和Redis
	fmt.Println("初始化全局配置...")
	config.InitConfig()
	config.InitDB()
	config.InitRedis()

	// 运行路由测试
	//log.Println("开始运行路由测试...")
	//tests.RunRouteTests(r)

	// 初始化 Gin 引擎
	r := gin.Default()

	// 设置路由分组
	routes.SetupAddressBookRoutes(r)  // 通讯录模块
	routes.SetupGroupChatRoutes(r)    // 群聊模块
	routes.SetupMessageListRoutes(r)  // 消息列表模块
	routes.SetupPersonalHomeRoutes(r) // 个人主页模块
	routes.SetupRouter(r)             // 用户账号相关模块
	routes.SetupSingleChatRoutes(r)   // 单聊模块
	routes.SetupWorkspaceRoutes(r)    // 工作区模块

	// 启动服务
	log.Println("监听端口:8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("初始化服务器失败: %v", err)
	}
}
