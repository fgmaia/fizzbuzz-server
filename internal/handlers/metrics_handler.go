package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler exposes Prometheus metrics
func MetricsHandler(c *fiber.Ctx) error {
	// Convert the Prometheus handler to a Fiber handler using the adaptor
	return adaptor.HTTPHandler(promhttp.Handler())(c)
}