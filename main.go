package main

import (
	"context"
	"crud_api/config"
	"crud_api/internal/handler"
	"crud_api/internal/middleware"
	"crud_api/internal/repository"
	"crud_api/internal/usecase"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to DB
	db := config.ConnectToDB()
	defer db.Close(context.Background())

	repo := repository.NewPgxUserRepo(db)
	usercase := usecase.NewUserUsecase(repo)
	handler := handler.NewUserHandler(usercase)

	r := mux.NewRouter()
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")

	authR := r.PathPrefix("/").Subrouter()
	authR.Use(middleware.Authentiction)

	authR.HandleFunc("/users", handler.GetAllUser).Methods("GET")
	authR.HandleFunc("/user/{id}", handler.GetUserByID).Methods("GET")
	authR.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PATCH")
	authR.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	log.Println("Listening on port", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
