package controllers

import (
	"context"
	"encoding/json"
	"file-sharing-backend/config"
	"file-sharing-backend/models"
	"file-sharing-backend/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// SearchFilesByName handles search requests for files by name
func SearchFilesByName(c *gin.Context) {
	// Extract user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt := userID.(int)

	// Get query parameters
	fileName := c.DefaultQuery("name", "")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File name is required"})
		return
	}

	// Create Redis cache key
	cacheKey := fmt.Sprintf("search:user:%d:name:%s", userIDInt, fileName)
	ctx := context.Background()

	// Check Redis cache
	if cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		var files []models.File
		if err := json.Unmarshal([]byte(cachedData), &files); err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Results from cache", "files": files})
			return
		}
	}

	// Fetch files from DB
	files, err := repositories.SearchFilesByName(userIDInt, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	// Cache the result for future requests
	filesJSON, _ := json.Marshal(files)
	_ = config.RedisClient.Set(ctx, cacheKey, filesJSON, 10*time.Minute).Err()

	c.JSON(http.StatusOK, gin.H{"message": "Search results", "files": files})
}

// SearchFilesBySize handles search requests for files by size
func SearchFilesBySize(c *gin.Context) {
	// Extract user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt := userID.(int)

	// Get file size parameter
	fileSizeStr := c.DefaultQuery("size", "")
	if fileSizeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size is required"})
		return
	}
	fileSize, err := strconv.Atoi(fileSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file size"})
		return
	}

	// Create Redis cache key
	cacheKey := fmt.Sprintf("search:user:%d:size:%d", userIDInt, fileSize)
	ctx := context.Background()

	// Check Redis cache
	if cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		var files []models.File
		if err := json.Unmarshal([]byte(cachedData), &files); err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Results from cache", "files": files})
			return
		}
	}

	// Fetch files from DB
	files, err := repositories.SearchFilesBySize(userIDInt, fileSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	// Cache the result for future requests
	filesJSON, _ := json.Marshal(files)
	_ = config.RedisClient.Set(ctx, cacheKey, filesJSON, 10*time.Minute).Err()

	c.JSON(http.StatusOK, gin.H{"message": "Search results", "files": files})
}

// SearchFilesByDate handles search requests for files by upload date
func SearchFilesByDate(c *gin.Context) {
	// Extract user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt := userID.(int)

	// Get upload date parameter
	uploadDate := c.DefaultQuery("date", "")
	if uploadDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Upload date is required"})
		return
	}

	// Create Redis cache key
	cacheKey := fmt.Sprintf("search:user:%d:date:%s", userIDInt, uploadDate)
	ctx := context.Background()

	// Check Redis cache
	if cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		var files []models.File
		if err := json.Unmarshal([]byte(cachedData), &files); err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Results from cache", "files": files})
			return
		}
	}

	// Fetch files from DB
	files, err := repositories.SearchFilesByDate(userIDInt, uploadDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	// Cache the result for future requests
	filesJSON, _ := json.Marshal(files)
	_ = config.RedisClient.Set(ctx, cacheKey, filesJSON, 10*time.Minute).Err()

	c.JSON(http.StatusOK, gin.H{"message": "Search results", "files": files})
}
// SearchFilesHandler handles file search requests based on name, size, and upload date
func SearchFilesHandler(c *gin.Context) {
	// Retrieve the user ID from JWT or context (you can assume the user is authenticated)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt := userID.(int)

	// Get query parameters for search
	fileName := c.DefaultQuery("name", "")
	fileSizeStr := c.DefaultQuery("size", "0")
	uploadDate := c.DefaultQuery("date", "")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Convert size, limit, and offset to integers
	fileSize, err := strconv.Atoi(fileSizeStr)
	if err != nil {
		fileSize = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Create cache key for Redis
	cacheKey := fmt.Sprintf("search:user:%d:name:%s:size:%d:date:%s:limit:%d:offset:%d", userIDInt, fileName, fileSize, uploadDate, limit, offset)
	ctx := context.Background()

	// Check Redis cache for results
	if cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		var files []models.File
		if err := json.Unmarshal([]byte(cachedData), &files); err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Results from cache", "files": files})
			return
		}
	}

	// Fetch files from the repository
	files, err := repositories.SearchFiles(userIDInt, fileName, fileSize, uploadDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	// Cache the results in Redis for 10 minutes
	filesJSON, _ := json.Marshal(files)
	_ = config.RedisClient.Set(ctx, cacheKey, filesJSON, 10*time.Minute).Err()

	// Return the results
	c.JSON(http.StatusOK, gin.H{"message": "Search results", "files": files})
}