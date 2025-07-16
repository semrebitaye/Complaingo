package tests

import (
	"Complaingo/testutils"
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidLoginCase(t *testing.T) {
	testutils.CleanTestDB()
	defer testutils.CleanTestDB()

	// create user with known credentials
	db := testutils.GetTestDB()
	_, err := db.Exec(context.Background(), `
	INSERT INTO roles(id, name) VALUES (1, 'admin'),(2, 'user') ON CONFLICT DO NOTHING;
	INSERT INTO users (first_name, last_name, email, password, role_id)
	VALUES ('test', 'user', 'test@gmail.com', '$2a$10$zWpP41dt6cSZwChZ6GYA4u11wIfwk8nlDQD.3uFPReU9G13H3ItDW', 1)
	`)
	assert.NoError(t, err)

	tests := []struct {
		name          string
		payload       string
		wantCode      int
		wantErrorType string
	}{
		{
			name:          "Wrong Password",
			payload:       `{"email": "test@gmail.com", "password": "unauthorized"}`,
			wantCode:      http.StatusBadRequest,
			wantErrorType: "invalid_payload",
		},
		{
			name:          "Non-Existent email",
			payload:       `{"email": "notfound@gmail.com", "password": "anything"}`,
			wantCode:      http.StatusBadRequest,
			wantErrorType: "invalid_payload",
		},
		{
			name:          "Missing email",
			payload:       `{"password": "anything"}`,
			wantCode:      http.StatusBadRequest,
			wantErrorType: "email",
		},
		{
			name:          "Missing password",
			payload:       `{"email": "test@gmail.com",}`,
			wantCode:      http.StatusBadRequest,
			wantErrorType: "Invalid login data",
		},
		{
			name:          "Empty payload",
			payload:       `{}`,
			wantCode:      http.StatusBadRequest,
			wantErrorType: "email",
		},
		{
			name:          "Invalid json",
			payload:       `Invalid-json`,
			wantCode:      http.StatusBadRequest,
			wantErrorType: "Invalid login data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare request with correct content type
			req, err := http.NewRequest("POST", testServer.URL+"/login", bytes.NewBufferString(tt.payload))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.wantCode, resp.StatusCode)
			assert.Contains(t, string(body), tt.wantErrorType)
		})
	}
}
