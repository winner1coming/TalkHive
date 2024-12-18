package main

import (
	"TalkHive/config"
	"TalkHive/routes"
	"fmt"
	"github.com/gin-contrib/cors"
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

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174"}, // 前端的URL
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "tid"}, // 允许的请求头中添加 tid
		AllowCredentials: true,
	}))

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务
	log.Println("监听端口:8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("初始化服务器失败: %v", err)
	}
}
