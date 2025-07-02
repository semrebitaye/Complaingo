package main

import (
	"context"
	"crud_api/config"
	"crud_api/internal/handler"
	"crud_api/internal/kafka"
	"crud_api/internal/middleware"
	"crud_api/internal/notifier"
	"crud_api/internal/rabbitmq"
	"crud_api/internal/redis"
	"crud_api/internal/repository"
	"crud_api/internal/usecase"
	"crud_api/internal/websocket"
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

	redis.ConnectRedis()

	rabbit := rabbitmq.NewProducer("amqp://guest:guest@localhost:5672/", "notifications")
	consumer := rabbitmq.NewConsumer("amqp://guest:guest@localhost:5672/", "notifications")
	rabbit.SendMessage("Hello from go rabbit")
	consumer.StartConsuming()

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

	messageRepo := repository.NewMessageRepository(db)
	kafkaConsumer := kafka.NewKafkaConsumer([]string{"localhost:9092"}, "chat-messages", "chat-group")
	kafkaCtx, kafkaStop := context.WithCancel(context.Background())
	kafkaConsumer.StartConsuming(kafkaCtx)
	kafkaProducer := kafka.NewKafkaProducer([]string{"localhost:9092"}, "chat-messages")
	websocketHandler := websocket.NewwebsocketHandler(messageRepo, kafkaProducer)

	authR.HandleFunc("/ws", websocketHandler.HandleWebsocket).Methods("GET")

	authR.Handle("/complaints", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.CreateComplaint))).Methods("POST")
	authR.Handle("/complaints/user/{id}", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.GetComplaintByRole))).Methods("GET")
	authR.Handle("/complaints/{id}/resolve", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.UserMarkResolved))).Methods("PATCH")
	authR.Handle("/complaints", middleware.RBAC("admin")(http.HandlerFunc(complaintHandler.GetAllComplaintByRole))).Methods("GET")
	authR.Handle("/complaints/{id}/status", middleware.RBAC("admin")(http.HandlerFunc(complaintHandler.AdminUpdateComplaints))).Methods("PATCH")

	authR.Handle("/complaints/{id}/messages", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.InsertCoplaintMessage))).Methods("POST")
	authR.Handle("/complaints/{id}/messages", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.GetMessagesByComplaint))).Methods("GET")
	authR.Handle("/complaints/{id}/reply", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.ReplyToMessage))).Methods("POST")

	docRepo := repository.NewDocumentRepository(db)
	docUC := usecase.NewDocumentUsecase(docRepo)
	docHandle := handler.NewDocumentHandler(docUC)

	authR.Handle("/documents", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandle.Uplod))).Methods("POST")
	authR.Handle("/documents/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandle.GetDocumentByID))).Methods("GET")
	authR.Handle("/documents/user/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandle.GetDocumentByUser))).Methods("GET")
	authR.Handle("/documents/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandle.DeleteDocument))).Methods("DELETE")

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

	kafkaStop() //gracefully stop kafka consumer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Graceful shotdown failed: %v\n", err)
	}

	log.Println("Server shoutdown complete")
}
