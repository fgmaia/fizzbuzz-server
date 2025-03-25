package handlers

import (
	"context"
	"fizzbuzz-server/internal/apps"
	"fizzbuzz-server/internal/entities"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func FizzbuzzHandler(c *fiber.Ctx) error {

	req := entities.FizzBuzzRequest{}

	// Bind query parameters
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid parameter format",
		})
	}

	// Set default values if empty
	if req.Str1 == "" {
		req.Str1 = "fizz"
	}
	if req.Str2 == "" {
		req.Str2 = "buzz"
	}

	// Validate
	validate := c.Locals("validator").(*validator.Validate)
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Generate result with context for potential tracing in the service layer
	result := generateFizzBuzzWithContext(c.Context(), req)

	// Update stats
	updateStats(req)

	return c.JSON(fiber.Map{
		"result": result,
	})
}

// generateFizzBuzzWithContext calls the FizzBuzz service with context for tracing
func generateFizzBuzzWithContext(ctx context.Context, req entities.FizzBuzzRequest) []string {
	result := apps.App().FizzBuzzService.GenerateFizzBuzz(req.Int1, req.Int2, req.Limit, req.Str1, req.Str2)
	return result
}

// updateStats updates the request statistics
func updateStats(req entities.FizzBuzzRequest) {
	stats.Mutex.Lock()
	defer stats.Mutex.Unlock()

	key := strings.Join([]string{
		strconv.Itoa(req.Int1),
		strconv.Itoa(req.Int2),
		strconv.Itoa(req.Limit),
		req.Str1,
		req.Str2,
	}, ",")
	stats.Counts[key]++
}
