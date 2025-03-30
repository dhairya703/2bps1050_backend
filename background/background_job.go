package background

import (
	"context"
	"file-sharing-backend/config"
	"file-sharing-backend/repositories"
	"log"
	"strings"
	"time"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// DeleteExpiredFilesJob runs periodically to delete expired files
func DeleteExpiredFilesJob() {
	ticker := time.NewTicker(24 * time.Hour) // Runs once every 24 hours
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Get the current time to check against file expiry
			currentTime := time.Now()

			// Fetch expired files from the database
			files, err := repositories.GetExpiredFiles(currentTime)
			if err != nil {
				log.Println("Error fetching expired files:", err)
				continue
			}

			// Process each expired file
			for _, file := range files {
				// Delete the file from Cloudinary
				err := deleteFileFromCloudinary(file.CloudinaryURL)
				if err != nil {
					log.Println("Error deleting file from Cloudinary:", err)
					continue
				}

				// Delete the corresponding metadata from PostgreSQL
				err = repositories.DeleteFileMetadata(file.ID)
				if err != nil {
					log.Println("Error deleting file metadata from DB:", err)
					continue
				}

				log.Printf("Successfully deleted expired file %d: %s\n", file.ID, file.FileName)
			}
		}
	}
}

// deleteFileFromCloudinary removes the file from Cloudinary using the URL
func deleteFileFromCloudinary(url string) error {
	// Extract the public ID from the Cloudinary URL
	publicID := extractPublicIDFromURL(url)

	// Delete the file from Cloudinary
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := config.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// extractPublicIDFromURL extracts the public ID from Cloudinary URL
func extractPublicIDFromURL(url string) string {
	// Split the URL to get the public ID
	segments := strings.Split(url, "/")
	publicIDWithExtension := segments[len(segments)-1]
	publicID := strings.Split(publicIDWithExtension, ".")[0]
	return publicID
}
