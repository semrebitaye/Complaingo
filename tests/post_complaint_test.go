package tests

import (
	"Complaingo/testutils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func createAdminUser(t *testing.T) (int, string) {
	db := testutils.GetTestDB()

	// Ensure admin role exists
	_, err := db.Exec(context.Background(), `
        INSERT INTO roles (id, name) VALUES (1, 'admin')
        ON CONFLICT DO NOTHING`)
	assert.NoError(t, err)

	// Create admin user
	var adminID int
	email := fmt.Sprintf("admin-%d@gmail.com", time.Now().UnixNano())
	err = db.QueryRow(context.Background(), `
        INSERT INTO users 
            (first_name, last_name, email, password, role_id) 
        VALUES 
            ($1, $2, $3, $4, $5)
        RETURNING id`,
		"Admin", "User", email, "$2a$10$hashedpassword", 1,
	).Scan(&adminID)
	assert.NoError(t, err)

	return adminID, getTestJWT(adminID, "admin")
}

func createTestUser(t *testing.T) (int, string) {
	db := testutils.GetTestDB()

	// Ensure user role exists
	_, err := db.Exec(context.Background(), `
        INSERT INTO roles (id, name) VALUES (2, 'user')
        ON CONFLICT DO NOTHING`)
	assert.NoError(t, err)

	// Create regular user
	var userID int
	email := fmt.Sprintf("user-%d@example.com", time.Now().UnixNano())
	err = db.QueryRow(context.Background(), `
        INSERT INTO users 
            (first_name, last_name, email, password, role_id) 
        VALUES 
            ($1, $2, $3, $4, $5)
        RETURNING id`,
		"Test", "User", email, "$2a$10$hashedpassword", 2, // role_id 2 = user
	).Scan(&userID)
	assert.NoError(t, err)

	return userID, getTestJWT(userID, "user")
}

func getTestJWT(userID int, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	signed, err := token.SignedString([]byte("test_secret"))
	if err != nil {
		panic(fmt.Sprintf("Failed to generate JWT: %v", err))
	}
	return signed
}

func TestCreateComplaint(t *testing.T) {
	// verify DB connection
	db := testutils.GetTestDB()
	var exists bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public')").Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, "Database tables not found")

	_, token := createTestUser(t)

	payload := map[string]interface{}{
		"user_id": 1,
		"subject": "Test Complaint",
		"message": "Test description",
		"status":  "Created",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", testServer.URL+"/complaints", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
