package main

import (
	"Complaingo/config"
	"Complaingo/internal/kafka"
	"Complaingo/internal/rabbitmq"
	"Complaingo/internal/redis"
	"Complaingo/internal/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to PostgreSQL DB
	db := config.ConnectToDB()
	defer db.Close(context.Background())

	// Connect to Redis
	redis.ConnectRedis()

	// Setup RabbitMQ
	rabbit := rabbitmq.NewProducer("amqp://guest:guest@localhost:5672/", "notifications")
	consumer := rabbitmq.NewConsumer("amqp://guest:guest@localhost:5672/", "notifications")
	rabbit.SendMessage("Hello from go rabbit")
	consumer.StartConsuming()

	// Setup Kafka
	kafkaConsumer := kafka.NewKafkaConsumer([]string{"localhost:9092"}, "chat-messages", "chat-group")
	kafkaCtx, kafkaStop := context.WithCancel(context.Background())
	kafkaConsumer.StartConsuming(kafkaCtx)
	kafkaProducer := kafka.NewKafkaProducer([]string{"localhost:9092"}, "chat-messages")

	// initialize router
	r := router.NewRouter(cfg, db, kafkaProducer)

	// start HTTP server
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

	// Graceful shutdown

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdownCtx.Done()
	fmt.Println("Shutting down server...")

	kafkaStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Graceful shutdown failed: %v\n", err)
	}
	log.Println("Server shutdown complete")
}
