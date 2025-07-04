package usecase

import (
	"context"
	"Complaingo/internal/domain/models"
	appErrors "Complaingo/internal/errors"
	"Complaingo/internal/middleware"
	"Complaingo/internal/notifier"
	"Complaingo/internal/rabbitmq"
	"Complaingo/internal/repository"
	"encoding/json"
	"fmt"
	"time"

	"github.com/joomcode/errorx"
)

type ComplaintUsecase struct {
	complaintRepo repository.ComplaintRepository
	messageRepo   repository.ComplaintMessageRepository
	notifier      notifier.Notifier
}

func NewComplaintUsecase(cr repository.ComplaintRepository, cm repository.ComplaintMessageRepository, n notifier.Notifier) *ComplaintUsecase {
	return &ComplaintUsecase{
		complaintRepo: cr,
		messageRepo:   cm,
		notifier:      n,
	}
}

func (cr *ComplaintUsecase) CreateComplaint(ctx context.Context, c *models.Complaints) error {
	if err := cr.complaintRepo.CreateComplaint(ctx, c); err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserDuplicate) {
			return err
		}
		return appErrors.ErrDbFailure.Wrap(err, "usecase: unable to create user")
	}

	// publish to rabbitmq
	message := models.NotificationMessage{
		Type:      "complaint_created",
		UserID:    c.UserID,
		Complient: c.Subject,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// serialize and send
	messageJson, _ := json.Marshal(message)
	prod := rabbitmq.NewProducer("amqp://guest:guest@localhost:5672/", "notifications")
	prod.SendMessage(string(messageJson))

	return nil
}

func (cr *ComplaintUsecase) GetComplaintByRole(ctx context.Context, UserID int) ([]*models.Complaints, error) {
	complaints, err := cr.complaintRepo.GetComplaintByRole(ctx, UserID)
	if err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.ErrUnauthorized.Wrap(err, "usecase: complaint not found")
		}
		return nil, appErrors.ErrDbFailure.Wrap(err, "usecase: unexpected db")
	}
	return complaints, nil
}

func (cr *ComplaintUsecase) UserMarkResolved(ctx context.Context, complaintID int) error {
	userID := middleware.GetUserId(ctx)
	fmt.Printf("🔍 Logged-in user ID: %d\n", userID)
	complaint, err := cr.complaintRepo.GetComplaintByID(ctx, complaintID)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "complaint not found")
	}

	fmt.Printf("📦 Complaint fetched: %+v\n", complaint)
	if complaint.UserID != userID {
		return appErrors.ErrDbFailure.New("user can only update their own complaint")
	}

	fmt.Printf("Updating complaint ID %d to Resolved\n", complaintID)
	return cr.complaintRepo.UpdateComplaints(ctx, complaintID, "Resolved")
}

func (cr *ComplaintUsecase) GetAllComplaintByRole(ctx context.Context) ([]*models.Complaints, error) {
	return cr.complaintRepo.GetAllComplaintByRole(ctx)
}

func (cr *ComplaintUsecase) AdminUpdateComplaints(ctx context.Context, complaintID int, status string) error {
	validStatus := map[string]bool{
		"Created":  true,
		"Accepted": true,
		"Resolved": true,
		"Rejected": true,
	}
	if !validStatus[status] {
		return appErrors.ErrInvalidPayload.New("Invalid complaint status")
	}

	return cr.complaintRepo.UpdateComplaints(ctx, complaintID, status)
}

// complaint_messages table
func (cr *ComplaintUsecase) InsertCoplaintMessage(ctx context.Context, cm *models.ComplaintMessages) error {
	if err := cr.messageRepo.InsertCoplaintMessage(ctx, cm); err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserDuplicate) {
			return err
		}
		return appErrors.ErrDbFailure.Wrap(err, "usecase: unable to create user")
	}

	return nil
}

func (cr *ComplaintUsecase) ReplyToMessage(ctx context.Context, msg *models.ComplaintMessages) error {
	if msg.ParentID != nil {
		parentMsg, err := cr.messageRepo.GetMessageByID(ctx, *msg.ParentID)
		if err != nil {
			return appErrors.ErrUserNotFound.Wrap(err, "parent message not found")
		}
		if parentMsg.ComplaintID != msg.ComplaintID {
			return appErrors.ErrInvalidPayload.New("parent message must belongs to the same complaint")
		}
	}

	role := middleware.GetUserRole(ctx)

	if role == "admin" && msg.FileUrl != "" {
		return appErrors.ErrUnauthorized.New("admins are not allowed to attach files")
	}

	msg.SenderID = middleware.GetUserId(ctx)
	err := cr.messageRepo.AddMessage(ctx, msg)
	if err != nil {
		return appErrors.ErrDbFailure.Wrap(err, "failed to save reply")
	}

	if role == "user" {
		cr.notifier.SendToAdmins(msg)
	}
	if role == "admin" {
		cr.notifier.SendToUser(msg.SenderID, msg)
	}

	return nil
}

func (cr *ComplaintUsecase) GetMessagesByComplaint(ctx context.Context, complaintID int) ([]*models.ComplaintMessages, error) {
	complaints, err := cr.messageRepo.GetMessagesByComplaint(ctx, complaintID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound.Wrap(err, "usecase: message not found")
	}

	return complaints, nil
}
