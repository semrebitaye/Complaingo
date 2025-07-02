package models

type NotificationMessage struct {
	Type      string `json:"type"`
	UserID    int    `json:"user_id"`
	Complient string `json:"complient"`
	Timestamp string `json:"timestamp"`
}
