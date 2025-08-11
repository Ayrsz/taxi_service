package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// "taxi_service/config"
	// "taxi_service/database"
	"taxi_service/routes"

	"github.com/gofiber/fiber/v2"
)

// SetupTestApp creates a new Fiber app for testing
func SetupTestApp(t *testing.T) *fiber.App {
	// Não precisamos mais de configuração de banco para JSON
	app := fiber.New()
	routes.SetupRoutes(app)

	return app
}

// CleanupTestApp cleans up after tests
func CleanupTestApp(t *testing.T) {
	// Não precisa mais de cleanup de banco
}

// MakeRequest makes an HTTP request to the test app
func MakeRequest(t *testing.T, app *fiber.App, method, path string, body interface{}) *http.Response {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	return resp
}

// ParseResponseBody parses the response body into the given interface
func ParseResponseBody(t *testing.T, resp *http.Response, v interface{}) {
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
}
