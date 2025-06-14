package utility

import (
	"crud_api/config"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

var cfg = config.LoadConfig()

func GenerateJWT(userID int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"emil":    email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))

}

func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	//Decode/validateit
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
	// Check for parsing error or invalid token
	if err != nil || !token.Valid {
		log.Printf("Token parse error: %v\n", err)
		return nil, err
	}

	// Safe type assertion for claims
	claims, ok := token.Claims.(jwt.MapClaims)
	// check exp date
	if !ok {
		log.Println("token expired")
		return nil, err
	}
	return claims, nil

}
