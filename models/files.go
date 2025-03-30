package models

import "time"

type File struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	FileName      string    `json:"file_name"`
	UploadDate    time.Time `json:"upload_date"`
	ExpiryDate    time.Time `json:"upload_date"`
	Size          int64     `json:"size"`
	CloudinaryURL string    `json:"cloudinary_url"`
}
