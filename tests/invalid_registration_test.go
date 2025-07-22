package tests

import (
	"Complaingo/testutils"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegistrationError(t *testing.T) {
	testutils.CleanTestDB()
	tests := []struct {
		name       string
		payload    map[string]string
		wantStatus int
	}{
		{
			name: "missing email",
			payload: map[string]string{
				"first_name": "Test",
				"last_name":  "User",
				"password":   "password123",
				"role":       "admin",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email format",
			payload: map[string]string{
				"first_name": "Test",
				"last_name":  "User",
				"email":      "invalid-email",
				"password":   "password123",
				"role":       "admin",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "weak password",
			payload: map[string]string{
				"first_name": "Test",
				"last_name":  "User",
				"email":      "weakPass@gmail.com",
				"password":   "123",
				"role":       "admin",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "unsupported role",
			payload: map[string]string{
				"first_name": "Test",
				"last_name":  "User",
				"email":      "rolefail@gmail.com",
				"password":   "password123",
				"role":       "superadmin",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			payload: map[string]string{
				"first_name": "Test",
				"last_name":  "User",
				"email":      "register_test@gmail.com", //same as success test
				"password":   "password123",
				"role":       "admin",
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req, err := http.NewRequest("POST", testServer.URL+"/register", bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)
		})
	}

}
