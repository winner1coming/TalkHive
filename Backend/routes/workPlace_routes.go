package routes

import (
	"TalkHive/controllers"

	"github.com/gin-gonic/gin"
)

func workSpaceRoutes(r *gin.Engine) {
	r.POST("/addLinks", controllers.AddLinks)
	r.POST("/delLinks", controllers.DelLinks)
}
