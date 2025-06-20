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

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid user payload"))
		return
	}

	if err := h.usecase.RegisterUser(r.Context(), &u); err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetAllUser(r.Context())
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid id"))
		return
	}
	user, err := h.usecase.GetUserByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	u := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalide user data"))
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("id invalid"))
		return
	}

	u.ID = id
	err = h.usecase.UpdateUser(r.Context(), u)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("id invalid"))
	}

	err = h.usecase.DeleteUser(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User Deleted Successfully")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid login data"))
		return
	}

	token, err := h.usecase.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
