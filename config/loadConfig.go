package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl      string
	JWTSecret  string
	ServerPort string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalf("url not found: %v\n", err)
	}

	jwtSecret := os.Getenv("secret")
	if jwtSecret == "" {
		log.Fatalf("secret not found: %v\n", err)
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Fatalf("port not found: %v\n", err)
	}

	return &Config{
		DBUrl:      dbUrl,
		JWTSecret:  jwtSecret,
		ServerPort: serverPort,
	}

}
