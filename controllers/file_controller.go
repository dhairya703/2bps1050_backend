package controllers

import (
	"context"
	"encoding/json"
	"file-sharing-backend/config"
	"file-sharing-backend/models"
	"file-sharing-backend/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetFiles retrieves metadata of uploaded files (with Redis caching)
func GetFiles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert userID to int
	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Generate Redis cache key
	redisKey := "files:user:" + strconv.Itoa(userIDInt)
	ctx := context.Background()

	// Check Redis cache first
	cachedData, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil { // Cache hit, return cached data
		var files []models.File
		if jsonErr := json.Unmarshal([]byte(cachedData), &files); jsonErr == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Files retrieved from cache", "files": files})
			return
		}
	}

	// Cache miss, fetch from DB
	files, err := repositories.GetFilesByUser(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	// Store result in Redis (cache for 5 minutes)
	filesJSON, jsonErr := json.Marshal(files)
	if jsonErr == nil { // Only cache if JSON encoding succeeds
		_ = config.RedisClient.Set(ctx, redisKey, filesJSON, 5*time.Minute).Err()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files retrieved", "files": files})
}
