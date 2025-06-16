package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func ConnectToDB() *pgx.Conn {
	cfg := LoadConfig()
	conn, err := pgx.Connect(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)
	}
	log.Println("Connected to db successfully")

	return conn
}
