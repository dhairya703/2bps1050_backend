package controllers

import (
	"context"
	"fmt"
	"file-sharing-backend/config"
	"file-sharing-backend/models"
	"file-sharing-backend/repositories"
	"github.com/gin-gonic/gin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"net/http"
	"mime/multipart"
	"strconv"
	"time"
)

// FileUploadResult stores the result of the upload operation
type FileUploadResult struct {
	URL      string
	FileName string
	Size     int64
	Err      error
}

// UploadFile handles file uploads concurrently
func UploadFile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDInt, _ := strconv.Atoi(fmt.Sprintf("%v", userID)) // Convert to int

	// Parse multipart form (max 100MB)
	err := c.Request.ParseMultipartForm(100 << 20) // 100MB limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large or invalid format"})
		return
	}

	form := c.Request.MultipartForm
	files := form.File["file"] // Get uploaded files

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Channel to receive upload results
	resultChan := make(chan FileUploadResult, len(files))

	// Upload each file concurrently
	for _, file := range files {
		go func(fileHeader *multipart.FileHeader) {
			uploadURL, err := processFileUpload(c, fileHeader)
			resultChan <- FileUploadResult{URL: uploadURL, FileName: fileHeader.Filename, Size: fileHeader.Size, Err: err}
		}(file)
	}

	// Collect results & store metadata
	var uploadedFiles []string
	for range files {
		result := <-resultChan
		if result.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}

		// Set the expiry date as 30 days after upload
		expiryDate := time.Now().Add(30 * 24 * time.Hour)

		// Store metadata in the database
		fileRecord := models.File{
			UserID:       userIDInt,
			FileName:     result.FileName,
			Size:         result.Size,
			CloudinaryURL: result.URL,
			UploadDate:    time.Now(),
			ExpiryDate:    expiryDate,
		}

		err := repositories.SaveFileMetadata(&fileRecord)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store metadata"})
			return
		}

		uploadedFiles = append(uploadedFiles, result.URL)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files uploaded & stored successfully", "urls": uploadedFiles})
}

// processFileUpload uploads a file to Cloudinary
func processFileUpload(c *gin.Context, fileHeader *multipart.FileHeader) (string, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Upload to Cloudinary
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadResult, err := config.Cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
