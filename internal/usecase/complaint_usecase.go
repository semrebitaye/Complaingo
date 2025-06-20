package usecase

import (
	"context"
	"crud_api/internal/domain/models"
	appErrors "crud_api/internal/errors"
	"crud_api/internal/repository"
	contexthelper "crud_api/internal/utility/context_helper"

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
	// authorization check
	if contexthelper.IsAdmin(ctx) {
		return appErrors.ErrUnauthorized.New("only users can create complaints")
	}

	if err := cr.complaintRepo.CreateComplaint(ctx, c); err != nil {
		if errorx.IsOfType(err, appErrors.ErrUserDuplicate) {
			return err
		}
		return appErrors.ErrDbFailure.Wrap(err, "usecase: unable to create user")
	}

	return nil
}

func (cr *ComplaintUsecase) GetComplaintByRole(ctx context.Context, UserID int) ([]*models.Complaints, error) {
	if contexthelper.IsAdmin(ctx) {
		return nil, appErrors.ErrUnauthorized.New("this is for only users")
	}

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
	if contexthelper.IsAdmin(ctx) {
		return appErrors.ErrUnauthorized.New("this is only for users")
	}

	userID := contexthelper.GetUserId(ctx)
	complaint, err := cr.complaintRepo.GetComplaintByID(&ctx, complaintID)
	if err != nil {
		appErrors.ErrDbFailure.Wrap(err, "complaint not found")
	}

	if complaint.UserID != userID {
		return appErrors.ErrDbFailure.New("user can only update their own complaint")
	}

	return cr.complaintRepo.UpdateComplaints(ctx, complaintID, "Resolved")
}

func (cr *ComplaintUsecase) GetAllComplaintByRole(ctx context.Context) ([]*models.Complaints, error) {
	if !contexthelper.IsAdmin(ctx) {
		return nil, appErrors.ErrUnauthorized.New("this operations is only for admins")
	}
	return cr.complaintRepo.GetAllComplaintByRole(ctx)
}

func (cr *ComplaintUsecase) AdminUpdateComplaints(ctx context.Context, complaintID int, status string) error {
	if !contexthelper.IsAdmin(ctx) {
		return appErrors.ErrUnauthorized.New("only admins have permits to do this operation")
	}

	if status != "Accepted" && status != "Resolved" && status != "Rejected" {
		return appErrors.ErrInvalidPayload.New("status not valid")
	}

	return cr.complaintRepo.UpdateComplaints(ctx, complaintID, status)
}

// complaint_messages table
func (cr *ComplaintUsecase) ReplayToMessage(ctx context.Context, msg *models.ComplaintMessages) error {
	if msg.ParentID != nil {
		parentMsg, err := cr.messageRepo.GetMessageByID(ctx, *msg.ParentID)
		if err != nil {
			return appErrors.ErrUserNotFound.Wrap(err, "parent message not found")
		}
		if parentMsg.ComplaintID != msg.ComplaintID {
			return appErrors.ErrInvalidPayload.New("parent message must belongs to the same complaint")
		}
	}

	role := contexthelper.GetUserRole(ctx)

	if role == "admin" && msg.FileUrl != nil {
		return appErrors.ErrUnauthorized.New("admins are not allowed to attach files")
	}

	msg.SenderID = contexthelper.GetUserId(ctx)
	return cr.messageRepo.AddMessage(ctx, msg)
}

func (cr *ComplaintUsecase) GetMessagesByComplaint(ctx context.Context, complaintID int) ([]*models.ComplaintMessages, error) {
	complaints, err := cr.messageRepo.GetMessagesByComplaint(ctx, complaintID)
	if err != nil {
		return nil, appErrors.ErrUserNotFound.Wrap(err, "usecase: message not found")
	}

	return complaints, nil
}
