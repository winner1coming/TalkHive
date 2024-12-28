package routes

import "github.com/gin-gonic/gin"

// SetupRoutes 集中设置所有路由
func SetupRoutes(r *gin.Engine) {
	AuthRoutes(r)
	ContactListRoutes(r)
	SettingsRoutes(r)
}
