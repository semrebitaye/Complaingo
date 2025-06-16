package handler

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/internal/middleware"
	"crud_api/internal/usecase"
	"encoding/json"
	"net/http"
	"strconv"

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
		middleware.WriteError(w, http.StatusBadRequest, "user data invalid")
		return
	}

	if err := h.usecase.RegisterUser(context.Background(), &u); err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "failed to register the user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetAllUser(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "failed to fetch the user")
		return
	}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "id not found")
		return
	}
	user, err := h.usecase.GetUserByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusNotFound, "Failed to get the user by the req id")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	u := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "user data invalid")
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "id not found")
		return
	}

	u.ID = id
	err = h.usecase.UpdateUser(r.Context(), u)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "Failed to update the user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "Failed to extract the id from the req url")
	}

	err = h.usecase.DeleteUser(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, http.StatusNotFound, "Failed to delete the user by the req id")
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
		middleware.WriteError(w, http.StatusBadRequest, "login request invalid")
		return
	}

	token, err := h.usecase.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "user credential invalid")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
