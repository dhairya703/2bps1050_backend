package repositories

import (
	"file-sharing-backend/config"
	"file-sharing-backend/models"
	"log"
)

func SaveFileMetadata(file *models.File) error {
    query := `INSERT INTO files (user_id, file_name, upload_date, size, cloudinary_url,expiry_date) 
              VALUES ($1, $2, $3, $4, $5,$6)`

    _, err := config.DB.Exec(query, file.UserID, file.FileName, file.UploadDate, file.Size, file.CloudinaryURL,file.ExpiryDate)
    return err
}
// GetFilesByUser fetches all files uploaded by a user
func GetFilesByUser(userID int) ([]models.File, error) {
	query := "SELECT id, user_id, file_name, upload_date, size, cloudinary_url,expiry_date FROM files WHERE user_id = $1"
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		err := rows.Scan(&file.ID, &file.UserID, &file.FileName, &file.UploadDate, &file.Size, &file.CloudinaryURL)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// GetFileByID retrieves a file by ID
func GetFileByID(fileID int) (*models.File, error) {
	query := "SELECT id, user_id, file_name, upload_date, size, cloudinary_url,expiry_date FROM files WHERE id = $1"
	row := config.DB.QueryRow(query, fileID)

	var file models.File
	err := row.Scan(&file.ID, &file.UserID, &file.FileName, &file.UploadDate, &file.Size, &file.CloudinaryURL)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
// DeleteFileMetadata deletes the metadata of a file from the database
func DeleteFileMetadata(fileID int) error {
	// SQL query to delete file metadata by ID
	query := "DELETE FROM files WHERE id = $1"
	_, err := config.DB.Exec(query, fileID) // Executes the delete operation
	if err != nil {
		log.Println("Error deleting file metadata:", err)
		return err
	}
	return nil
}