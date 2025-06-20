package models

import "time"

type Complaints struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
