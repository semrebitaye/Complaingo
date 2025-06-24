package models

import "time"

type ComplaintMessages struct {
	ID          int       `json:"id"`
	ComplaintID int       `json:"complaint_id"`
	SenderID    int       `json:"sender_id"`
	ParentID    *int      `json:"parent_id,omitempty"`
	Message     string    `json:"message"`
	FileUrl     string    `json:"file_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
