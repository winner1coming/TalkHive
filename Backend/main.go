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

	// 初始化 Gin 引擎
	r := gin.Default()
	routes.SetupRoutes(r)

	// 启动服务
	log.Println("监听端口:8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("初始化服务器失败: %v", err)
	}
}
