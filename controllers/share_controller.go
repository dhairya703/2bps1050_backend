package controllers

import (
	"context"
	"file-sharing-backend/config"
	"file-sharing-backend/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// ShareFile generates a public link for a file (stores in Redis with expiration)
func ShareFile(c *gin.Context) {
	fileIDStr := c.Param("file_id")

	// Convert fileID to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Check Redis cache first
	redisKey := "shared:file:" + fileIDStr
	ctx := context.Background()

	cachedURL, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil { // Cache hit, return cached URL
		c.JSON(http.StatusOK, gin.H{"message": "File URL retrieved from cache", "shared_url": cachedURL})
		return
	}

	// Cache miss, fetch from DB
	file, err := repositories.GetFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Store shared URL in Redis (expires in 24 hours)
	sharedURL := file.CloudinaryURL
	_ = config.RedisClient.Set(ctx, redisKey, sharedURL, 24*time.Hour).Err()

	c.JSON(http.StatusOK, gin.H{"message": "File shared successfully", "shared_url": sharedURL})
}
