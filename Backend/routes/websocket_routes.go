package routes

import (
	"TalkHive/controllers"
	"github.com/gin-gonic/gin"
)

func websocketRoutes(r *gin.Engine) {
	r.GET("/ws/websocketMessage", controllers.HandleConnections)
}
