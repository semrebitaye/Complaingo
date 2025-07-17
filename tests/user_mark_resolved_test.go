package tests

import (
	"Complaingo/testutils"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestUserResolved(t *testing.T) (int, string) {
	userID := 1
	_, err := testutils.GetTestDB().Exec(context.Background(),
		`INSERT INTO users (id, email, password, role, first_name, last_name) 
         VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, "user@example.com", "hashedpass", "user", "Test", "User")
	assert.NoError(t, err)

	token := getTestJWT(userID, "user")
	return userID, token
}

func TestUserMarkResolved(t *testing.T) {
	testutils.CleanTestDB()
	testutils.InitTestSchema()

	// 1. Create admin and regular user
	userID, userToken := createTestUserResolved(t)

	// 2,insert complaint into DB
	complaintID := testutils.InsertComplaint(userID, "Feature request", "Add export to pdf", "Created")

	// 3, send patch request to mark complaint as resolved
	url := fmt.Sprintf("%s/complaints/%d/resolve", testServer.URL, complaintID)
	req, err := http.NewRequest("PATCH", url, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+userToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 4, check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 5. Verify updated status in DB
	var status string
	err = testutils.GetTestDB().QueryRow(context.Background(),
		"SELECT status FROM complaints WHERE id = $1", complaintID).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "Resolved", status)
}
