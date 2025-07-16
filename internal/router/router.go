package router

import (
	"Complaingo/config"
	"Complaingo/internal/handler"
	"Complaingo/internal/kafka"
	"Complaingo/internal/middleware"
	"Complaingo/internal/notifier"
	"Complaingo/internal/repository"
	"Complaingo/internal/usecase"
	websocket "Complaingo/internal/websockets"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func NewRouter(cfg *config.Config, db *pgx.Conn, kafkaProducer *kafka.KafkaProducer) *mux.Router {
	r := mux.NewRouter()

	// serve static files
	fs := http.FileServer(http.Dir("/uploads"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	// ==== auth and users ====
	repo := repository.NewPgxUserRepo(db)
	usercase := usecase.NewUserUsecase(repo)
	userHandler := handler.NewUserHandler(usercase)

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	authR := r.PathPrefix("/").Subrouter()
	authR.Use(middleware.Authentication)

	authR.Handle("/ask-ai", middleware.RBAC("admin", "user")(http.HandlerFunc(handler.AIChatHandler))).Methods("POST")
	authR.Handle("/users", middleware.RBAC("admin", "user")(http.HandlerFunc(userHandler.GetAllUser))).Methods("GET")
	authR.Handle("/user/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(userHandler.GetUserByID))).Methods("GET")
	authR.Handle("/users/{id}", middleware.RBAC("admin")(http.HandlerFunc(userHandler.UpdateUser))).Methods("PATCH")
	authR.Handle("/users/{id}", middleware.RBAC("admin")(http.HandlerFunc(userHandler.DeleteUser))).Methods("DELETE")

	//  === complaint and complain message ===
	complaintRepo := repository.NewPgxComplaintRepo(db)
	complaintMessageRepo := repository.NewPgxComplaintMessageRepo(db)
	notif := &notifier.RealTimeNotifier{}
	complaintUC := usecase.NewComplaintUsecase(complaintRepo, complaintMessageRepo, notif)
	complaintHandler := handler.NewComplaintHandler(complaintUC)

	authR.Handle("/complaints", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.CreateComplaint))).Methods("POST")
	authR.Handle("/complaints/user/{id}", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.GetComplaintByRole))).Methods("GET")
	authR.Handle("/complaints/{id}/resolve", middleware.RBAC("user")(http.HandlerFunc(complaintHandler.UserMarkResolved))).Methods("PATCH")
	authR.Handle("/complaints", middleware.RBAC("admin")(http.HandlerFunc(complaintHandler.GetAllComplaintByRole))).Methods("GET")
	authR.Handle("/complaints/{id}/status", middleware.RBAC("admin")(http.HandlerFunc(complaintHandler.AdminUpdateComplaints))).Methods("PATCH")

	authR.Handle("/complaints/{id}/messages", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.InsertCoplaintMessage))).Methods("POST")
	authR.Handle("/complaints/{id}/messages", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.GetMessagesByComplaint))).Methods("GET")
	authR.Handle("/complaints/{id}/reply", middleware.RBAC("admin", "user")(http.HandlerFunc(complaintHandler.ReplyToMessage))).Methods("POST")

	//  === document ===
	docRepo := repository.NewDocumentRepository(db)
	docUC := usecase.NewDocumentUsecase(docRepo)
	docHandler := handler.NewDocumentHandler(docUC)

	authR.Handle("/documents", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandler.Uplod))).Methods("POST")
	authR.Handle("/documents/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandler.GetDocumentByID))).Methods("GET")
	authR.Handle("/documents/user/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandler.GetDocumentByUser))).Methods("GET")
	authR.Handle("/documents/{id}", middleware.RBAC("admin", "user")(http.HandlerFunc(docHandler.DeleteDocument))).Methods("DELETE")

	// === websocket ===
	msgRepo := repository.NewMessageRepository(db)
	wsHandler := websocket.NewwebsocketHandler(msgRepo, kafkaProducer)
	authR.HandleFunc("/ws", wsHandler.HandleWebsocket).Methods("GET")

	return r
}
