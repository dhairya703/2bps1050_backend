package repositories

import (
	"file-sharing-backend/config"
	"file-sharing-backend/models"
	"time"
)

// GetExpiredFiles fetches all files that have expired
func GetExpiredFiles(currentTime time.Time) ([]models.File, error) {
	// SQL query to fetch files where expiry_date is less than the current time
	query := "SELECT id, user_id, file_name, upload_date, size, cloudinary_url, expiry_date FROM files WHERE expiry_date < $1"
	rows, err := config.DB.Query(query, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		// Scan the row values into the file struct
		err := rows.Scan(&file.ID, &file.UserID, &file.FileName, &file.UploadDate, &file.Size, &file.CloudinaryURL, &file.ExpiryDate)
		if err != nil {
			return nil, err
		}
		// Append the file to the slice
		files = append(files, file)
	}

	// Return the list of expired files
	return files, nil
}
