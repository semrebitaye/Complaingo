package repository

import (
	"context"
	"crud_api/internal/domain/models"
)

type ComplaintRepository interface {
	CreateComplaint(ctx context.Context, c *models.Complaints) error                  //user only
	GetComplaintByRole(ctx context.Context, UserID int) ([]*models.Complaints, error) // user only
	GetComplaintByID(ctx context.Context, complaintID int) (*models.Complaints, error)
	UpdateComplaints(ctx context.Context, complaintID int, status string) error
	GetAllComplaintByRole(ctx context.Context) ([]*models.Complaints, error) //admin olny
}

type ComplaintMessageRepository interface {
	InsertCoplaintMessage(ctx context.Context, cm *models.ComplaintMessages) error
	AddMessage(ctx context.Context, cm *models.ComplaintMessages) error
	GetMessageByID(ctx context.Context, messageID int) (*models.ComplaintMessages, error)
	GetMessagesByComplaint(ctx context.Context, complaintID int) ([]*models.ComplaintMessages, error)
}
