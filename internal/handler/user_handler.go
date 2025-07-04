package handler

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/middleware"
	"Complaingo/internal/redis"
	"Complaingo/internal/usecase"
	"Complaingo/internal/validation"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	appErrors "Complaingo/internal/errors"

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

	// check redis catch first
	chachKey := fmt.Sprintf("user:%d", id)
	catchUserJson, err := redis.RDB.Get(redis.Ctx, chachKey).Result()

	if err == nil {
		// found on catch return it
		var cathUser models.User
		if err := json.Unmarshal([]byte(catchUserJson), &cathUser); err == nil {
			middleware.WriteSuccess(w, cathUser, "user fetched from catch")
			return
		}
	}

	// if not found in cache, fetch from DB
	user, err := h.usecase.GetUserByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	// cache the result for future use
	userJson, _ := json.Marshal(user)
	redis.RDB.Set(redis.Ctx, chachKey, userJson, time.Minute*10)

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

	resp, err := h.usecase.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	// cache the user login in redis
	userJson, _ := json.Marshal(resp.User)
	redis.RDB.Set(redis.Ctx, fmt.Sprintf("User:%d", resp.User.ID), userJson, time.Minute*10)

	middleware.WriteSuccess(w, resp.Token, "User login successfully")
}
