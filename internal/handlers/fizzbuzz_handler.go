package handlers

import (
	"context"
	"fizzbuzz-server/internal/apps"
	"fizzbuzz-server/internal/entities"
	"fizzbuzz-server/internal/metrics"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Tracer name for the FizzBuzz handler
var tracer = otel.Tracer("fizzbuzz-server/internal/handlers")

func FizzbuzzHandler(c *fiber.Ctx) error {
	startTime := time.Now()
	status := "success"

	// Create a span for the handler
	ctx, span := tracer.Start(c.Context(), "FizzbuzzHandler")
	defer span.End()

	req := entities.FizzBuzzRequest{}

	// Bind query parameters
	if err := c.QueryParser(&req); err != nil {
		status = "error_parse"
		span.SetStatus(codes.Error, "Invalid parameter format")
		span.RecordError(err)
		
		// Record metrics for error
		metrics.FizzBuzzRequests.WithLabelValues(status).Inc()
		metrics.FizzBuzzRequestDuration.WithLabelValues(status).Observe(time.Since(startTime).Seconds())
		
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

	// Add request parameters as span attributes
	span.SetAttributes(
		attribute.Int("fizzbuzz.int1", req.Int1),
		attribute.Int("fizzbuzz.int2", req.Int2),
		attribute.Int("fizzbuzz.limit", req.Limit),
		attribute.String("fizzbuzz.str1", req.Str1),
		attribute.String("fizzbuzz.str2", req.Str2),
	)

	// Validate
	validate := c.Locals("validator").(*validator.Validate)
	if err := validate.Struct(req); err != nil {
		status = "error_validation"
		span.SetStatus(codes.Error, "Validation error")
		span.RecordError(err)
		
		// Record metrics for validation error
		metrics.FizzBuzzRequests.WithLabelValues(status).Inc()
		metrics.FizzBuzzRequestDuration.WithLabelValues(status).Observe(time.Since(startTime).Seconds())
		
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Generate result with context for potential tracing in the service layer
	result := generateFizzBuzzWithContext(ctx, req)

	// Update stats
	updateStats(req)

	// Record successful metrics
	metrics.FizzBuzzRequests.WithLabelValues(status).Inc()
	metrics.FizzBuzzRequestDuration.WithLabelValues(status).Observe(time.Since(startTime).Seconds())
	metrics.FizzBuzzResultSize.Observe(float64(len(result)))

	span.SetStatus(codes.Ok, "Success")
	return c.JSON(fiber.Map{
		"result": result,
	})
}

// generateFizzBuzzWithContext calls the FizzBuzz service with context for tracing
func generateFizzBuzzWithContext(ctx context.Context, req entities.FizzBuzzRequest) []string {
	_, span := tracer.Start(ctx, "GenerateFizzBuzz")
	defer span.End()

	result := apps.App().FizzBuzzService.GenerateFizzBuzz(req.Int1, req.Int2, req.Limit, req.Str1, req.Str2)

	// Add result size as an attribute
	resultSize := len(result)
	span.SetAttributes(attribute.Int("fizzbuzz.result.size", resultSize))

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