package config

import (
	"log"
	"os"

	appErrors "Complaingo/internal/errors"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl      string
	JWTSecret  string
	ServerPort string
}

func LoadConfig() *Config {
	envFile := ".env"
	if os.Getenv("ENV") == "test" {
		envFile = "../.env.test"
	}

	// Check if custom config path was provided
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		envFile = configPath
	}

	// Verify file exists
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Printf("Warning: Config file %s not found, using environment variables", envFile)
	} else {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Printf("Warning: Error loading %s: %v", envFile, err)
		}
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
		panic(appErrors.ErrInvalidPayload.New("port must be a number"))
	}

	return &Config{
		DBUrl:      dbUrl,
		JWTSecret:  jwtSecret,
		ServerPort: serverPort,
	}
}
