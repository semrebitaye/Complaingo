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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
