package tests

import (
	redis "Complaingo/internal/redis"
	"Complaingo/testutils"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	// 1, clean database
	testutils.CleanTestDB()
	defer testutils.CleanTestDB()

	redis.ConnectRedis()
	db := testutils.GetTestDB()

	// 2, insert required role and test user
	_, err := db.Exec(context.Background(), `
	INSERT INTO roles (id, name) VALUES (1, 'admin') ON CONFLICT DO NOTHING;
    INSERT INTO users (first_name, last_name, email, password, role_id)
    VALUES ('alemu', 'tadese', 'alex@gmail.com', '$2a$10$aY5Qa2sV0oz6GFeOTn0pFOwmHfxEqmEGGdztj2NSA4xfZgBo2GQvW', 1)`)
	assert.NoError(t, err)

	// 3, prepare login payload
	payload := map[string]string{
		"email":    "alex@gmail.com",
		"password": "alexman",
	}
	body, _ := json.Marshal(payload)

	// 4, send login request
	req, err := http.NewRequest("POST", testServer.URL+"/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// 5, decode response
	var response testutils.GenericAPIResponse[string]
	err = json.Unmarshal(respBody, &response)
	assert.NoError(t, err)

	// assert values
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotEmpty(t, response.Data, "Token should not be empty")
}
