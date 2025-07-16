package tests

import (
	"Complaingo/testutils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	// clean old data
	testutils.CleanTestDB()
	defer testutils.CleanTestDB()

	// 2, insert test user into database
	db := testutils.GetTestDB()
	_, err := db.Exec(context.Background(), `
	INSERT INTO roles (id, name) 
	VALUES (1, 'admin'), (2, 'user') 
	ON CONFLICT DO NOTHING;
	INSERT INTO users (first_name, last_name, email, password, role_id) 
	VALUES 
	('dawit', 'abebe', 'dave@gmail.com', 'devaman', 2), 
	('goshu', 'yirga', 'goshu@gmail.com', 'goshumanega', 1)
	`)
	assert.NoError(t, err)

	token := getTestJWT(1, "admin")

	// 3, make get/users request
	req, err := http.NewRequest("GET", testServer.URL+"/users", nil)
	assert.NoError(t, err)
	// set header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// 4,send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 5,decode response
	var response testutils.GenericAPIResponse[[]map[string]interface{}]
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	fmt.Printf("Returned Users: %+v\n", response.Data)
	// 6, assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, len(response.Data))

	emails := []string{
		response.Data[0]["email"].(string),
		response.Data[1]["email"].(string),
	}

	assert.Contains(t, emails, "dave@gmail.com")
	assert.Contains(t, emails, "goshu@gmail.com")
}
