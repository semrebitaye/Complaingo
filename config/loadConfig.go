package config

import (
	"log"
	"os"

	appErrors "crud_api/internal/errors"

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
		panic(appErrors.ErrInvalidPayload.New("url not found"))
	}

	jwtSecret := os.Getenv("secret")
	if jwtSecret == "" {
		panic(appErrors.ErrInvalidPayload.New("Secret not found"))
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		panic(appErrors.ErrInvalidPayload.Wrap(err, "port must be a number"))
	}

	return &Config{
		DBUrl:      dbUrl,
		JWTSecret:  jwtSecret,
		ServerPort: serverPort,
	}

}
