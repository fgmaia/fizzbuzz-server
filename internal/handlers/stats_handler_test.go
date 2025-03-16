package handlers_test

import (
	"encoding/json"
	"fizzbuzz-server/internal/handlers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestStatsHandler_NoRequests(t *testing.T) {
	// Setup a new app for this test to ensure clean stats
	app := fiber.New()
	validate := validator.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("validator", validate)
		return c.Next()
	})
	app.Get("/stats", handlers.Stats)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response handlers.StatsResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the response - should have zero hits since no requests were made
	assert.Equal(t, 0, response.MostFrequentRequest.Hits)
}

func TestStatsHandler_WithRequests(t *testing.T) {
	// Setup
	app := fiber.New()
	validate := validator.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("validator", validate)
		return c.Next()
	})
	app.Get("/fizzbuzz", handlers.FizzbuzzHandler)
	app.Get("/stats", handlers.Stats)

	// Make several FizzBuzz requests with the same parameters
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		resp.Body.Close()
	}

	// Make a different FizzBuzz request
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=2&int2=7&limit=10&str1=hello&str2=world", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	resp.Body.Close()

	// Now check the stats
	statsReq := httptest.NewRequest(http.MethodGet, "/stats", nil)
	statsReq.Header.Set("Content-Type", "application/json")
	statsResp, err := app.Test(statsReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statsResp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(statsResp.Body)
	assert.NoError(t, err)

	var response handlers.StatsResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the most frequent request
	assert.Equal(t, 3, response.MostFrequentRequest.Int1)
	assert.Equal(t, 5, response.MostFrequentRequest.Int2)
	assert.Equal(t, 15, response.MostFrequentRequest.Limit)
	assert.Equal(t, "fizz", response.MostFrequentRequest.Str1)
	assert.Equal(t, "buzz", response.MostFrequentRequest.Str2)
	assert.Equal(t, 3, response.MostFrequentRequest.Hits)
}

func TestStatsHandler_MultipleTopRequests(t *testing.T) {
	// Setup
	app := fiber.New()
	validate := validator.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("validator", validate)
		return c.Next()
	})
	app.Get("/fizzbuzz", handlers.FizzbuzzHandler)
	app.Get("/stats", handlers.Stats)

	// Make several FizzBuzz requests with different parameters, same number of times
	for i := 0; i < 2; i++ {
		// First set of parameters
		req1 := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil)
		req1.Header.Set("Content-Type", "application/json")
		resp1, _ := app.Test(req1)
		resp1.Body.Close()

		// Second set of parameters
		req2 := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=2&int2=7&limit=10&str1=hello&str2=world", nil)
		req2.Header.Set("Content-Type", "application/json")
		resp2, _ := app.Test(req2)
		resp2.Body.Close()
	}

	// Now check the stats
	statsReq := httptest.NewRequest(http.MethodGet, "/stats", nil)
	statsReq.Header.Set("Content-Type", "application/json")
	statsResp, err := app.Test(statsReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statsResp.StatusCode)

	// Parse the response
	body, err := io.ReadAll(statsResp.Body)
	assert.NoError(t, err)

	var response handlers.StatsResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify that one of the most frequent requests is returned
	// Since both have the same frequency, either could be returned
	assert.Equal(t, 2, response.MostFrequentRequest.Hits)

	// Check that the returned parameters match one of our sets
	isFirstSet := response.MostFrequentRequest.Int1 == 3 &&
		response.MostFrequentRequest.Int2 == 5 &&
		response.MostFrequentRequest.Limit == 15 &&
		response.MostFrequentRequest.Str1 == "fizz" &&
		response.MostFrequentRequest.Str2 == "buzz"

	isSecondSet := response.MostFrequentRequest.Int1 == 2 &&
		response.MostFrequentRequest.Int2 == 7 &&
		response.MostFrequentRequest.Limit == 10 &&
		response.MostFrequentRequest.Str1 == "hello" &&
		response.MostFrequentRequest.Str2 == "world"

	assert.True(t, isFirstSet || isSecondSet, "The most frequent request should match one of our test sets")
}
