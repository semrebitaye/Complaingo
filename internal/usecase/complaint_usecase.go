package usecase

import (
	"context"
	"crud_api/internal/domain/models"
	appErrors "crud_api/internal/errors"
	"crud_api/internal/middleware"
	"crud_api/internal/repository"
	"fmt"

	"github.com/joomcode/errorx"
)

type ComplaintUsecase struct {
	complaintRepo repository.ComplaintRepository
	messageRepo   repository.ComplaintMessageRepository
}

func NewComplaintUsecase(cr repository.ComplaintRepository, cm repository.ComplaintMessageRepository) *ComplaintUsecase {
	return &ComplaintUsecase{
		complaintRepo: cr,
		messageRepo:   cm,
	}
}

func (cr *ComplaintUsecase) CreateComplaint(ctx context.Context, c *models.Complaints) error {
	if err := cr.complaintRepo.CreateComplaint(ctx, c); err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserDuplicate) {
			return err
		}
		return appErrors.ErrDbFailure.Wrap(err, "usecase: unable to create user")
	}

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
		fmt.Println("❌ Error fetching complaint:", err)
		return appErrors.ErrDbFailure.Wrap(err, "complaint not found")
	}

	fmt.Printf("📦 Complaint fetched: %+v\n", complaint)
	if complaint.UserID != userID {
		fmt.Printf("❌ User ID %d does not own complaint %d (owned by %d)\n", userID, complaintID, complaint.UserID)
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
	return cr.messageRepo.AddMessage(ctx, msg)
}

func (cr *ComplaintUsecase) GetMessagesByComplaint(ctx context.Context, complaintID int) ([]*models.ComplaintMessages, error) {
	complaints, err := cr.messageRepo.GetMessagesByComplaint(ctx, complaintID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound.Wrap(err, "usecase: message not found")
	}

	return complaints, nil
}
