package tests

import (
	"Complaingo/config"
	"Complaingo/internal/router"
	"Complaingo/testutils"
	"context"
	"net/http/httptest"
	"os"
	"testing"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	// get the shared DB connection
	db := testutils.GetTestDB()
	// clean old data before starting tests
	testutils.CleanTestDB()
	// close DB after all tests complete
	defer db.Close(context.Background())
	// cleanup leftover data after all tests
	defer testutils.CleanTestDB()

	// Setup test server
	// load .env.test environment
	cfg := config.LoadConfig()
	// build full http.Handler with routes and middleware
	r := router.NewRouter(cfg, db, nil)
	// start a test server
	testServer = httptest.NewServer(r)
	// shuts it down after tests
	defer testServer.Close()

	// run all tests in the current package and exit the process to determine pass or fail
	os.Exit(m.Run())
}
