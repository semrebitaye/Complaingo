package handler

import (
	"context"
	"crud_api/internal/domain/models"
	"crud_api/internal/usecase"
	"encoding/json"
	"net/http"
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
		http.Error(w, "Failed to decode the user", http.StatusBadRequest)
		return
	}

	if err := h.usecase.RegisterUser(context.Background(), &u); err != nil {
		http.Error(w, "Failed to register the user"+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetAllUser(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch user"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
