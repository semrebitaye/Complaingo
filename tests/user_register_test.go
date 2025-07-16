package tests

import (
	"Complaingo/testutils"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegistration(t *testing.T) {
	// Clean the database
	testutils.CleanTestDB()

	// define fake user to register
	payload := map[string]string{
		"first_name": "Test",
		"last_name":  "User",
		"email":      "register_test@gmail.com",
		"password":   "password123",
		"role":       "admin",
	}

	// Send a POST /register request
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", testServer.URL+"/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// send the HTTP request
	client := &http.Client{}
	resp, Err := client.Do(req)
	assert.NoError(t, Err)
	defer resp.Body.Close()

	// Assert HTTP Status Code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Validate user is inserted in the DB
	var userID int
	err := testutils.GetTestDB().QueryRow(context.Background(),
		"SELECT id FROM users WHERE email=$1", payload["email"]).Scan(&userID)
	assert.NoError(t, err)
	assert.NotZero(t, userID)
}
