package main

import (
	"chatroom/config"
	"chatroom/routes"
	"log"
)

func main() {
	// 初始化数据库
	configPath := "config/config.yaml"
	config.InitDB(configPath)

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
