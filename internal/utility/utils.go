package utility

import (
	"crud_api/config"
	"time"

	appErrors "crud_api/internal/errors"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", appErrors.ErrDbFailure.New("usecase: Failed to generate password")
	}
	return string(hash), nil
}

func ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

var cfg = config.LoadConfig()

func GenerateJWT(userID int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", appErrors.ErrUnauthorized.Wrap(err, "Failed to sign JWT token")
	}

	return signed, nil

}

func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	//Decode/validateit
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appErrors.ErrUnauthorized.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})
	// Check for parsing error or invalid token
	if err != nil || !token.Valid {
		return nil, appErrors.ErrUnauthorized.Wrap(err, "Invalid or expired token")
	}

	// Safe type assertion for claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, appErrors.ErrUnauthorized.New("token claim not ok")
	}

	// Explicit expiration check
	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return nil, appErrors.ErrUnauthorized.New("token expired")
	}

	return claims, nil

}
