package models

import "time"

type Document struct {
	ID         int32     `json:"id"`
	UserID     int       `json:"user_id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	UploadedAt time.Time `json:"uploaded_at"`
}
