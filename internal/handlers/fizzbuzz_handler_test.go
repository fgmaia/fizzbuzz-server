package handlers_test

import (
	"encoding/json"
	"fizzbuzz-server/internal/apps"
	"fizzbuzz-server/internal/handlers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	// Initialize the application
	apps.App()

	// Create a new Fiber app
	app := fiber.New()

	// Add validator middleware
	validate := validator.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("validator", validate)
		return c.Next()
	})

	// Register the FizzBuzz handler
	app.Get("/fizzbuzz", handlers.FizzbuzzHandler)

	return app
}

func TestFizzbuzzHandler_ValidRequest(t *testing.T) {
	// Setup
	app := setupApp()

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string][]string
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the response
	result, exists := response["result"]
	assert.True(t, exists)
	assert.Equal(t, 15, len(result))

	// Check specific values in the result
	assert.Equal(t, "1", result[0])
	assert.Equal(t, "2", result[1])
	assert.Equal(t, "fizz", result[2])
	assert.Equal(t, "4", result[3])
	assert.Equal(t, "buzz", result[4])
	assert.Equal(t, "fizz", result[5])
	assert.Equal(t, "7", result[6])
	assert.Equal(t, "8", result[7])
	assert.Equal(t, "fizz", result[8])
	assert.Equal(t, "buzz", result[9])
	assert.Equal(t, "11", result[10])
	assert.Equal(t, "fizz", result[11])
	assert.Equal(t, "13", result[12])
	assert.Equal(t, "14", result[13])
	assert.Equal(t, "fizzbuzz", result[14])
}

func TestFizzbuzzHandler_MissingParameters(t *testing.T) {
	// Setup
	app := setupApp()

	// Create a test request with missing parameters
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]string
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the error message
	errorMsg, exists := response["error"]
	assert.True(t, exists)
	assert.Contains(t, errorMsg, "validation")
}

func TestFizzbuzzHandler_InvalidParameters(t *testing.T) {
	// Setup
	app := setupApp()

	// Create a test request with invalid parameters
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=0&int2=5&limit=15&str1=fizz&str2=buzz", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]string
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the error message
	errorMsg, exists := response["error"]
	assert.True(t, exists)
	assert.Contains(t, errorMsg, "validation")
}

func TestFizzbuzzHandler_DefaultValues(t *testing.T) {
	// Setup
	app := setupApp()

	// Create a test request without str1 and str2
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string][]string
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the response
	result, exists := response["result"]
	assert.True(t, exists)
	assert.Equal(t, 15, len(result))

	// Check that default values were used
	assert.Equal(t, "fizzbuzz", result[14]) // Should be "fizzbuzz" at position 15
}

func TestFizzbuzzHandler_LargeLimit(t *testing.T) {
	// Setup
	app := setupApp()

	// Create a test request with a large limit
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=100&str1=fizz&str2=buzz", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string][]string
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the response
	result, exists := response["result"]
	assert.True(t, exists)
	assert.Equal(t, 100, len(result))
}

func TestFizzbuzzHandler_ExceedMaxLimit(t *testing.T) {
	// Setup
	app := setupApp()

	// Create a test request with a limit exceeding the maximum
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=20000&str1=fizz&str2=buzz", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string]string
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the error message
	errorMsg, exists := response["error"]
	assert.True(t, exists)
	assert.Contains(t, errorMsg, "validation")
}

func TestFizzbuzzHandler_CustomStrings(t *testing.T) {
	// Setup
	app := setupApp()

	// Create a test request with custom strings
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15&str1=hello&str2=world", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response map[string][]string
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the response
	result, exists := response["result"]
	assert.True(t, exists)
	assert.Equal(t, 15, len(result))

	// Check specific values in the result
	assert.Equal(t, "1", result[0])
	assert.Equal(t, "2", result[1])
	assert.Equal(t, "hello", result[2])
	assert.Equal(t, "4", result[3])
	assert.Equal(t, "world", result[4])
	assert.Equal(t, "hello", result[5])
	assert.Equal(t, "7", result[6])
	assert.Equal(t, "8", result[7])
	assert.Equal(t, "hello", result[8])
	assert.Equal(t, "world", result[9])
	assert.Equal(t, "11", result[10])
	assert.Equal(t, "hello", result[11])
	assert.Equal(t, "13", result[12])
	assert.Equal(t, "14", result[13])
	assert.Equal(t, "helloworld", result[14])
}
