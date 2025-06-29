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

type MessageEntity struct {
	ID         int     `json:"id"`
	FromUserID int     `json:"from_user_id"`
	ToUserID   *int    `json:"to_user_id"`
	ToRole     *string `json:"to_role"`
	Channel    *string `json:"channel"`
	Message    string  `json:"message"`
	CeatedAt   string  `json:"created_at"`
}
