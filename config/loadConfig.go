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
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("url not found")
	}

	jwtSecret := os.Getenv("secret")
	if jwtSecret == "" {
		log.Fatal("secret not found")
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Fatal("port not found")
	}

	return &Config{
		DBUrl:      dbUrl,
		JWTSecret:  jwtSecret,
		ServerPort: serverPort,
	}

}
