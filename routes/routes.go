package routes

import (
	"file-sharing-backend/controllers"
	"file-sharing-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}

	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware()) 
	{
		protected.POST("/upload", controllers.UploadFile) 
		protected.GET("/profile", controllers.Profile) 
		protected.GET("/files", controllers.GetFiles)
		protected.GET("/share/:file_id", controllers.ShareFile)
		protected.GET("search", controllers.SearchFilesHandler)
		protected.GET("search/name", controllers.SearchFilesByName)
		protected.GET("search/size", controllers.SearchFilesBySize)
		protected.GET("search/date", controllers.SearchFilesByDate)// Search API	
	}
}
