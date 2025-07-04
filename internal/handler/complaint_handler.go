package handler

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/middleware"
	"Complaingo/internal/redis"
	"Complaingo/internal/usecase"
	"Complaingo/internal/utility"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	appErrors "Complaingo/internal/errors"

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

	middleware.WriteSuccess(w, c, "Complaint Creted Successfully")
}

func (uc *ComplaintHandler) GetComplaintByRole(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	user_id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	// check redis cache
	cachKey := fmt.Sprintf("complaints:%d", user_id)
	cachedComplaint, err := redis.RDB.Get(redis.Ctx, cachKey).Result()
	if err == nil {
		// found on cache
		var complaints []models.Complaints
		if err := json.Unmarshal([]byte(cachedComplaint), &complaints); err == nil {
			middleware.WriteSuccess(w, complaints, "Complient from cache")
			return
		}

	}

	// not found in cache fetch from DB
	complaint, err := uc.usecase.GetComplaintByRole(r.Context(), user_id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	// cache the database result for future use
	complientJson, _ := json.Marshal(complaint)
	redis.RDB.Set(redis.Ctx, cachKey, complientJson, time.Minute*10)

	middleware.WriteSuccess(w, complaint, "complaint get successfully by pk user_id")
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

	middleware.WriteSuccess(w, complaint_id, "Complait updated successfully")
}

func (uc *ComplaintHandler) GetAllComplaintByRole(w http.ResponseWriter, r *http.Request) {
	complaints, err := uc.usecase.GetAllComplaintByRole(r.Context())
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteSuccess(w, complaints, "All compliants successfully")
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

	middleware.WriteSuccess(w, body.Status, "Complaint Updated Successfully")
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

	middleware.WriteSuccess(w, message, "Message created Successfully")
}

func (uc *ComplaintHandler) GetMessagesByComplaint(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	complaintID, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}

	// check redis cache
	cacheKey := fmt.Sprintf("Message:%d", complaintID)
	cachedMessage, err := redis.RDB.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		var message []models.ComplaintMessages
		if err := json.Unmarshal([]byte(cachedMessage), &message); err == nil {
			middleware.WriteSuccess(w, message, "Feched from cache")
			return
		}
	}

	// not found in cache -> fetch from DB
	message, err := uc.usecase.GetMessagesByComplaint(r.Context(), complaintID)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	// save in cache for future
	messageJson, _ := json.Marshal(message)
	redis.RDB.Set(redis.Ctx, cacheKey, messageJson, time.Minute*10)

	middleware.WriteSuccess(w, message, "Message successfully fetched by complaint id")
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

	middleware.WriteSuccess(w, msg, "Reply Message added successfully")
}
