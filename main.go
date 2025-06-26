package main

import (
	"context"
	"crud_api/config"
	"crud_api/internal/handler"
	"crud_api/internal/middleware"
	"crud_api/internal/notifier"
	"crud_api/internal/repository"
	"crud_api/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	websocketpkg "crud_api/internal/websocket"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to DB
	db := config.ConnectToDB()
	defer db.Close(context.Background())

	repo := repository.NewPgxUserRepo(db)
	usercase := usecase.NewUserUsecase(repo)
	userHandler := handler.NewUserHandler(usercase)

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("/uploads"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	authR := r.PathPrefix("/").Subrouter()
	authR.Use(middleware.Authentication)

	authR.Handle("/users", middleware.RBAC("admin", "user")(http.HandlerFunc(userHandler.GetAllUser))).Methods("GET")
	authR.Handle("/user/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(userHandler.GetUserByID))).Methods("GET")
	authR.Handle("/users/{id}", middleware.RBAC("admin")(http.HandlerFunc(userHandler.UpdateUser))).Methods("PATCH")
	authR.Handle("/users/{id}", middleware.RBAC("admin")(http.HandlerFunc(userHandler.DeleteUser))).Methods("DELETE")

	complaintRepo := repository.NewPgxComplaintRepo(db)
	complaintMessageRepo := repository.NewPgxComplaintMessageRepo(db)
	notifier := &notifier.RealTimeNotifier{}
	complaintUsecase := usecase.NewComplaintUsecase(complaintRepo, complaintMessageRepo, notifier)

	complaintHandler := handler.NewComplaintHandler(complaintUsecase)

	authR.HandleFunc("/ws", websocketpkg.HandleWebsocket).Methods("GET")

	authR.Handle("/complaints", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.CreateComplaint))).Methods("POST")
	authR.Handle("/complaints/user/{id}", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.GetComplaintByRole))).Methods("GET")
	authR.Handle("/complaints/{id}/resolve", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.UserMarkResolved))).Methods("PATCH")
	authR.Handle("/complaints", middleware.RBAC("admin")(http.HandlerFunc(complaintHandler.GetAllComplaintByRole))).Methods("GET")
	authR.Handle("/complaints/{id}/status", middleware.RBAC("admin")(http.HandlerFunc(complaintHandler.AdminUpdateComplaints))).Methods("PATCH")

	authR.Handle("/complaints/{id}/messages", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.InsertCoplaintMessage))).Methods("POST")
	authR.Handle("/complaints/{id}/messages", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.GetMessagesByComplaint))).Methods("GET")
	authR.Handle("/complaints/{id}/reply", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.ReplyToMessage))).Methods("POST")

	// create server
	srv := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: middleware.RecoverMiddleware(r),
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
