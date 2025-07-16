package tests

import (
	"Complaingo/internal/domain/models"
	"Complaingo/testutils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetComplaintsAsAdmin(t *testing.T) {
	testutils.CleanTestDB()
	testutils.InitTestSchema()

	// 1. Create admin and regular user
	_, adminToken := createAdminUser(t)
	userID, _ := createTestUser(t)

	// 2. Insert a complaint
	_, err := testutils.GetTestDB().Exec(context.Background(), `
        INSERT INTO complaints 
            (user_id, subject, message, status) 
        VALUES 
            ($1, $2, $3, $4)`,
		userID, "Test Complaint", "Test description", "Created")
	assert.NoError(t, err)

	// 3. Make request as admin
	req, err := http.NewRequest("GET", testServer.URL+"/complaints", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 4. Check status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Admin should get a successful response")

	// 5. Parse body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respPayload testutils.GenericAPIResponse[[]models.Complaints]

	err = json.Unmarshal(body, &respPayload)
	if err != nil {
		t.Fatalf("Failed to parse response: %v\nResponse body: %s", err, string(body))
	}

	assert.GreaterOrEqual(t, len(respPayload.Data), 1, "Should return at least one complaint")

}
