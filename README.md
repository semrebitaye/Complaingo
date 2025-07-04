## Complaingo
Complaingo is an Online Complaint Management System built with Go (Golang) that enables users and admins to efficiently manage complaints, communicate in real time, and streamline customer service workflows. This system is designed with a modern architecture incorporating WebSockets, Kafka, Redis, RabbitMQ, and OpenAI for intelligent feedback support.

### Features
#### Role-Based Access Control (RBAC): 
   Admin and User roles with secure access
#### Real-Time Communication: 
    WebSocket chat between users and admins
#### Pub/Sub Channels: 
    Dynamic message broadcasting using custom channels
#### Complaint Submission and Resolution: 
    User complaint creation, admin status updates
#### Document Upload: 
    Upload and retrieve documents tied to users
#### Authentication: 
    JWT-based login & registration
#### Redis Caching: 
    Speed up repeated queries (e.g., user or complaint data)
#### RabbitMQ: 
    Message queue for system notifications
#### OpenAI Integration: 
    Generate smart responses or summaries (API key required)
#### Kafka Integration: 
    Message streaming and decoupled architecture support

### Tech Stack

#### Layer                                                   Technology

    Backend                                                 Go (Golang)
    Framework                                               Gorilla Mux (Routing)
    DB                                                      PostgreSQL
    Real-Time                                               Gorilla WebSocket   
    Caching                                                 Redis
    Messaging                                               Kafka, RabbitMQ
    AI Integration                                          OpenAI API
    Auth                                                    JWT

### Test Endpoints via Postman
POST /register – Register user

POST /login – Login and get JWT

GET /ws – Connect to WebSocket for real-time chat

POST /complaints – Submit a complaint

POST /documents – Upload a document

### Project Structure
complaingo/
├── api_gateway/        # API Gateway service (reverse proxy, rate limit)
├── config/             # Configuration files
├── db/                 # SQL migration files or database setup
├── internal/           # Internal application logic
│   ├── domain/         # Domain models
│   ├── errors/         # Custom error types
│   ├── validation/     # Request validation logic
│   ├── utility/        # Shared helper functions
│   ├── handler/        # HTTP route handlers
│   ├── middleware/     # Auth, logging, recovery, RBAC
│   ├── repository/     # Database access logic
│   ├── usecase/        # Business logic and workflows
│   ├── websocket/      # Real-time chat handlers
│   ├── kafka/          # Kafka integration (producer/consumer)
│   ├── rabbitmq/       # RabbitMQ integration
│   └── notifier/       # Real-time notifier interfaces
├── uploads_doc/        # Directory for uploaded documents
├── .env                # Environment variable definitions
├── .gitignore          # Git ignore file
├── docker-compose.yaml            # Compose file for main services
├── docker-compose.kafka.yml       # Compose file for Kafka
├── Dockerfile          # Docker build instructions
├── go.mod              # Go module definition
├── go.sum              # Go module checksum file
└── main.go             # Main application entry point