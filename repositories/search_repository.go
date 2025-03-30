package repositories

import (
	"file-sharing-backend/config"
	"file-sharing-backend/models"
	"strings"
	"strconv"
)

// SearchFilesByName retrieves files for the user based on the file name
func SearchFilesByName(userID int, fileName string) ([]models.File, error) {
	var files []models.File
	query := "SELECT id, user_id, file_name, upload_date, size, cloudinary_url FROM files WHERE user_id = $1 AND file_name LIKE $2"
	rows, err := config.DB.Query(query, userID, "%"+fileName+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// SearchFilesBySize retrieves files for the user based on the file size
func SearchFilesBySize(userID int, fileSize int) ([]models.File, error) {
	var files []models.File
	query := "SELECT id, user_id, file_name, upload_date, size, cloudinary_url FROM files WHERE user_id = $1 AND size = $2"
	rows, err := config.DB.Query(query, userID, fileSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// SearchFilesByDate retrieves files for the user based on the upload date
func SearchFilesByDate(userID int, uploadDate string) ([]models.File, error) {
	var files []models.File
	query := "SELECT id, user_id, file_name, upload_date, size, cloudinary_url FROM files WHERE user_id = $1 AND upload_date = $2"
	rows, err := config.DB.Query(query, userID, uploadDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
// SearchFiles retrieves files for the user based on name, size, or upload date with pagination
func SearchFiles(userID int, fileName string, fileSize int, uploadDate string, limit, offset int) ([]models.File, error) {
	var files []models.File
	var queryParts []string
	var args []interface{}
	query := "SELECT id, user_id, file_name, upload_date, size, cloudinary_url FROM files WHERE user_id = $1"
	args = append(args, userID)

	// Add filters to query dynamically
	if fileName != "" {
		queryParts = append(queryParts, "file_name ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, "%"+fileName+"%")
	}
	if fileSize > 0 {
		queryParts = append(queryParts, "size = $"+strconv.Itoa(len(args)+1))
		args = append(args, fileSize)
	}
	if uploadDate != "" {
		queryParts = append(queryParts, "upload_date = $"+strconv.Itoa(len(args)+1))
		args = append(args, uploadDate)
	}

	// If filters were added, append them to the base query
	if len(queryParts) > 0 {
		query += " AND " + strings.Join(queryParts, " AND ")
	}

	// Add pagination (limit and offset)
	query += " LIMIT $"+strconv.Itoa(len(args)+1) + " OFFSET $"+strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	// Execute the query
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results into the files slice
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