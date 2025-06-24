package handler

import (
	"crud_api/internal/domain/models"
	"crud_api/internal/middleware"
	"crud_api/internal/usecase"
	"crud_api/internal/utility"
	"encoding/json"
	"net/http"
	"strconv"

	appErrors "crud_api/internal/errors"

	"github.com/gorilla/mux"
)

type ComplaintHandler struct {
	usecase *usecase.ComplaintUsecase
}

func NewComplaintHandler(usecase *usecase.ComplaintUsecase) *ComplaintHandler {
	return &ComplaintHandler{usecase: usecase}
}

func (uc *ComplaintHandler) CreateComplaint(w http.ResponseWriter, r *http.Request) {
	var c models.Complaints
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.Wrap(err, "Invalid Payload"))
		return
	}

	err = uc.usecase.CreateComplaint(r.Context(), &c)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrDbFailure.Wrap(err, "failed to create complaint"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (uc *ComplaintHandler) GetComplaintByRole(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	user_id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	complaint, err := uc.usecase.GetComplaintByRole(r.Context(), user_id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(complaint)
}

func (uc *ComplaintHandler) UserMarkResolved(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	complaint_id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	if err := uc.usecase.UserMarkResolved(r.Context(), complaint_id); err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      complaint_id,
		"Message": "Complaint updated Successfully",
	})
}

func (uc *ComplaintHandler) GetAllComplaintByRole(w http.ResponseWriter, r *http.Request) {
	complaints, err := uc.usecase.GetAllComplaintByRole(r.Context())
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(complaints)
}

func (uc *ComplaintHandler) AdminUpdateComplaints(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	complaintID, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.Wrap(err, "Invalid status input"))
		return
	}

	err = uc.usecase.AdminUpdateComplaints(r.Context(), complaintID, body.Status)
	if err != nil {
		middleware.WriteError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"data": "Complaint updated successfully",
	})
}

// complaint_messages table
func (uc *ComplaintHandler) InsertCoplaintMessage(w http.ResponseWriter, r *http.Request) {
	complaintIdStr := mux.Vars(r)["id"]
	complaintID, err := strconv.Atoi(complaintIdStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrDbFailure.New("invalid id"))
		return
	}

	var cm models.ComplaintMessages
	if err := json.NewDecoder(r.Body).Decode(&cm); err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.Wrap(err, "Invalid payload"))
		return
	}

	senderID := middleware.GetUserId(r.Context())
	message := &models.ComplaintMessages{
		ComplaintID: complaintID,
		SenderID:    senderID,
		ParentID:    cm.ParentID,
		Message:     cm.Message,
		FileUrl:     cm.FileUrl,
	}

	err = uc.usecase.InsertCoplaintMessage(r.Context(), message)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (uc *ComplaintHandler) GetMessagesByComplaint(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	complaintID, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	message, err := uc.usecase.GetMessagesByComplaint(r.Context(), complaintID)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (uc *ComplaintHandler) ReplyToMessage(w http.ResponseWriter, r *http.Request) {
	// parse complaintID from url
	complaintIdStr := mux.Vars(r)["id"]
	complaintID, err := strconv.Atoi(complaintIdStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	// parse the form (text and file)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.Wrap(err, "can't parse multipart form"))
		return
	}

	// extract text message
	message := r.FormValue("message")
	if message == "" {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.Wrap(err, "message field is required"))
		return
	}

	// handle file
	var fileUrl string
	file, handler, err := r.FormFile("file")
	if err == nil && handler != nil {
		defer file.Close()
		fileUrl, err = utility.SaveUploadFile(file, *handler)
		if err != nil {
			middleware.WriteError(w, err)
			return
		}
	}

	// prepare message model
	msg := models.ComplaintMessages{
		Message:     message,
		ComplaintID: complaintID,
		FileUrl:     fileUrl,
	}

	if err := uc.usecase.ReplyToMessage(r.Context(), &msg); err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Message": "Reply Message added successfully",
		"data":    msg,
	})
}
