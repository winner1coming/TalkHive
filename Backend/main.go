package main

import (
	"chatroom/config"
	"chatroom/routes"
	"log"
)

func main() {
	// 初始化全局配置，包括数据库和 Redis
	config.InitConfig()
	config.InitDB()
	config.InitRedis()

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
