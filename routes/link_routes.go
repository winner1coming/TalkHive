package routes

func LinksRoutes(r *gin.Engine) {
	r.POST("/addLinks", controllers.AddLinks)
	r.GET("/delLinks", controllers.DelLinks)
}
