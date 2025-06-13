package main

import (
	"context"
	"crud_api/config"
	"crud_api/internal/handler"
	"crud_api/internal/repository"
	"crud_api/internal/usecase"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnvVariable()

	db := config.ConnectDB()
	defer db.Close(context.Background())

	repo := repository.NewPgxUserRepo(db)
	usercase := usecase.NewUserUsecase(repo)
	handler := handler.NewUserHandler(usercase)

	r := mux.NewRouter()
	r.HandleFunc("/register", handler.Register).Methods("POST")
	r.HandleFunc("/users", handler.GetAllUser).Methods("GET")
	r.HandleFunc("/user/{id}", handler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PATCH")
	r.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
