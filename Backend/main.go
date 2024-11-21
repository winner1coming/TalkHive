package main

import (
	"chatroom/config"
	"chatroom/routes"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务
	r.Run(":8080")
}
