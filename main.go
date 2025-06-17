package main

import (
	"context"
	"crud_api/config"
	"crud_api/internal/handler"
	"crud_api/internal/middleware"
	"crud_api/internal/repository"
	"crud_api/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	authR.Use(middleware.Authentication)

	authR.HandleFunc("/users", handler.GetAllUser).Methods("GET")
	authR.HandleFunc("/user/{id}", handler.GetUserByID).Methods("GET")
	authR.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PATCH")
	authR.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	// log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))

	// create server
	srv := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Println("Listening on port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Stopped Listening: %v\n", err)
		}
	}()

	shoutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shoutdown.Done()
	fmt.Println("Souting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Graceful shotdown failed: %v\n", err)
	}

	log.Println("Server shoutdown complete")
}
