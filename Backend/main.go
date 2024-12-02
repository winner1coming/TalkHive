package main

import (
	"TalkHive/config"
	"TalkHive/routes"
	"TalkHive/test"
	"log"
)

func main() {
	// 初始化全局配置，包括数据库和Redis
	config.InitConfig()
	config.InitDB()
	config.InitRedis()

	// 设置路由
	r := routes.SetupRouter()

	// 运行路由测试
	log.Println("开始运行路由测试...")
	tests.RunRouteTests(r)

	// 启动服务
	log.Println("监听端口:8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("初始化服务器失败: %v", err)
	}
}
