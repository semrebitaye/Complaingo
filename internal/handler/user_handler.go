package handler

import (
	"crud_api/internal/domain/models"
	"crud_api/internal/middleware"
	"crud_api/internal/usecase"
	"crud_api/internal/validation"
	"encoding/json"
	"log"
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

	middleware.WriteSuccess(w, u, "User Registered Successfully")
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	log.Println("🧠 User role from context:", middleware.GetUserRole(r.Context()))

	users, err := h.usecase.GetAllUser(r.Context())
	if err != nil {
		log.Panicln("users not found ", users)
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteSuccess(w, users, "All users retrieved Successfully")
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

	middleware.WriteSuccess(w, user, "user successfully get by pk id")
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

	middleware.WriteSuccess(w, u, "User updated successfully")
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

	middleware.WriteSuccess(w, nil, "User Deleted Successfully")
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body validation.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.WriteError(w, appErrors.ErrInvalidPayload.New("Invalid login data"))
		return
	}

	token, err := h.usecase.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteSuccess(w, token, "User login successfully")
}
