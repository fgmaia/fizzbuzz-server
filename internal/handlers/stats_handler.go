package handlers

import (
	"fizzbuzz-server/internal/apps"
	"fizzbuzz-server/internal/metrics"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// statsTracer is the tracer for the stats handler
var statsTracer = otel.Tracer("fizzbuzz-server/internal/handlers/stats")

func Stats(c *fiber.Ctx) error {
	startTime := time.Now()
	status := "success"

	// Create a span for the stats handler
	ctx, span := statsTracer.Start(c.Context(), "StatsHandler")
	defer span.End()

	// Acquire read lock with tracing
	lockSpan := trace.SpanFromContext(ctx)
	lockSpan.AddEvent("Acquiring stats read lock")
	stats.Mutex.RLock()
	defer stats.Mutex.RUnlock()
	lockSpan.AddEvent("Acquired stats read lock")

	// Record the number of stats entries
	entriesCount := len(stats.Counts)
	span.SetAttributes(attribute.Int("stats.count", entriesCount))
	
	// Update Prometheus gauge for stats entries count
	metrics.StatsEntriesCount.Set(float64(entriesCount))

	if entriesCount == 0 {
		span.SetStatus(codes.Ok, "No requests yet")
		
		// Record metrics for empty stats
		metrics.StatsRequests.WithLabelValues(status).Inc()
		metrics.StatsRequestDuration.WithLabelValues(status).Observe(time.Since(startTime).Seconds())
		
		return c.JSON(fiber.Map{
			"message": "No requests have been made yet",
		})
	}

	// Find most frequent request with tracing
	_, mostFrequentSpan := statsTracer.Start(ctx, "FindMostFrequentRequest")
	var maxKey string
	maxCount := 0
	for key, count := range stats.Counts {
		if count > maxCount {
			maxCount = count
			maxKey = key
		}
	}
	mostFrequentSpan.SetAttributes(
		attribute.String("stats.most_frequent.key", maxKey),
		attribute.Int("stats.most_frequent.count", maxCount),
	)
	mostFrequentSpan.End()

	// Parse the key back into components with tracing
	_, parseSpan := statsTracer.Start(ctx, "ParseStatsKey")
	parts, err := apps.App().StatsService.ParseStatsKey(maxKey)
	if err != nil {
		status = "error_parse"
		parseSpan.SetStatus(codes.Error, "Failed to parse stats key")
		parseSpan.RecordError(err)
		parseSpan.End()

		span.SetStatus(codes.Error, "Internal stats error")
		span.RecordError(err)
		
		// Record metrics for parse error
		metrics.StatsRequests.WithLabelValues(status).Inc()
		metrics.StatsRequestDuration.WithLabelValues(status).Observe(time.Since(startTime).Seconds())
		
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Internal stats error",
		})
	}
	parseSpan.End()

	// Build response with tracing
	_, responseSpan := statsTracer.Start(ctx, "BuildStatsResponse")
	resp := StatsResponse{}
	resp.MostFrequentRequest.Int1 = parts.Int1
	resp.MostFrequentRequest.Int2 = parts.Int2
	resp.MostFrequentRequest.Limit = parts.Limit
	resp.MostFrequentRequest.Str1 = parts.Str1
	resp.MostFrequentRequest.Str2 = parts.Str2
	resp.MostFrequentRequest.Hits = maxCount

	// Add response attributes to the span
	responseSpan.SetAttributes(
		attribute.Int("stats.response.int1", parts.Int1),
		attribute.Int("stats.response.int2", parts.Int2),
		attribute.Int("stats.response.limit", parts.Limit),
		attribute.String("stats.response.str1", parts.Str1),
		attribute.String("stats.response.str2", parts.Str2),
		attribute.Int("stats.response.hits", maxCount),
	)
	responseSpan.End()

	// Update Prometheus gauge for most frequent request hits
	metrics.MostFrequentRequestHits.Set(float64(maxCount))
	
	// Record successful metrics
	metrics.StatsRequests.WithLabelValues(status).Inc()
	metrics.StatsRequestDuration.WithLabelValues(status).Observe(time.Since(startTime).Seconds())

	span.SetStatus(codes.Ok, "Success")
	return c.JSON(resp)
}