package main

import (
	"file-sharing-backend/middlewares"

	// "time"
	"file-sharing-backend/config"
	"file-sharing-backend/routes"
	// "file-sharing-backend/background"
	"github.com/gin-gonic/gin"
		"log"
		"fmt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Error loading .env file: %v", err)
	}
	// Connect to database before starting server

	config.ConnectDatabase()
	config.InitCloudinary()
	config.InitRedis()
	// go background.DeleteExpiredFilesJob()
	// time.Sleep(1 * time.Minute)

	// Add any other logic to check if expired files were deleted (e.g., querying the database)
	// You can run a manual query here to confirm the files have been deleted

	fmt.Println("Finished testing background job.")
	router := gin.Default()
	router.Use(middlewares.AuthMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "File Sharing API is running!"})
	})

	routes.SetupRoutes(router)

	router.Run(":8080")
}
