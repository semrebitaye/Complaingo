package utility

import (
	"crud_api/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"emil":    email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	cfg := config.LoadConfig()

	tocken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tocken.SignedString([]byte(cfg.JWTSecret))

}
