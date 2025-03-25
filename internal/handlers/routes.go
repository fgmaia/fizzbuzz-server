package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RegisterRoutes(fiberApp *fiber.App) {
	// Define the rate limiting configuration
	rateLimitConfig := limiter.Config{
		Max:        100,             // Maximum number of requests
		Expiration: 1 * time.Second, // Time frame for the rate limit
	}

	// API Documentation route
	fiberApp.Get("/docs", DocHandler)
	// FizzBuzz endpoint
	fiberApp.Get("/fizzbuzz", FizzbuzzHandler)
	// Stats endpoint
	fiberApp.Get("/stats", Stats)
	// Prometheus metrics endpoint
	fiberApp.Get("/metrics", MetricsHandler)

	_ = rateLimitConfig
	/*
		fiberApp.Get("senior-rh-emp/:cpf", limiter.New(rateLimitConfig), TokenMiddleware(func(ctx *fiber.Ctx) error {
			ulog.Info("senior-rh-emp", "cpf", ctx.Params("cpf"))
			return SeniorRhEmp(ctx)
		}))
	*/
}
