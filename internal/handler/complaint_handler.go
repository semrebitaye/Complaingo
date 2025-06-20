package handler

import (
	"crud_api/internal/domain/models"
	"crud_api/internal/middleware"
	"crud_api/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	appErrors "crud_api/internal/errors"

	"github.com/gorilla/mux"
)

type ComplaintHandler struct {
	usecase usecase.ComplaintUsecase
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
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (uc *ComplaintHandler) GetComplaintByRole(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	complaint, err := uc.usecase.GetComplaintByRole(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(complaint)
}

func (uc *ComplaintHandler) UserMarkResolved(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	if err := uc.usecase.UserMarkResolved(r.Context(), id); err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
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
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (uc *ComplaintHandler) ReplayToMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.ComplaintMessages
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrDbFailure.Wrap(err, "Invalid payload"))
		return
	}
	if err := uc.usecase.ReplayToMessage(r.Context(), &msg); err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Message": "Replay Message added successfully",
		"data":    msg,
	})
}
