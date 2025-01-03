package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	AuthRoutes(r)
	//ChatListRoutes(r)
	ContactListRoutes(r)
	SettingsRoutes(r)
	SetupWorkspaceRoutes(r)
}
