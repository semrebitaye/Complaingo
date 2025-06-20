package models

import "time"

type ComplaintMessages struct {
	ID          int
	ComplaintID int
	SenderID    int
	ParentID    *int
	Message     string
	FileUrl     *string
	CreatedAt   time.Time
}
