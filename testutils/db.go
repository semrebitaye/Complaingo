package testutils

import (
	"Complaingo/config"
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

var (
	testDB     *pgx.Conn
	dbInitOnce sync.Once
)

// GenericAPIResponse
type GenericAPIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func InitTestSchema() {
	// get database connection
	db := GetTestDB()

	// Create roles table first
	_, err := db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS roles (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) NOT NULL UNIQUE
        )`)
	if err != nil {
		log.Fatalf("Failed to create roles table: %v", err)
	}

	// Insert required roles
	_, err = db.Exec(context.Background(), `
        INSERT INTO roles (id, name) 
        VALUES (1, 'admin'), (2, 'user')
        ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Fatalf("Failed to seed roles: %v", err)
	}

	// Create users table with FK constraint
	_, err = db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(100) NOT NULL,
            last_name VARCHAR(100) NOT NULL,
            email VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            role_id INTEGER NOT NULL REFERENCES roles(id)
        )`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create complaints table with FK constraint
	_, err = db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS complaints (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL REFERENCES users(id),
            subject VARCHAR(255) NOT NULL,
            description TEXT NOT NULL,
            status VARCHAR(50) DEFAULT 'pending',
            created_at TIMESTAMP DEFAULT NOW()
        )`)
	if err != nil {
		log.Fatalf("Failed to create complaints table: %v", err)
	}
}

func GetTestDB() *pgx.Conn {
	dbInitOnce.Do(func() { //ensures this block runs once only
		os.Setenv("ENV", "test")
		cfg := config.LoadConfig()

		// Verify connection string
		if cfg.DBUrl == "" {
			log.Fatal("DATABASE_URL is empty - check your .env.test file")
		}

		var err error
		var attempts int
		maxAttempts := 10
		delay := time.Second * 2

		for attempts < maxAttempts {
			testDB, err = pgx.Connect(context.Background(), cfg.DBUrl)
			if err == nil {
				// Verify we can query
				var dbName string
				err = testDB.QueryRow(context.Background(), "SELECT current_database()").Scan(&dbName)
				if err == nil {
					log.Printf("Connected to database: %s", dbName)
					break
				}
			}

			log.Printf("Attempt %d failed: %v", attempts+1, err)
			time.Sleep(delay)
			attempts++
		}

		if err != nil {
			log.Fatalf("Failed to connect after %d attempts. Verify:\n"+
				"1. Docker container is running\n"+
				"2. Database name matches exactly\n"+
				"3. Credentials are correct\n"+
				"4. Port mapping is correct\n"+
				"Error: %v", maxAttempts, err)
		}
	})
	return testDB
}

func CleanTestDB() {
	db := GetTestDB()
	_, err := db.Exec(context.Background(), `
		TRUNCATE 
			users, 
			complaints, 
			complaint_messages
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		log.Printf("failed to clean test database: %v", err)
	}
}
